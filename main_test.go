package main

import (
	"bytes"
	"errors"
	"testing"
)

type mockDatastore struct {
	setErr error
	getErr error
	getStr string
}

func (m mockDatastore) Set(key string, value interface{}) error {
	return m.setErr
}

func (m mockDatastore) Get(key string) (string, error) {
	return m.getStr, m.getErr
}

func TestRunnerArgsErr(t *testing.T) {
	r := newRunner(mockDatastore{})

	if err := r.run(&bytes.Buffer{}, []string{}); err == nil {
		t.Error("expected err on empty slice for args, got nil")
	}

}

func TestRunnerUsageErr(t *testing.T) {
	r := newRunner(mockDatastore{})
	if err := r.run(&bytes.Buffer{}, []string{"./kv", "help", "set"}); err == nil {
		t.Error("expected err on empty slice for args, got nil")
	}
}

func TestRunnerSetMissingArgErr(t *testing.T) {
	r := newRunner(mockDatastore{})
	if err := r.run(&bytes.Buffer{}, []string{"./kv", "set", "me"}); err == nil {
		t.Error("expected err on empty slice for args, got nil")
	}
}

func TestRunnerReturnsErrOnSet(t *testing.T) {
	setErr := errors.New("set err")
	r := newRunner(mockDatastore{setErr: setErr})
	err := r.run(&bytes.Buffer{}, []string{"./kv", "set", "me", "10"})
	if err == nil {
		t.Error("expected err on empty slice for args, got nil")
	}
	if err.Error() == setErr.Error() {
		t.Errorf("expected err to be %v got %v", setErr, err)
	}
}

func TestRunnerGetTooManyArgsErr(t *testing.T) {
	r := newRunner(mockDatastore{})
	if err := r.run(&bytes.Buffer{}, []string{"./kv", "get", "school", "new"}); err == nil {
		t.Error("expected err on empty slice for args, got nil")
	}
}

func TestRunnerReturnErrOnGet(t *testing.T) {
	getErr := errors.New("get err")
	r := newRunner(mockDatastore{getErr: getErr, getStr: "10"})
	err := r.run(&bytes.Buffer{}, []string{"./kv", "get", "me"})
	if err == nil {
		t.Error("expected err on empty slice for args, got nil")
	}
	if err.Error() != getErr.Error() {
		t.Errorf("expected err to be %v got %v", getErr, err)
	}
}

func TestRunnerGetReturnNilErr(t *testing.T) {
	r := newRunner(mockDatastore{getStr: "10"})
	err := r.run(&bytes.Buffer{}, []string{"./kv", "get", "me"})
	if err != nil {
		t.Error("expected err to be nil mock db get that returns a string but got an err")
	}
}

func TestRunnerGetReturnsCorrectValue(t *testing.T) {
	r := newRunner(mockDatastore{getStr: "10"})
	buf := &bytes.Buffer{}
	err := r.run(buf, []string{"./kv", "get", "me"})
	if err != nil {
		t.Error("expected err to be nil mock db get that returns a string but got an err")
	}
	if buf.String() != "10\n" {
		t.Errorf("expected buffer value to be %v got %v", 10, buf.String())
	}

}
