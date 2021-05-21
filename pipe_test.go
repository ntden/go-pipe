package pipe

import (
	"errors"
	"reflect"
	"testing"
)

func TestPipe_Execute(t *testing.T) {
	fn1 := func(a int, b float64) (int, error) {
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / int(b), nil
	}
	fn2 := func(a int) string {
		if a%2 == 0 {
			return "even"
		}
		return "odd"
	}
	fn3 := func(s string) []byte {
		return []byte(s)
	}

	// Create a new pipe.
	p, err := New(fn1, fn2)
	if err != nil {
		t.Fatalf("unexpected error creating a new pipe: %v", err)
	}

	err = p.Add(fn3)
	if err != nil {
		t.Fatalf("unexpected error adding functions to pipe: %v", err)
	}

	tests := []struct {
		inputArgs      []interface{}
		expectedOutput []interface{}
		expectError    bool
	}{
		{
			inputArgs:      []interface{}{10, 3.0},
			expectedOutput: []interface{}{[]byte("odd")},
			expectError:    false,
		},
		{
			inputArgs:      []interface{}{10},
			expectedOutput: nil,
			expectError:    true,
		},
		{
			inputArgs:      []interface{}{10, 3.0, 1},
			expectedOutput: []interface{}{[]byte("odd")},
			expectError:    false,
		},
		{
			inputArgs:      []interface{}{3.22},
			expectedOutput: nil,
			expectError:    true,
		},
	}

	for i, test := range tests {
		output, err := p.Execute(test.inputArgs...)
		if test.expectError && err == nil {
			t.Errorf("test %d: expected an error but got nil", i)
			continue
		}
		if !test.expectError && err != nil {
			t.Errorf("test %d: unexpected error: %v", i, err)
			continue
		}
		if !reflect.DeepEqual(output, test.expectedOutput) {
			t.Errorf("test %d: output mismatch: expected %v, got %v", i, test.expectedOutput, output)
		}
	}
}
