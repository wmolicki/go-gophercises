package parser

import (
	"os"
	"strings"
	"testing"
)

func eq(t *testing.T, a, b []Link) bool {
	t.Helper()

	if len(a) != len(b) {
		return false
	}

	for i, elem := range a {
		if b[i] != elem {
			return false
		}
	}

	return true
}

type TestCase struct {
	description string
	input       string
	want        []Link
}

func TestSimpleMarkup(t *testing.T) {

	cases := []TestCase{
		TestCase{
			"empty",
			"",
			[]Link{},
		},
		TestCase{
			"single node with markup",
			"<a href=\"wp.pl\">I am <strong>awesome</strong></a>",
			[]Link{Link{Href: "wp.pl", Text: "I am awesome"}},
		},
		TestCase{
			"text after markup",
			"<a href=\"onet.pl\">I am <strong>strong</strong> and awesome</a>",
			[]Link{Link{Href: "onet.pl", Text: "I am strong and awesome"}},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.description, func(t *testing.T) {
			reader := strings.NewReader(testCase.input)
			got := Parse(reader)
			want := testCase.want

			if !eq(t, want, got) {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}

type FileTestCase struct {
	description string
	filename    string
	want        []Link
}

func TestExampleFiles(t *testing.T) {

	cases := []FileTestCase{
		FileTestCase{
			"complex html",
			"resources/ex3.html",
			[]Link{
				Link{Text: "Login", Href: "#"},
				Link{Text: "Lost? Need help?", Href: "/lost"},
				Link{Text: "@marcusolsson", Href: "https://twitter.com/marcusolsson"},
			},
		},
		FileTestCase{
			"unnest markup two links simple html",
			"resources/ex2.html",
			[]Link{
				Link{Text: "Check me out on twitter", Href: "https://www.twitter.com/joncalhoun"},
				Link{Text: "Gophercises is on Github!", Href: "https://github.com/gophercises"},
			},
		},
		FileTestCase{
			"single link simple page",
			"resources/ex1.html",
			[]Link{Link{Text: "A link to another page", Href: "/other-page"}},
		},
		FileTestCase{
			"commented out text",
			"resources/ex4.html",
			[]Link{Link{Text: "dog cat after", Href: "/dog-cat"}},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.description, func(t *testing.T) {
			reader, err := os.Open(testCase.filename)
			if err != nil {
				t.Errorf("could not load test file %s: %v", testCase.filename, err)
			}

			got := Parse(reader)

			if !eq(t, testCase.want, got) {
				t.Errorf("want: %v, got: %v", testCase.want, got)
			}
		})
	}

}
