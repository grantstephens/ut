package ut

import (
	"io"
	"testing"
)

type MockReader struct {
	CallTracker
}

func NewMockReader(t *testing.T) *MockReader {
	return &MockReader{NewCallRecords(t)}
}

func (m *MockReader) Read(p []byte) (n int, err error) {
	r := m.TrackCall("Read", p)
	return r[0].(int), NilOrError(r[1])
}

func UnderTest(r io.Reader) bool {
	p := make([]byte, 10)
	n, _ := r.Read(p)

	return n >= 1 && p[0] == 37
}

func TestUnderTest(t *testing.T) {

	// Define the tests we're going to run.
	tests := []struct {
		bytezero byte
		n        int
		expRet   bool
	}{
		{bytezero: 37, n: 1, expRet: true},
		{bytezero: 37, n: 2, expRet: true},
		{bytezero: 38, n: 2, expRet: false},
		{bytezero: 0, n: 2, expRet: false},
		{bytezero: 37, n: 0, expRet: false},
		{bytezero: 0, n: 0, expRet: false},
	}

	for _, test := range tests {
		// Set up the mock
		m := NewMockReader(t)

		// Parameters for AddCall can either be: values, which are compared against the actual parameter;
		// or functions, which can check and act on the parameter as they like
		checkReadParam := func(p interface{}) {
			buf := p.([]byte)
			if len(buf) != 10 {
				t.Fatalf("should have read 10 bytes")
			}
			buf[0] = test.bytezero
		}

		// Note the calls we expect to happen when we run our test
		m.AddCall("Read", checkReadParam).SetReturns(test.n, error(nil))

		// Test the function
		if UnderTest(m) != test.expRet {
			t.Fatalf("return not as expected")
		}

		// Check the method calls we expected actually happened
		m.AssertDone()
	}
}
