// Copyright 2022 Dolthub, Inc.
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

// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package serial

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type PrivilegeSetColumn struct {
	_tab flatbuffers.Table
}

func GetRootAsPrivilegeSetColumn(buf []byte, offset flatbuffers.UOffsetT) *PrivilegeSetColumn {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &PrivilegeSetColumn{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsPrivilegeSetColumn(buf []byte, offset flatbuffers.UOffsetT) *PrivilegeSetColumn {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &PrivilegeSetColumn{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *PrivilegeSetColumn) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *PrivilegeSetColumn) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *PrivilegeSetColumn) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *PrivilegeSetColumn) Privs(j int) int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetInt32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *PrivilegeSetColumn) PrivsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *PrivilegeSetColumn) MutatePrivs(j int, n int32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateInt32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func PrivilegeSetColumnStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func PrivilegeSetColumnAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(name), 0)
}
func PrivilegeSetColumnAddPrivs(builder *flatbuffers.Builder, privs flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(privs), 0)
}
func PrivilegeSetColumnStartPrivsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func PrivilegeSetColumnEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type PrivilegeSetTable struct {
	_tab flatbuffers.Table
}

func GetRootAsPrivilegeSetTable(buf []byte, offset flatbuffers.UOffsetT) *PrivilegeSetTable {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &PrivilegeSetTable{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsPrivilegeSetTable(buf []byte, offset flatbuffers.UOffsetT) *PrivilegeSetTable {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &PrivilegeSetTable{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *PrivilegeSetTable) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *PrivilegeSetTable) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *PrivilegeSetTable) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *PrivilegeSetTable) Privs(j int) int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetInt32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *PrivilegeSetTable) PrivsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *PrivilegeSetTable) MutatePrivs(j int, n int32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateInt32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *PrivilegeSetTable) Columns(obj *PrivilegeSetColumn, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *PrivilegeSetTable) ColumnsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func PrivilegeSetTableStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func PrivilegeSetTableAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(name), 0)
}
func PrivilegeSetTableAddPrivs(builder *flatbuffers.Builder, privs flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(privs), 0)
}
func PrivilegeSetTableStartPrivsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func PrivilegeSetTableAddColumns(builder *flatbuffers.Builder, columns flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(columns), 0)
}
func PrivilegeSetTableStartColumnsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func PrivilegeSetTableEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type PrivilegeSetDatabase struct {
	_tab flatbuffers.Table
}

func GetRootAsPrivilegeSetDatabase(buf []byte, offset flatbuffers.UOffsetT) *PrivilegeSetDatabase {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &PrivilegeSetDatabase{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsPrivilegeSetDatabase(buf []byte, offset flatbuffers.UOffsetT) *PrivilegeSetDatabase {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &PrivilegeSetDatabase{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *PrivilegeSetDatabase) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *PrivilegeSetDatabase) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *PrivilegeSetDatabase) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *PrivilegeSetDatabase) Privs(j int) int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetInt32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *PrivilegeSetDatabase) PrivsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *PrivilegeSetDatabase) MutatePrivs(j int, n int32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateInt32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *PrivilegeSetDatabase) Tables(obj *PrivilegeSetTable, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *PrivilegeSetDatabase) TablesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func PrivilegeSetDatabaseStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func PrivilegeSetDatabaseAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(name), 0)
}
func PrivilegeSetDatabaseAddPrivs(builder *flatbuffers.Builder, privs flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(privs), 0)
}
func PrivilegeSetDatabaseStartPrivsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func PrivilegeSetDatabaseAddTables(builder *flatbuffers.Builder, tables flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(tables), 0)
}
func PrivilegeSetDatabaseStartTablesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func PrivilegeSetDatabaseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type PrivilegeSet struct {
	_tab flatbuffers.Table
}

func GetRootAsPrivilegeSet(buf []byte, offset flatbuffers.UOffsetT) *PrivilegeSet {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &PrivilegeSet{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsPrivilegeSet(buf []byte, offset flatbuffers.UOffsetT) *PrivilegeSet {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &PrivilegeSet{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *PrivilegeSet) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *PrivilegeSet) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *PrivilegeSet) GlobalStatic(j int) int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetInt32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *PrivilegeSet) GlobalStaticLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *PrivilegeSet) MutateGlobalStatic(j int, n int32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateInt32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *PrivilegeSet) GlobalDynamic(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *PrivilegeSet) GlobalDynamicLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *PrivilegeSet) Databases(obj *PrivilegeSetDatabase, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *PrivilegeSet) DatabasesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func PrivilegeSetStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func PrivilegeSetAddGlobalStatic(builder *flatbuffers.Builder, globalStatic flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(globalStatic), 0)
}
func PrivilegeSetStartGlobalStaticVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func PrivilegeSetAddGlobalDynamic(builder *flatbuffers.Builder, globalDynamic flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(globalDynamic), 0)
}
func PrivilegeSetStartGlobalDynamicVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func PrivilegeSetAddDatabases(builder *flatbuffers.Builder, databases flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(databases), 0)
}
func PrivilegeSetStartDatabasesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func PrivilegeSetEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type User struct {
	_tab flatbuffers.Table
}

func GetRootAsUser(buf []byte, offset flatbuffers.UOffsetT) *User {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &User{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsUser(buf []byte, offset flatbuffers.UOffsetT) *User {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &User{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *User) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *User) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *User) User() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *User) Host() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *User) PrivilegeSet(obj *PrivilegeSet) *PrivilegeSet {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(PrivilegeSet)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *User) Plugin() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *User) Password() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *User) PasswordLastChanged() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *User) MutatePasswordLastChanged(n int64) bool {
	return rcv._tab.MutateInt64Slot(14, n)
}

func (rcv *User) Locked() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *User) MutateLocked(n bool) bool {
	return rcv._tab.MutateBoolSlot(16, n)
}

func (rcv *User) Attributes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *User) Identity() []byte {
	// TODO: offset?
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func UserStart(builder *flatbuffers.Builder) {
	builder.StartObject(9)
}
func UserAddUser(builder *flatbuffers.Builder, user flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(user), 0)
}
func UserAddHost(builder *flatbuffers.Builder, host flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(host), 0)
}
func UserAddPrivilegeSet(builder *flatbuffers.Builder, privilegeSet flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(privilegeSet), 0)
}
func UserAddPlugin(builder *flatbuffers.Builder, plugin flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(plugin), 0)
}
func UserAddPassword(builder *flatbuffers.Builder, password flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(password), 0)
}
func UserAddPasswordLastChanged(builder *flatbuffers.Builder, passwordLastChanged int64) {
	builder.PrependInt64Slot(5, passwordLastChanged, 0)
}
func UserAddLocked(builder *flatbuffers.Builder, locked bool) {
	builder.PrependBoolSlot(6, locked, false)
}
func UserAddAttributes(builder *flatbuffers.Builder, attributes flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(7, flatbuffers.UOffsetT(attributes), 0)
}
func UserAddIdentity(builder *flatbuffers.Builder, identity flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(8, flatbuffers.UOffsetT(identity), 0)
}
func UserEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type RoleEdge struct {
	_tab flatbuffers.Table
}

func GetRootAsRoleEdge(buf []byte, offset flatbuffers.UOffsetT) *RoleEdge {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RoleEdge{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsRoleEdge(buf []byte, offset flatbuffers.UOffsetT) *RoleEdge {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &RoleEdge{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *RoleEdge) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RoleEdge) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RoleEdge) FromHost() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RoleEdge) FromUser() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RoleEdge) ToHost() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RoleEdge) ToUser() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RoleEdge) WithAdminOption() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *RoleEdge) MutateWithAdminOption(n bool) bool {
	return rcv._tab.MutateBoolSlot(12, n)
}

func RoleEdgeStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func RoleEdgeAddFromHost(builder *flatbuffers.Builder, fromHost flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(fromHost), 0)
}
func RoleEdgeAddFromUser(builder *flatbuffers.Builder, fromUser flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(fromUser), 0)
}
func RoleEdgeAddToHost(builder *flatbuffers.Builder, toHost flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(toHost), 0)
}
func RoleEdgeAddToUser(builder *flatbuffers.Builder, toUser flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(toUser), 0)
}
func RoleEdgeAddWithAdminOption(builder *flatbuffers.Builder, withAdminOption bool) {
	builder.PrependBoolSlot(4, withAdminOption, false)
}
func RoleEdgeEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type MySQLDb struct {
	_tab flatbuffers.Table
}

func GetRootAsMySQLDb(buf []byte, offset flatbuffers.UOffsetT) *MySQLDb {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MySQLDb{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsMySQLDb(buf []byte, offset flatbuffers.UOffsetT) *MySQLDb {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MySQLDb{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *MySQLDb) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MySQLDb) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MySQLDb) User(obj *User, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *MySQLDb) UserLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MySQLDb) RoleEdges(obj *RoleEdge, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *MySQLDb) RoleEdgesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func MySQLDbStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func MySQLDbAddUser(builder *flatbuffers.Builder, user flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(user), 0)
}
func MySQLDbStartUserVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func MySQLDbAddRoleEdges(builder *flatbuffers.Builder, roleEdges flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(roleEdges), 0)
}
func MySQLDbStartRoleEdgesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func MySQLDbEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
