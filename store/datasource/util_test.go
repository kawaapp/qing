package datasource

import "testing"

func TestJoinIntArray(t *testing.T) {
	if get := joinIntArray([]int64{}); "" != get {
		t.Error("int array should be empty, get:", get)
	}

	if get := joinIntArray([]int64{1}); "1" != get {
		t.Error("int array should be 1, get:", get)
	}

	if get := joinIntArray([]int64{1, 2}); "1,2" != get {
		t.Error("int array should be 1,2 get:", get)
	}
}
