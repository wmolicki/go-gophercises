package main

import (
	"io"
	"strings"
	"testing"
)

type TestCsv struct {
	content  string
	expected string
}

type StubCSVReader struct {
	testCsv TestCsv
	read    bool
}

func (s *StubCSVReader) Read(p []byte) (n int, err error) {
	if s.read {
		return 0, io.EOF
	}
	i := 0
	line := s.testCsv.content + "," + s.testCsv.expected
	for ; i < len(p) && i < len(line); i++ {
		p[i] = byte(line[i])
	}
	s.read = true

	return i, nil
}

func TestProblemsLoaderLoadsSingleLine(t *testing.T) {
	cases := []TestCsv{
		TestCsv{"1+1", "2"},
		TestCsv{"3+4", "7"},
	}

	for _, testCase := range cases {
		t.Run("a test", func(t *testing.T) {

			//stubReader := &StubCSVReader{testCase, false}
			problemsLoader := NewReaderProblemsLoader(strings.NewReader(testCase.content + "," + testCase.expected))

			problems := problemsLoader.Load()

			got := problems[0]
			want := Problem{testCase.content, testCase.expected}

			if got != want {
				t.Errorf("got %v, wanted %v", got, want)
			}
		})
	}

}
