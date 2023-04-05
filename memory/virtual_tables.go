package memory

import (
	"bytes"
	"io"

	"github.com/dolthub/go-mysql-server/sql"
)

type VirtualTable struct {
	name     string
	database string
	schema   sql.PrimaryKeySchema
}

type virtualSpatialTableIter struct {
	columns                []int
	rows                   []sql.Row
	pos                    int
	ord                    int
	minX, minY, maxX, maxY float64
}

var _ sql.RowIter = (*virtualSpatialTableIter)(nil)
var _ sql.RowIter2 = (*virtualSpatialTableIter)(nil)

func (i *virtualSpatialTableIter) Next(ctx *sql.Context) (sql.Row, error) {
	return nil, io.EOF
}

func (i *virtualSpatialTableIter) Next2(ctx *sql.Context, frame *sql.RowFrame) error {
	r, err := i.Next(ctx)
	if err != nil {
		return err
	}

	for _, v := range r {
		x, err := sql.ConvertToValue(v)
		if err != nil {
			return err
		}
		frame.Append(x)
	}

	return nil
}

func (i *virtualSpatialTableIter) Close(ctx *sql.Context) error {
	return nil
}

func NewVirtualTable(name string, schema sql.PrimaryKeySchema) (*VirtualTable, error) {
	return &VirtualTable{
		name:   name,
		schema: schema,
	}, nil
}

func (t *VirtualTable) GetPartition(key string) []sql.Row {
	return []sql.Row{}
}

// PartitionRows implements sql.Table
func (t *VirtualTable) PartitionRows(*sql.Context, sql.Partition) (sql.RowIter, error) {
	return &virtualSpatialTableIter{}, nil
}

// LookupPartitions implements sql.IndexedTable
func (*VirtualTable) LookupPartitions(*sql.Context, sql.IndexLookup) (sql.PartitionIter, error) {
	panic("unimplemented")
}

// Name implements sql.Table
func (t VirtualTable) Name() string {
	return t.name
}

// String implements sql.Table
func (t VirtualTable) String() string {
	return t.name
}

// Collation implements sql.Table
func (VirtualTable) Collation() sql.CollationID {
	return sql.Collation_utf8mb4_general_ci
}

// Partitions implements sql.Table
func (VirtualTable) Partitions(*sql.Context) (sql.PartitionIter, error) {
	return &virtualPartitionIter{}, nil
}

// Schema implements sql.Table
func (t *VirtualTable) Schema() sql.Schema {
	return t.schema.Schema
}

func (t *VirtualTable) AddColumn(ctx *sql.Context, column *sql.Column, order *sql.ColumnOrder) error {
	return sql.ErrUnsupportedFeature.New("ADD_COLUMN")
}
func (t *VirtualTable) AddForeignKey(ctx *sql.Context, fk sql.ForeignKeyConstraint) error {
	return sql.ErrUnsupportedFeature.New("ADD_FOREIGN_KEY")
}
func (t *VirtualTable) AnalyzeTable(ctx *sql.Context) error {
	return sql.ErrUnsupportedFeature.New("AnalyzeTable")
}
func (t *VirtualTable) AutoIncrementSetter(*sql.Context) sql.AutoIncrementSetter {
	return nil
}
func (t *VirtualTable) CreateCheck(_ *sql.Context, check *sql.CheckDefinition) error {
	return sql.ErrUnsupportedFeature.New("CREATE_CHECK")
}
func (t *VirtualTable) CreateIndex(ctx *sql.Context, idx sql.IndexDef) error {
	return sql.ErrUnsupportedFeature.New("CREATE_INDEX")
}
func (t *VirtualTable) CreateIndexForForeignKey(ctx *sql.Context, idx sql.IndexDef) error {
	return sql.ErrUnsupportedFeature.New("CreateIndexForForeignKey")
}
func (t *VirtualTable) CreatePrimaryKey(ctx *sql.Context, columns []sql.IndexColumn) error {
	return sql.ErrUnsupportedFeature.New("CreatePrimaryKey")
}
func (t *VirtualTable) DataLength(ctx *sql.Context) (uint64, error) {
	return 0, nil
}
func (t *VirtualTable) DebugString() string {
	return t.name
}
func (t *VirtualTable) Deleter(*sql.Context) sql.RowDeleter {
	return nil
}
func (t *VirtualTable) DropCheck(ctx *sql.Context, chName string) error {
	return sql.ErrUnsupportedFeature.New("DropCheck")
}
func (t *VirtualTable) DropColumn(ctx *sql.Context, columnName string) error {
	return sql.ErrUnsupportedFeature.New("DropColumn")
}
func (t *VirtualTable) DropForeignKey(ctx *sql.Context, fkName string) error {
	return sql.ErrUnsupportedFeature.New("DropForeignKey")
}
func (t *VirtualTable) DropIndex(ctx *sql.Context, indexName string) error {
	return sql.ErrUnsupportedFeature.New("DropIndex")
}
func (t *VirtualTable) DropPrimaryKey(ctx *sql.Context) error {
	return sql.ErrUnsupportedFeature.New("DropPrimaryKey")
}
func (t *VirtualTable) EnablePrimaryKeyIndexes() {

}
func (t *VirtualTable) Filters() []sql.Expression {
	return nil
}
func (t *VirtualTable) GetChecks(_ *sql.Context) ([]sql.CheckDefinition, error) {
	return nil, nil
}
func (t *VirtualTable) GetDeclaredForeignKeys(ctx *sql.Context) ([]sql.ForeignKeyConstraint, error) {
	return nil, nil
}
func (t *VirtualTable) GetForeignKeyEditor(ctx *sql.Context) sql.ForeignKeyEditor {
	return nil
}
func (t *VirtualTable) GetIndexes(ctx *sql.Context) ([]sql.Index, error) {
	return nil, nil
}
func (t *VirtualTable) GetNextAutoIncrementValue(ctx *sql.Context, insertVal interface{}) (uint64, error) {
	return 0, sql.ErrUnsupportedGISType.New("GetNextAutoIncrementValue")
}
func (t *VirtualTable) GetReferencedForeignKeys(ctx *sql.Context) ([]sql.ForeignKeyConstraint, error) {
	return nil, nil
}
func (t *VirtualTable) HandledFilters(filters []sql.Expression) []sql.Expression {
	return nil
}
func (t *VirtualTable) IndexKeyValues(ctx *sql.Context, colNames []string) (sql.PartitionIndexKeyValueIter, error) {
	return nil, nil
}
func (t *VirtualTable) IndexedAccess(i sql.IndexLookup) sql.IndexedTable {
	return t
}
func (t *VirtualTable) Insert(ctx *sql.Context, row sql.Row) error {
	return sql.ErrUnsupportedFeature.New("DropPrimaryKey")
}
func (t *VirtualTable) Inserter(*sql.Context) sql.RowInserter {
	return nil
}

func (t *VirtualTable) ModifyColumn(ctx *sql.Context, columnName string, column *sql.Column, order *sql.ColumnOrder) error {
	return sql.ErrUnsupportedFeature.New("ModifyColumn")
}
func (t *VirtualTable) ModifyDefaultCollation(ctx *sql.Context, collation sql.CollationID) error {
	return sql.ErrUnsupportedFeature.New("ModifyDefaultCollation")
}
func (t *VirtualTable) ModifyStoredCollation(ctx *sql.Context, collation sql.CollationID) error {
	return sql.ErrUnsupportedFeature.New("ModifyStoredCollation")
}
func (t *VirtualTable) PartitionCount(ctx *sql.Context) (int64, error) {
	return 1, nil
}
func (t *VirtualTable) PartitionRows2(ctx *sql.Context, partition sql.Partition) (sql.RowIter2, error) {
	return nil, sql.ErrUnsupportedFeature.New("PartitionRows2")
}
func (t *VirtualTable) PeekNextAutoIncrementValue(*sql.Context) (uint64, error) {
	return 0, sql.ErrUnsupportedFeature.New("PeekNextAutoIncrementValue")
}
func (t *VirtualTable) PrimaryKeySchema() sql.PrimaryKeySchema {
	return sql.NewPrimaryKeySchema(nil)
}
func (t *VirtualTable) Projections() []string {
	return nil
}
func (t *VirtualTable) RenameIndex(ctx *sql.Context, fromIndexName string, toIndexName string) error {
	return sql.ErrUnsupportedFeature.New("RenameIndex")
}
func (t *VirtualTable) Replacer(*sql.Context) sql.RowReplacer {
	return nil
}
func (t *VirtualTable) RowCount(ctx *sql.Context) (uint64, error) {
	return 0, nil
}
func (t *VirtualTable) SetForeignKeyResolved(ctx *sql.Context, fkName string) error {
	return sql.ErrUnsupportedFeature.New("SetForeignKeyResolved")
}
func (t *VirtualTable) Truncate(ctx *sql.Context) (int, error) {
	return 0, sql.ErrUnsupportedFeature.New("Truncate")
}
func (t *VirtualTable) UpdateForeignKey(ctx *sql.Context, fkName string, fk sql.ForeignKeyConstraint) error {
	return sql.ErrUnsupportedFeature.New("UpdateForeignKey")
}
func (t *VirtualTable) Updater(*sql.Context) sql.RowUpdater {
	return nil
}
func (t *VirtualTable) WithDriverIndexLookup(lookup sql.DriverIndexLookup) sql.Table {
	return t
}
func (t *VirtualTable) WithProjections(cols []string) sql.Table {
	return t
}

type virtualPartitionIter struct {
	pos int
}

type virtualPartition struct{}

var defaultKey = bytes.NewBufferString("default").Bytes()

// Key implements sql.Partition
func (p *virtualPartition) Key() []byte {
	return defaultKey
}

// Close implements sql.PartitionIter
func (*virtualPartitionIter) Close(*sql.Context) error {
	return nil
}

func (i *virtualPartitionIter) Next(*sql.Context) (sql.Partition, error) {
	if i.pos == 0 {
		i.pos++
		return &virtualPartition{}, nil
	}
	return nil, io.EOF
}
