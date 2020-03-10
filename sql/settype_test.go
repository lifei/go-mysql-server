package sql

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetCompare(t *testing.T) {
	tests := []struct {
		vals []string
		collation Collation
		val1 interface{}
		val2 interface{}
		expectedCmp int
	}{
		{[]string{"one", "two"}, Collation_Default, nil, 1, -1},
		{[]string{"one", "two"}, Collation_Default, "one", nil, 1},
		{[]string{"one", "two"}, Collation_Default, nil, nil, 0},
		{[]string{"one", "two"}, Collation_Default, 0, "one", -1},
		{[]string{"one", "two"}, Collation_Default, 1, "two", -1},
		{[]string{"one", "two"}, Collation_Default, 2, []byte("one"), 1},
		{[]string{"one", "two"}, Collation_Default, "one", "", 1},
		{[]string{"one", "two"}, Collation_Default, "one", 1, 0},
		{[]string{"one", "two"}, Collation_Default, "one", "two", -1},
		{[]string{"two", "one"}, Collation_binary, "two", "one", -1},
		{[]string{"one", "two"}, Collation_Default, 3, "one,two", 0},
		{[]string{"one", "two"}, Collation_Default, "two,one,two", "one,two", 0},
		{[]string{"one", "two"}, Collation_Default, "two", "", 1},
		{[]string{"one", "two"}, Collation_Default, "one,two", "two", 1},
		{[]string{"a", "b", "c"}, Collation_Default, "a,b", "b,c", -1},
		{[]string{"a", "b", "c"}, Collation_Default, "a,b,c", "c,c,b", 1},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v %v %v %v", test.vals, test.collation, test.val1, test.val2), func(t *testing.T) {
			typ := MustCreateSetType(test.vals, test.collation)
			cmp, err := typ.Compare(test.val1, test.val2)
			require.NoError(t, err)
			assert.Equal(t, test.expectedCmp, cmp)
		})
	}
}

func TestSetCompareErrors(t *testing.T) {
	tests := []struct {
		vals []string
		collation Collation
		val1 interface{}
		val2 interface{}
	}{
		{[]string{"one", "two"}, Collation_Default, "three", "two"},
		{[]string{"one", "two"}, Collation_Default, time.Date(2019, 12, 12, 12, 12, 12, 0, time.UTC), []byte("one")},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v %v %v %v", test.vals, test.collation, test.val1, test.val2), func(t *testing.T) {
			typ := MustCreateSetType(test.vals, test.collation)
			_, err := typ.Compare(test.val1, test.val2)
			require.Error(t, err)
		})
	}
}

func TestSetCreate(t *testing.T) {
	tests := []struct {
		vals []string
		collation Collation
		expectedVals map[string]uint64
		expectedErr bool
	}{
		{[]string{"one"}, Collation_Default,
			map[string]uint64{"one": 1}, false},
		{[]string{" one ", "  two  "}, Collation_Default,
			map[string]uint64{" one": 1, "  two": 2}, false},
		{[]string{"a", "b", "c"}, Collation_Default,
			map[string]uint64{"a": 1, "b": 2, "c": 4}, false},
		{[]string{"one", "one "}, Collation_binary, map[string]uint64{"one": 1, "one ": 2}, false},
		{[]string{"one", "One"}, Collation_binary, map[string]uint64{"one": 1, "One": 2}, false},

		{[]string{}, Collation_Default, nil, true},
		{[]string{"one", "one"}, Collation_Default, nil, true},
		{[]string{"one", "one"}, Collation_binary, nil, true},
		{[]string{"one", "One"}, Collation_Default, nil, true},
		{[]string{"one", "one "}, Collation_Default, nil, true},
		{[]string{"one", "two,"}, Collation_Default, nil, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v %v", test.vals, test.collation), func(t *testing.T) {
			typ, err := CreateSetType(test.vals, test.collation)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				concreteType, ok := typ.(setType)
				require.True(t, ok)
				assert.Equal(t, test.collation, typ.Collation())
				for val, bit := range test.expectedVals {
					bitField, err := concreteType.convertStringToBitField(val)
					if assert.NoError(t, err) {
						assert.Equal(t, bit, bitField)
					}
					str, err := concreteType.convertBitFieldToString(bit)
					if assert.NoError(t, err) {
						assert.Equal(t, val, str)
					}
				}
			}
		})
	}
}

func TestSetCreateTooLarge(t *testing.T) {
	vals := make([]string, 65)
	for i := range vals {
		vals[i] = strconv.Itoa(i)
	}
	_, err := CreateSetType(vals, Collation_Default)
	require.Error(t, err)
}

func TestSetConvert(t *testing.T) {
	tests := []struct {
		vals []string
		collation Collation
		val interface{}
		expectedVal interface{}
		expectedErr bool
	}{
		{[]string{"one", "two"}, Collation_Default, nil, nil, false},
		{[]string{"one", "two"}, Collation_Default, "", "", false},
		{[]string{"one", "two"}, Collation_Default, int(0), "", false},
		{[]string{"one", "two"}, Collation_Default, int8(2), "two", false},
		{[]string{"one", "two"}, Collation_Default, int16(1), "one", false},
		{[]string{"one", "two"}, Collation_binary, int32(2), "two", false},
		{[]string{"one", "two"}, Collation_Default, int64(1), "one", false},
		{[]string{"one", "two"}, Collation_Default, uint(2), "two", false},
		{[]string{"one", "two"}, Collation_binary, uint8(1), "one", false},
		{[]string{"one", "two"}, Collation_Default, uint16(2), "two", false},
		{[]string{"one", "two"}, Collation_binary, uint32(3), "one,two", false},
		{[]string{"one", "two"}, Collation_Default, uint64(2), "two", false},
		{[]string{"one", "two"}, Collation_Default, "one", "one", false},
		{[]string{"one", "two"}, Collation_Default, []byte("two"), "two", false},
		{[]string{"one", "two"}, Collation_Default, "one,two", "one,two", false},
		{[]string{"one", "two"}, Collation_binary, "two,one", "one,two", false},
		{[]string{"one", "two"}, Collation_Default, "one,two,one", "one,two", false},
		{[]string{"one", "two"}, Collation_binary, "two,one,two", "one,two", false},
		{[]string{"one", "two"}, Collation_Default, "two,one,two", "one,two", false},
		{[]string{"a", "b", "c"}, Collation_Default, "b,c  ,a", "a,b,c", false},
		{[]string{"one", "two"}, Collation_Default, "ONE", "one", false},
		{[]string{"ONE", "two"}, Collation_Default, "one", "ONE", false},

		{[]string{"one", "two"}, Collation_Default, 4, nil, true},
		{[]string{"one", "two"}, Collation_Default, "three", nil, true},
		{[]string{"one", "two"}, Collation_Default, "one,two,three", nil, true},
		{[]string{"a", "b", "c"}, Collation_binary, "b,c  ,a", nil, true},
		{[]string{"one", "two"}, Collation_binary, "ONE", nil, true},
		{[]string{"one", "two"}, Collation_Default, time.Date(2019, 12, 12, 12, 12, 12, 0, time.UTC), nil, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v | %v | %v", test.vals, test.collation, test.val), func(t *testing.T) {
			typ := MustCreateSetType(test.vals, test.collation)
			val, err := typ.Convert(test.val)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedVal, val)
			}
		})
	}
}

func TestSetString(t *testing.T) {
	tests := []struct {
		vals []string
		collation Collation
		expectedStr string
	}{
		{[]string{"one"}, Collation_Default, "SET('one')"},
		{[]string{"مرحبا", "こんにちは"}, Collation_Default, "SET('مرحبا,こんにちは')"},
		{[]string{" hi ", "  lo  "}, Collation_Default, "SET(' hi,  lo')"},
		{[]string{" hi ", "  lo  "}, Collation_binary, "SET(' hi ,  lo  ') CHARACTER SET binary COLLATE binary"},
		{[]string{"a"}, Collation_Default.CharacterSet().BinaryCollation(),
			fmt.Sprintf("SET('a') COLLATE %v", Collation_Default.CharacterSet().BinaryCollation())},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v %v", test.vals, test.collation), func(t *testing.T) {
			typ := MustCreateSetType(test.vals, test.collation)
			assert.Equal(t, test.expectedStr, typ.String())
		})
	}
}