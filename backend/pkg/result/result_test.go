package result

import (
	"errors"
	"testing"
)

func TestSuccess(t *testing.T) {
	res := Success(100)
	if !res.IsSuccess() {
		t.Error("Expected Success to be IsSuccess")
	}
	if res.Value != 100 {
		t.Errorf("Expected 100, got %v", res.Value)
	}
}

func TestFailure(t *testing.T) {
	err := errors.New("test error")
	res := Failure[int](err)
	if !res.IsFailure() {
		t.Error("Expected Failure to be IsFailure")
	}
	if res.Error != err {
		t.Errorf("Expected %v, got %v", err, res.Error)
	}
}

func TestUnwrap(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected Unwrap to panic on failure")
		}
	}()

	res := Failure[int](errors.New("panic!"))
	res.Unwrap()
}

func TestErrs(t *testing.T) {
	res := Errs[int]("test error")
	if !res.IsFailure() {
		t.Error("Expected Errs to be IsFailure")
	}
	if res.Error.Error() != "test error" {
		t.Errorf("Expected 'test error', got %v", res.Error)
	}
}
