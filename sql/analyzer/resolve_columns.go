// Copyright 2020-2021 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dolthub/vitess/go/vt/sqlparser"
	"gopkg.in/src-d/go-errors.v1"

	"github.com/dolthub/go-mysql-server/internal/similartext"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
	"github.com/dolthub/go-mysql-server/sql/plan"
	"github.com/dolthub/go-mysql-server/sql/transform"
)

func validateUniqueTableNames(ctx *sql.Context, a *Analyzer, n sql.Node, scope *Scope, sel RuleSelector) (sql.Node, transform.TreeIdentity, error) {
	// getTableAliases will error if any table name / alias is repeated
	_, err := getTableAliases(n, scope)
	if err != nil {
		return nil, transform.SameTree, err
	}

	return n, transform.SameTree, err
}

// deferredColumn is a wrapper on UnresolvedColumn used to defer the resolution of the column because it may require
// some work done by other analyzer phases.
type deferredColumn struct {
	*expression.UnresolvedColumn
}

func (dc *deferredColumn) DebugString() string {
	return fmt.Sprintf("deferred(%s)", dc.UnresolvedColumn.String())
}

// IsNullable implements the Expression interface.
func (*deferredColumn) IsNullable() bool {
	return true
}

// Children implements the Expression interface.
func (*deferredColumn) Children() []sql.Expression { return nil }

// WithChildren implements the Expression interface.
func (dc *deferredColumn) WithChildren(children ...sql.Expression) (sql.Expression, error) {
	if len(children) != 0 {
		return nil, sql.ErrInvalidChildrenNumber.New(dc, len(children), 0)
	}
	return dc, nil
}

type tableCol struct {
	table string
	col   string
}

func newTableCol(table, col string) tableCol {
	return tableCol{
		table: strings.ToLower(table),
		col:   strings.ToLower(col),
	}
}

var _ sql.Tableable = tableCol{}
var _ sql.Nameable = tableCol{}

func (tc tableCol) Table() string {
	return tc.table
}

func (tc tableCol) Name() string {
	return tc.col
}

type indexedCol struct {
	*sql.Column
	index int
}

// column is the common interface that groups UnresolvedColumn and deferredColumn.
type column interface {
	sql.Nameable
	sql.Tableable
	sql.Expression
}

// scopeLevelSymbols tracks available table and column name symbols at a specific scope level for a query. Each nested
// subquery in a statement represents an additional scope level.
type scopeLevelSymbols struct {
	availableColumns   map[string][]string
	availableAliases   map[string][]*expression.Alias
	availableTables    map[string]string
	availableTableCols map[tableCol]struct{}
}

func newScopeLevelSymbols() *scopeLevelSymbols {
	return &scopeLevelSymbols{
		availableColumns:   make(map[string][]string),
		availableAliases:   make(map[string][]*expression.Alias),
		availableTables:    make(map[string]string),
		availableTableCols: make(map[tableCol]struct{}),
	}
}

// availableNames tracks available table and column name symbols at each scope level for a query, where level 0
// is the node being analyzed, and each additional level is one layer of query scope outward.
type availableNames map[int]*scopeLevelSymbols

// indexColumn adds a column with the given table and column name at the given scope level
func (a availableNames) indexColumn(table, col string, scopeLevel int) {
	col = strings.ToLower(col)
	_, ok := a[scopeLevel]
	if !ok {
		a[scopeLevel] = newScopeLevelSymbols()
	}
	tableLower := strings.ToLower(table)
	if !stringContains(a[scopeLevel].availableColumns[col], tableLower) {
		a[scopeLevel].availableColumns[col] = append(a[scopeLevel].availableColumns[col], tableLower)
		a[scopeLevel].availableTableCols[tableCol{table: tableLower, col: col}] = struct{}{}
	}
}

// levels returns a sorted list of the available nesting scope levels
func (a availableNames) levels() []int {
	levels := make([]int, len(a))
	i := 0
	for l := range a {
		levels[i] = l
		i++
	}
	sort.Ints(levels)
	return levels
}

// indexAlias adds an alias name to track at the specified scope level
func (a availableNames) indexAlias(e *expression.Alias, scopeLevel int) {
	name := strings.ToLower(e.Name())
	_, ok := a[scopeLevel]
	if !ok {
		a[scopeLevel] = newScopeLevelSymbols()
	}
	_, ok = a[scopeLevel].availableAliases[name]
	if !ok {
		a[scopeLevel].availableAliases[name] = make([]*expression.Alias, 0)
	}
	a[scopeLevel].availableAliases[name] = append(a[scopeLevel].availableAliases[name], e)
}

// indexTable adds a table with the given name at the specified scope level
func (a availableNames) indexTable(alias, name string, scopeLevel int) {
	alias = strings.ToLower(alias)
	_, ok := a[scopeLevel]
	if !ok {
		a[scopeLevel] = newScopeLevelSymbols()
	}
	a[scopeLevel].availableTables[alias] = strings.ToLower(name)
}

// nesting levels returns all levels present, from inner to outer
func (a availableNames) nestingLevels() []int {
	levels := make([]int, len(a))
	for level := range a {
		levels = append(levels, level)
	}
	sort.Ints(levels)
	return levels
}

func (a availableNames) tablesAtLevel(scopeLevel int) map[string]string {
	return a[scopeLevel].availableTables
}

func (a availableNames) allTables() []string {
	var allTables []string
	for _, symbols := range a {
		for name, table := range symbols.availableTables {
			allTables = append(allTables, name, table)
		}
	}
	return dedupStrings(allTables)
}

func (a availableNames) tablesForColumnAtLevel(column string, scopeLevel int) []string {
	return a[scopeLevel].availableColumns[column]
}

func (a availableNames) hasTableCol(tc tableCol) bool {
	for scopeLevel := range a {
		_, ok := a[scopeLevel].availableTableCols[tc]
		if ok {
			return true
		}
	}
	return false
}

func dedupStrings(in []string) []string {
	var seen = make(map[string]struct{})
	var result []string
	for _, s := range in {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}

// qualifyColumns assigns a table to any column expressions that don't have one already
func qualifyColumns(ctx *sql.Context, a *Analyzer, n sql.Node, scope *Scope, sel RuleSelector) (sql.Node, transform.TreeIdentity, error) {
	// Calculate the available symbols BEFORE we get into a transform function, since symbols need to be calculated
	// on the full scope; otherwise transform looks at sub-nodes and calculates a partial view of what is available.
	symbols := getAvailableNamesByScope(n, scope)

	return transform.Node(n, func(n sql.Node) (sql.Node, transform.TreeIdentity, error) {
		if _, ok := n.(sql.Expressioner); !ok || n.Resolved() {
			return n, transform.SameTree, nil
		}
		// Updates need to have check constraints qualified since multiple tables could be involved
		checkConstraintsUpdated := false
		if nn, ok := n.(*plan.Update); ok {
			newChecks, err := qualifyCheckConstraints(nn)
			if err != nil {
				return n, transform.SameTree, err
			}
			nn.Checks = newChecks
			n = nn
			checkConstraintsUpdated = true
		}

		newNode, identity, err := transform.OneNodeExprsWithNode(n, func(n sql.Node, e sql.Expression) (sql.Expression, transform.TreeIdentity, error) {
			return qualifyExpression(e, n, symbols)
		})
		if err != nil {
			return n, transform.SameTree, err
		}
		if checkConstraintsUpdated {
			identity = transform.NewTree
		}
		return newNode, identity, nil
	})
}

// qualifyCheckConstraints returns a new set of CheckConstraints created by taking the specified Update node's checks
// and qualifying them to the table involved in the update, including honoring any table aliases specified.
func qualifyCheckConstraints(update *plan.Update) (sql.CheckConstraints, error) {
	checks := update.Checks
	table, alias := getResolvedTableAndAlias(update.Child)

	newExprs := make([]sql.Expression, len(checks))
	for i, checkExpr := range checks.ToExpressions() {
		newExpr, _, err := transform.Expr(checkExpr, func(e sql.Expression) (sql.Expression, transform.TreeIdentity, error) {
			switch e := e.(type) {
			case *expression.UnresolvedColumn:
				if e.Table() == "" {
					tableName := table.Name()
					if alias != "" {
						tableName = alias
					}
					return expression.NewUnresolvedQualifiedColumn(tableName, e.Name()), transform.NewTree, nil
				}
			default:
				// nothing else needed for other types
			}
			return e, transform.SameTree, nil
		})
		if err != nil {
			return nil, err
		}
		newExprs[i] = newExpr
	}

	return checks.FromExpressions(newExprs)
}

func getAvailableNamesByScope(n sql.Node, scope *Scope) availableNames {
	symbols := make(availableNames)

	scopeNodes := make([]sql.Node, 0, 1+len(scope.InnerToOuter()))
	scopeNodes = append(scopeNodes, n)
	scopeNodes = append(scopeNodes, scope.InnerToOuter()...)
	currentScopeLevel := len(scopeNodes)

	// Scope 0 should be top level node, scope 1 should be the current/nested node
	// Examine all columns, from the innermost scope (this node) outward.
	getColumnsInNodes(n.Children(), symbols, currentScopeLevel-1) // TODO: Why don't we pass in the full node n?
	for _, currentScopeNode := range scopeNodes {
		getColumnsInNodes([]sql.Node{currentScopeNode}, symbols, currentScopeLevel-1)

		// TODO: Can we remove this somehow?
		currentScopeLevel--
	}

	// Get table names in all outer scopes and nodes. Inner scoped names will overwrite those from the outer scope.
	// note: we terminate the symbols for this level after finding the first column source
	for scopeLevel, n := range scopeNodes {
		transform.Inspect(n, func(n sql.Node) bool {
			switch n := n.(type) {
			case *plan.SubqueryAlias, *plan.ResolvedTable, *plan.ValueDerivedTable, *plan.RecursiveCte, *information_schema.ColumnsTable, *plan.IndexedTableAccess:
				name := strings.ToLower(n.(sql.Nameable).Name())
				symbols.indexTable(name, name, scopeLevel)
				return false
			case *plan.TableAlias:
				switch t := n.Child.(type) {
				case *plan.ResolvedTable, *plan.UnresolvedTable, *plan.SubqueryAlias,
					*plan.RecursiveTable, *information_schema.ColumnsTable, *plan.IndexedTableAccess:
					name := strings.ToLower(t.(sql.Nameable).Name())
					alias := strings.ToLower(n.Name())
					symbols.indexTable(alias, name, scopeLevel)
				}
				return false
			case *plan.GroupBy:
				// project aliases can overwrite lower namespaces, but importantly,
				// we do not terminate symbol generation.
				for _, e := range n.SelectedExprs {
					if a, ok := e.(*expression.Alias); ok {
						symbols.indexAlias(a, scopeLevel)
					}
				}
			case *plan.Project:
				// project aliases can overwrite lower namespaces, but importantly,
				// we do not terminate symbol generation.
				for _, e := range n.Projections {
					if a, ok := e.(*expression.Alias); ok {
						symbols.indexAlias(a, scopeLevel)
					}
				}
			}
			return true
		})
	}

	return symbols
}

func qualifyExpression(e sql.Expression, node sql.Node, symbols availableNames) (sql.Expression, transform.TreeIdentity, error) {
	switch col := e.(type) {
	case column:
		if col.Resolved() {
			return col, transform.SameTree, nil
		}

		// Skip qualification if an expression has already been identified as an alias reference
		if _, ok := col.(*expression.AliasReference); ok {
			return col, transform.SameTree, nil
		}

		// Skip this step for variables
		if strings.HasPrefix(col.Name(), "@") || strings.HasPrefix(col.Table(), "@") {
			return col, transform.SameTree, nil
		}

		// if there are no tables or columns anywhere in the query, just give up and let another part of the analyzer throw
		// an analysis error. (for some queries, like SHOW statements, this is expected and not an error)
		if len(symbols) == 0 {
			return col, transform.SameTree, nil
		}

		canAccessAliasAtCurrentScope := true
		if _, ok := node.(*plan.Filter); ok {
			// Expression aliases from the same scope are NOT allowed in where/filter clauses
			canAccessAliasAtCurrentScope = false
		}

		// If this column is already qualified, make sure the table name is known
		if col.Table() != "" {
			if validateQualifiedColumn(col, symbols) {
				return col, transform.SameTree, nil
			} else {
				similar := similartext.Find(symbols.allTables(), col.Table())
				return nil, transform.SameTree, sql.ErrTableNotFound.New(col.Table() + similar)
			}
		}

		// Look in all the scopes, inner to outer, to identify the column. Stop as soon as we have a scope with exactly 1
		// match for the column name. If any scope has ambiguity in available column names, that's an error.
		name := strings.ToLower(col.Name())
		levels := symbols.levels()

		// Reverse the levels so that we go from inner scope to outer scopes looking for a match
		sort.Sort(sort.Reverse(sort.IntSlice(levels)))
		maxLevel := levels[0]

		for _, scopeLevel := range levels {
			tablesForColumn := symbols.tablesForColumnAtLevel(name, scopeLevel)
			switch len(tablesForColumn) {
			case 0:
				// This column could be in an outer scope, keep going
				continue
			case 1:
				// This indicates we found a match with an alias definition.
				if tablesForColumn[0] == "" {
					if canAccessAliasAtCurrentScope || scopeLevel != maxLevel {
						// TODO: Make sure there is only one alias registered, otherwise it's an error (or a warning sometimes)
						return expression.NewAliasReference(col.Name()), transform.NewTree, nil
					}
				} else {
					return expression.NewUnresolvedQualifiedColumn(
						tablesForColumn[0],
						col.Name(),
					), transform.NewTree, nil
				}
			default:
				aliasesFound := make([]string, 0, len(tablesForColumn))
				tablesFound := make([]string, 0, len(tablesForColumn))
				for _, tableName := range tablesForColumn {
					if tableName == "" {
						aliasesFound = append(aliasesFound, tableName)
					} else {
						tablesFound = append(tablesFound, tableName)
					}
				}

				// TODO: How about if there is a column name that matches at the same scope, but there is an alias
				//       available at an outer scope that should be preferred? Concrete test cases?
				switch node.(type) {
				case *plan.Sort, *plan.Having:
					// For order by and having clauses... prefer an alias over a column when there is ambiguity
					if len(aliasesFound) == 0 {
						if len(tablesFound) == 1 {
							return expression.NewUnresolvedQualifiedColumn(tablesFound[0], col.Name()), transform.NewTree, nil
						} else if len(tablesFound) > 1 {
							// TODO: This error is specific to OrderBy
							return col, transform.SameTree, sql.ErrAmbiguousColumnInOrderBy.New(col.Name())
						}
					} else if len(aliasesFound) == 1 {
						// TODO: Why are we checking this, too? Need to explain this in a comment...
						if availableAliases, ok := symbols[scopeLevel].availableAliases[col.Name()]; ok {
							if len(availableAliases) > 1 {
								return col, transform.SameTree, sql.ErrAmbiguousColumnInOrderBy.New(col.Name())
							}
						}
						return expression.NewAliasReference(col.Name()), transform.NewTree, nil
					} else if len(aliasesFound) > 1 {
						// TODO: This error is specific to OrderBy
						return col, transform.SameTree, sql.ErrAmbiguousColumnInOrderBy.New(col.Name())
					}
				default:
					// TODO: Need to do more testing on this code and to clean up this nested switch in a switch; hard
					//       to reason about and messy
					// otherwise, prefer the table column...
					if len(tablesFound) == 1 {
						return expression.NewUnresolvedQualifiedColumn(tablesFound[0], col.Name()), transform.NewTree, nil
					} else if len(aliasesFound) == 1 {
						return expression.NewAliasReference(col.Name()), transform.NewTree, nil
					}
				}

				return nil, transform.SameTree, sql.ErrAmbiguousColumnName.New(col.Name(), strings.Join(tablesForColumn, ", "))
			}
		}

		if !canAccessAliasAtCurrentScope {
			// return an error if we can't find a column after searching all levels and we can't hope that this might be an alias
			// TODO: When could this actually be an alias that we don't know about yet?
			// TODO: is it safer to just return a deferred column? Should we do this for
			//return col, transform.SameTree, sql.ErrColumnNotFound.New(col.Name())
			return &deferredColumn{expression.NewUnresolvedQualifiedColumn(col.Table(), col.Name())}, transform.NewTree, nil
		}

		// If there are no tables that have any column with the column name let's just return it as it is. This may be an
		// alias, so we'll wait for the reorder of the projection to resolve it.
		return col, transform.SameTree, nil
	case *expression.Star:
		// Make sure that any qualified stars reference known tables
		if col.Table != "" {
			//nestingLevels := symbols.nestingLevels()
			tableFound := false
			for level := range symbols {
				tables := symbols.tablesAtLevel(level)
				if _, ok := tables[strings.ToLower(col.Table)]; ok {
					tableFound = true
					break
				}
			}
			if !tableFound {
				return nil, transform.SameTree, sql.ErrTableNotFound.New(col.Table)
			}
		}
		return col, transform.SameTree, nil
	default:
		// If any other kind of expression has a star, just replace it
		// with an unqualified star because it cannot be expanded.
		return transform.Expr(e, func(e sql.Expression) (sql.Expression, transform.TreeIdentity, error) {
			if _, ok := e.(*expression.Star); ok {
				return expression.NewStar(), transform.NewTree, nil
			}
			return e, transform.SameTree, nil
		})
	}
}

// validateQualifiedColumn returns true if the table name of the specified column is a valid table name symbol, meaning
// it is available in some scope of the current statement. If a valid table name symbol can't be found, false is returned.
func validateQualifiedColumn(col column, symbols availableNames) bool {
	for scopeLevel := range symbols {
		tables := symbols.tablesAtLevel(scopeLevel)
		if _, ok := tables[strings.ToLower(col.Table())]; ok {
			return true
		}
	}

	if symbols.hasTableCol(tableCol{table: strings.ToLower(col.Table()), col: strings.ToLower(col.Name())}) {
		return true
	}

	return false
}

func getColumnsInNodes(nodes []sql.Node, names availableNames, scopeLevel int) {
	indexExpressions := func(exprs []sql.Expression) {
		for _, e := range exprs {
			switch e := e.(type) {
			case *expression.Alias:
				// Mark this as an available column; we'll record the Alias information later
				names.indexColumn("", e.Name(), scopeLevel)
			case *expression.GetField:
				names.indexColumn(e.Table(), e.Name(), scopeLevel)
			}
		}
	}

	for _, node := range nodes {
		switch n := node.(type) {
		case *plan.TableAlias, *plan.ResolvedTable, *plan.SubqueryAlias, *plan.ValueDerivedTable, *plan.RecursiveTable, *information_schema.ColumnsTable:
			for _, col := range n.Schema() {
				names.indexColumn(col.Source, col.Name, scopeLevel)
			}
		case *plan.Project:
			indexExpressions(n.Projections)
			getColumnsInNodes(n.Children(), names, scopeLevel)
		case *plan.GroupBy:
			indexExpressions(n.SelectedExprs)
			getColumnsInNodes(n.Children(), names, scopeLevel)
		case *plan.Window:
			indexExpressions(n.SelectExprs)
			getColumnsInNodes(n.Children(), names, scopeLevel)
		default:
			getColumnsInNodes(n.Children(), names, scopeLevel)
		}
	}
}

var errGlobalVariablesNotSupported = errors.NewKind("can't resolve global variable, %s was requested")

const (
	sessionTable  = "@@" + sqlparser.SessionStr
	sessionPrefix = sqlparser.SessionStr + "."
	globalPrefix  = sqlparser.GlobalStr + "."
)

// resolveColumns replaces UnresolvedColumn expressions with GetField expressions for the appropriate numbered field in
// the expression's child node.
func resolveColumns(ctx *sql.Context, a *Analyzer, n sql.Node, scope *Scope, sel RuleSelector) (sql.Node, transform.TreeIdentity, error) {
	span, ctx := ctx.Span("resolve_columns")
	defer span.End()

	return transform.Node(n, func(n sql.Node) (sql.Node, transform.TreeIdentity, error) {
		if n.Resolved() {
			return n, transform.SameTree, nil
		}

		if _, ok := n.(sql.Expressioner); !ok {
			return n, transform.SameTree, nil
		}

		// We need to use the schema, so all children must be resolved.
		// TODO: also enforce the equivalent constraint for outer scopes. More complicated, because the outer scope can't
		//  be Resolved() owing to a child expression (the one being evaluated) not being resolved yet.
		for _, c := range n.Children() {
			if !c.Resolved() {
				return n, transform.SameTree, nil
			}
		}

		columns, err := indexColumns(ctx, a, n, scope)
		if err != nil {
			return nil, transform.SameTree, err
		}

		return transform.OneNodeExprsWithNode(n, func(n sql.Node, e sql.Expression) (sql.Expression, transform.TreeIdentity, error) {
			uc, ok := e.(column)
			if !ok || e.Resolved() {
				return e, transform.SameTree, nil
			}

			return resolveColumnExpression(a, n, uc, columns)
		})
	})
}

// indexColumns returns a map of column identifiers to their index in the node's schema. Columns from outer scopes are
// included as well, with lower indexes (prepended to node schema) but lower precedence (overwritten by inner nodes in
// map)
func indexColumns(ctx *sql.Context, a *Analyzer, n sql.Node, scope *Scope) (map[tableCol]indexedCol, error) {
	var columns = make(map[tableCol]indexedCol)
	var idx int

	indexColumn := func(col *sql.Column) {
		columns[tableCol{
			table: strings.ToLower(col.Source),
			col:   strings.ToLower(col.Name),
		}] = indexedCol{col, idx}
		idx++
	}

	indexSchema := func(n sql.Schema) {
		for _, col := range n {
			indexColumn(col)
		}
	}

	var indexColumnExpr func(e sql.Expression)
	indexColumnExpr = func(e sql.Expression) {
		switch e := e.(type) {
		case *expression.Alias:
			// Aliases get indexed twice with the same index number: once with the aliased name and once with the
			// underlying name
			indexColumn(transform.ExpressionToColumn(e))
			idx--
			indexColumnExpr(e.Child)
		default:
			indexColumn(transform.ExpressionToColumn(e))
		}
	}

	indexChildNode := func(n sql.Node) {
		switch n := n.(type) {
		case *plan.Project:
			for _, e := range n.Projections {
				indexColumnExpr(e)
			}
		case *plan.GroupBy:
			for _, e := range n.SelectedExprs {
				indexColumnExpr(e)
			}
		case *plan.Window:
			for _, e := range n.SelectExprs {
				indexColumnExpr(e)
			}
		case *plan.Values:
			// values nodes don't have a schema to index like other nodes that provide columns
		default:
			indexSchema(n.Schema())
		}
	}

	// Index the columns in the outer scope, outer to inner. This means inner scope columns will overwrite the outer
	// ones of the same name. This matches the MySQL scope precedence rules.
	indexSchema(scope.Schema())

	// For the innermost scope (the node being evaluated), look at the schemas of the children instead of this node
	// itself. Skip this for DDL nodes that handle indexing separately.
	shouldIndexChildNode := true
	switch n.(type) {
	case *plan.AddColumn, *plan.ModifyColumn:
		shouldIndexChildNode = false
	case *plan.RecursiveCte, *plan.Union:
		shouldIndexChildNode = false
	}

	if shouldIndexChildNode {
		for _, child := range n.Children() {
			indexChildNode(child)
		}
	}

	// For certain DDL nodes, we have to do more work
	indexSchemaForDefaults := func(column *sql.Column, order *sql.ColumnOrder, sch sql.Schema) {
		tblSch := make(sql.Schema, len(sch))
		copy(tblSch, sch)
		if order == nil {
			tblSch = append(tblSch, column)
		} else if order.First {
			tblSch = append(sql.Schema{column}, tblSch...)
		} else { // must be After
			index := 1
			afterColumn := strings.ToLower(order.AfterColumn)
			for _, col := range tblSch {
				if strings.ToLower(col.Name) == afterColumn {
					break
				}
				index++
			}
			if index <= len(tblSch) {
				tblSch = append(tblSch, nil)
				copy(tblSch[index+1:], tblSch[index:])
				tblSch[index] = column
			}
		}
		for _, col := range tblSch {
			columns[tableCol{
				table: "",
				col:   strings.ToLower(col.Name),
			}] = indexedCol{col, idx}
			columns[tableCol{
				table: strings.ToLower(col.Source),
				col:   strings.ToLower(col.Name),
			}] = indexedCol{col, idx}
			idx++
		}
	}

	switch node := n.(type) {
	case *plan.CreateTable: // For this node in particular, the columns will only come into existence after the analyzer step, so we forge them here.
		for _, col := range node.CreateSchema.Schema {
			columns[tableCol{
				table: "",
				col:   strings.ToLower(col.Name),
			}] = indexedCol{col, idx}
			columns[tableCol{
				table: strings.ToLower(col.Source),
				col:   strings.ToLower(col.Name),
			}] = indexedCol{col, idx}
			idx++
		}
	case *plan.AddColumn: // Add/Modify need to have the full column set in order to resolve a default expression.
		tbl := node.Table
		indexSchemaForDefaults(node.Column(), node.Order(), tbl.Schema())
	case *plan.ModifyColumn:
		tbl := node.Table
		indexSchemaForDefaults(node.NewColumn(), node.Order(), tbl.Schema())
	case *plan.RecursiveCte, *plan.Union:
		// opaque nodes have derived schemas
		// TODO also subquery aliases?
		indexChildNode(node.(sql.BinaryNode).Left())
	}

	return columns, nil
}

func resolveColumnExpression(a *Analyzer, n sql.Node, e column, columns map[tableCol]indexedCol) (sql.Expression, transform.TreeIdentity, error) {
	name := strings.ToLower(e.Name())
	table := strings.ToLower(e.Table())
	col, ok := columns[tableCol{table, name}]
	if !ok {
		switch uc := e.(type) {
		case *expression.UnresolvedColumn:
			// Defer the resolution of the column to give the analyzer more
			// time to resolve other parts so this can be resolved.
			a.Log("deferring resolution of column %s", e)
			return &deferredColumn{uc}, transform.NewTree, nil
		default:
			if table != "" {
				return nil, transform.SameTree, sql.ErrTableColumnNotFound.New(e.Table(), e.Name())
			}

			// This means the expression is either a non-existent column or an alias defined in the same projection.
			// Check for the latter first.
			aliasesInNode := aliasesDefinedInNode(n)
			if stringContains(aliasesInNode, name) {
				return nil, transform.SameTree, sql.ErrMisusedAlias.New(name)
			}

			return nil, transform.SameTree, sql.ErrColumnNotFound.New(e.Name())
		}
	}

	a.Log("column %s resolved to GetFieldWithTable: idx %d, typ %s, table %s, name %s, nullable %t",
		e, col.index, col.Type, col.Source, col.Name, col.Nullable)
	return expression.NewGetFieldWithTable(
		col.index,
		col.Type,
		col.Source,
		col.Name,
		col.Nullable,
	), transform.NewTree, nil
}

// pushdownGroupByAliases reorders the aggregation in a groupby so aliases defined in it can be resolved in the grouping
// of the groupby. To do so, all aliases are pushed down to a projection node under the group by.
func pushdownGroupByAliases(ctx *sql.Context, a *Analyzer, n sql.Node, scope *Scope, sel RuleSelector) (sql.Node, transform.TreeIdentity, error) {
	if n.Resolved() {
		return n, transform.SameTree, nil
	}

	// replacedAliases is a map of original expression string to alias that has been pushed down below the GroupBy in
	// the new projection node.
	replacedAliases := make(map[string]string)
	var err error
	return transform.Node(n, func(n sql.Node) (sql.Node, transform.TreeIdentity, error) {
		// For any Expressioner node above the GroupBy, we need to apply the same alias replacement as we did in the
		// GroupBy itself.
		ex, ok := n.(sql.Expressioner)
		if ok && len(replacedAliases) > 0 {
			newExprs, same := replaceExpressionsWithAliases(ex.Expressions(), replacedAliases)
			if !same {
				n, err = ex.WithExpressions(newExprs...)
				return n, transform.NewTree, err
			}
		}

		g, ok := n.(*plan.GroupBy)
		if n.Resolved() || !ok || len(g.GroupByExprs) == 0 {
			return n, transform.SameTree, nil
		}
		if !ok || len(g.GroupByExprs) == 0 {
			return n, transform.SameTree, nil
		}

		// The reason we have two sets of columns, one for grouping and
		// one for aggregate is because an alias can redefine a column name
		// of the child schema. In the grouping, if that column is referenced
		// it refers to the alias, and not the one in the child. However,
		// in the aggregate, aliases in that same aggregate cannot be used,
		// so it refers to the column in the child node.
		var groupingColumns = make(map[string]*expression.UnresolvedColumn)
		for _, g := range g.GroupByExprs {
			for _, n := range findAllColumns(g) {
				groupingColumns[strings.ToLower(n.Name())] = n
			}
		}

		var selectedColumns = make(map[string]*expression.UnresolvedColumn)
		for _, agg := range g.SelectedExprs {
			// This alias is going to be pushed down, so don't bother gathering
			// its requirements.
			if alias, ok := agg.(*expression.Alias); ok {
				if _, ok := groupingColumns[strings.ToLower(alias.Name())]; ok {
					continue
				}
			}

			for _, n := range findAllColumns(agg) {
				selectedColumns[strings.ToLower(n.Name())] = n
			}
		}

		var newSelectedExprs []sql.Expression
		var projection []sql.Expression
		// Aliases will keep the aliases that have been pushed down and their
		// index in the new aggregate.
		var aliases = make(map[string]int)

		var needsReorder bool
		for _, expr := range g.SelectedExprs {
			alias, ok := expr.(*expression.Alias)
			// Note that aliases of aggregations cannot be used in the grouping
			// because the grouping is needed before computing the aggregation.
			if !ok || containsAggregation(alias) {
				newSelectedExprs = append(newSelectedExprs, expr)
				continue
			}

			name := strings.ToLower(alias.Name())
			// Only if the alias is required in the grouping set needsReorder
			// to true. If it's not required, there's no need for a reorder if
			// no other alias is required.
			_, ok = groupingColumns[name]
			if ok && groupingColumns[name].Table() == "" {
				aliases[name] = len(newSelectedExprs)
				needsReorder = true
				delete(groupingColumns, name)

				projection = append(projection, expr)
				replacedAliases[alias.Child.String()] = alias.Name()
				newSelectedExprs = append(newSelectedExprs, expression.NewAliasReference(alias.Name()))
			} else {
				newSelectedExprs = append(newSelectedExprs, expr)
			}
		}

		if !needsReorder {
			return n, transform.SameTree, nil
		}

		// Any replacements of aliases in the select expression must be mirrored in the group by, replacing any aliased
		// expressions with a reference to that alias. This is so that we can directly compare the group by an select
		// expressions for validation, which requires us to know that (table.column as col) and (table.column) are the
		// same expressions. So if we replace one, replace both.
		// TODO: this is pretty fragile and relies on string matching, need a better solution
		newGroupBys, _ := replaceExpressionsWithAliases(g.GroupByExprs, replacedAliases)

		// Instead of iterating columns directly, we want them sorted so the
		// executions of the rule are consistent.
		var missingCols = make([]*expression.UnresolvedColumn, 0, len(selectedColumns)+len(groupingColumns))
		for _, col := range selectedColumns {
			missingCols = append(missingCols, col)
		}
		for _, col := range groupingColumns {
			missingCols = append(missingCols, col)
		}

		sort.SliceStable(missingCols, func(i, j int) bool {
			return missingCols[i].Name() < missingCols[j].Name()
		})

		var renames = make(map[string]string)
		// All columns required by expressions in both grouping and aggregation
		// must also be projected in the new projection node or they will not
		// be able to resolve.
		for _, col := range missingCols {
			name := col.Name()
			// If an alias has been pushed down with the same name as a missing
			// column, there will be a conflict of names. We must find an unique name
			// for the missing column.
			if _, ok := aliases[name]; ok {
				for i := 1; ; i++ {
					name = fmt.Sprintf("%s_%02d", col.Name(), i)
					if _, ok := selectedColumns[name]; !ok {
						break
					} else if _, ok := groupingColumns[name]; !ok {
						break
					}
				}
			}

			if name == col.Name() {
				projection = append(projection, col)
			} else {
				renames[col.Name()] = name
				projection = append(projection, expression.NewAlias(name, col))
			}
		}

		// If there is any name conflict between columns we need to rename every
		// usage inside the aggregate.
		if len(renames) > 0 {
			for i, expr := range newSelectedExprs {
				var err error
				newSelectedExprs[i], _, err = transform.Expr(expr, func(e sql.Expression) (sql.Expression, transform.TreeIdentity, error) {
					col, ok := e.(*expression.UnresolvedColumn)
					if ok {
						// We need to make sure we don't rename the reference to the
						// pushed down alias.
						if to, ok := renames[col.Name()]; ok && aliases[col.Name()] != i {
							return expression.NewUnresolvedColumn(to), transform.NewTree, nil
						}
					}

					return e, transform.SameTree, nil
				})
				if err != nil {
					return nil, transform.SameTree, err
				}
			}
		}

		return plan.NewGroupBy(
			newSelectedExprs, newGroupBys,
			plan.NewProject(projection, g.Child),
		), transform.NewTree, nil
	})
}

// replaceExpressionsWithAliases replaces any expressions in the slice given that match the map of aliases given with
// their alias expression or alias name. This is necessary when pushing aliases down the tree, since we introduce a
// projection node that effectively erases the original columns of a table.
func replaceExpressionsWithAliases(exprs []sql.Expression, replacedAliases map[string]string) ([]sql.Expression, transform.TreeIdentity) {
	newMap := make(map[string]struct{})
	for _, aliasName := range replacedAliases {
		newMap[aliasName] = struct{}{}
	}

	var newExprs []sql.Expression
	var expr sql.Expression
	for i := range exprs {
		expr = exprs[i]
		if alias, ok := replacedAliases[expr.String()]; ok {
			if newExprs == nil {
				newExprs = make([]sql.Expression, len(exprs))
				copy(newExprs, exprs)
			}
			newExprs[i] = expression.NewAliasReference(alias)
		} else if uc, ok := expr.(*expression.UnresolvedColumn); ok && uc.Table() == "" {
			if _, ok := newMap[uc.Name()]; ok {
				if newExprs == nil {
					newExprs = make([]sql.Expression, len(exprs))
					copy(newExprs, exprs)
				}
				newExprs[i] = expression.NewAliasReference(uc.Name())
			}
		}
	}
	if len(newExprs) > 0 {
		return newExprs, transform.NewTree
	}
	return exprs, transform.SameTree
}

func findAllColumns(e sql.Expression) []*expression.UnresolvedColumn {
	var cols []*expression.UnresolvedColumn
	sql.Inspect(e, func(e sql.Expression) bool {
		col, ok := e.(*expression.UnresolvedColumn)
		if ok {
			cols = append(cols, col)
		}
		return true
	})
	return cols
}
