package curator

import (
	"errors"
	"testing"
)

var errTestDummy = errors.New("test dummy")

func TestDummyConn_Children(t *testing.T) {
	dummy := dummyConn{err: errTestDummy}
	_, _, err := dummy.Children("")
	if err != errTestDummy {
		t.Fatal("unexpected error")
	}
}

func TestDummyConn_ChildrenW(t *testing.T) {
	dummy := dummyConn{err: errTestDummy}
	_, _, _, err := dummy.ChildrenW("")
	if err != errTestDummy {
		t.Fatal("unexpected error")
	}
}

func TestDummyConn_Create(t *testing.T) {
	dummy := dummyConn{err: errTestDummy}
	_, err := dummy.Create("", nil, 0, nil)
	if err != errTestDummy {
		t.Fatal("unexpected error")
	}
}

func TestDummyConn_CreateProtectedEphemeralSequential(t *testing.T) {
	dummy := dummyConn{err: errTestDummy}
	_, err := dummy.CreateProtectedEphemeralSequential("", nil, nil)
	if err != errTestDummy {
		t.Fatal("unexpected error")
	}
}

func TestDummyConn_Delete(t *testing.T) {
	dummy := dummyConn{err: errTestDummy}
	err := dummy.Delete("", 0)
	if err != errTestDummy {
		t.Fatal("unexpected error")
	}
}

func TestDummyConn_Exists(t *testing.T) {
	dummy := dummyConn{err: errTestDummy}
	_, _, err := dummy.Exists("")
	if err != errTestDummy {
		t.Fatal("unexpected error")
	}
}

func TestDummyConn_ExistsW(t *testing.T) {
	dummy := dummyConn{err: errTestDummy}
	_, _, _, err := dummy.ExistsW("")
	if err != errTestDummy {
		t.Fatal("unexpected error")
	}
}

func TestDummyConn_Get(t *testing.T) {
	dummy := dummyConn{err: errTestDummy}
	_, _, err := dummy.Get("")
	if err != errTestDummy {
		t.Fatal("unexpected error")
	}
}

func TestDummyConn_GetACL(t *testing.T) {
	dummy := dummyConn{err: errTestDummy}
	_, _, err := dummy.GetACL("")
	if err != errTestDummy {
		t.Fatal("unexpected error")
	}
}

func TestDummyConn_GetW(t *testing.T) {
	dummy := dummyConn{err: errTestDummy}
	_, _, _, err := dummy.GetW("")
	if err != errTestDummy {
		t.Fatal("unexpected error")
	}
}

func TestDummyConn_Set(t *testing.T) {
	dummy := dummyConn{err: errTestDummy}
	_, err := dummy.Set("", nil, 0)
	if err != errTestDummy {
		t.Fatal("unexpected error")
	}
}

func TestDummyConn_SetACL(t *testing.T) {
	dummy := dummyConn{err: errTestDummy}
	_, err := dummy.SetACL("", nil, 0)
	if err != errTestDummy {
		t.Fatal("unexpected error")
	}
}
