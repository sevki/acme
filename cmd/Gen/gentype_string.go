// Code generated by "stringer -type genType cmd/Gen/tests.go"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[tests-0]
	_ = x[mocks-1]
}

const _genType_name = "testsmocks"

var _genType_index = [...]uint8{0, 5, 10}

func (i genType) String() string {
	if i < 0 || i >= genType(len(_genType_index)-1) {
		return "genType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _genType_name[_genType_index[i]:_genType_index[i+1]]
}
