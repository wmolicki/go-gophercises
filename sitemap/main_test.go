package main

import "testing"

type TestCase struct {
	Description string
	A           string
	B           string
	Result      bool
}


type NormalizeUrlTestCase struct {
	url string
	base string
	normalized string
}


func TestNormalizeUrl(t *testing.T) {
	testCases := []NormalizeUrlTestCase{
		NormalizeUrlTestCase{"/", "https://wp.pl", "https://wp.pl/"},
		NormalizeUrlTestCase{"/abc", "https://wp.pl", "https://wp.pl/abc"},
		NormalizeUrlTestCase{"/abc", "https://wp.pl/", "https://wp.pl/abc"},
		NormalizeUrlTestCase{"/abc#fragment", "https://wp.pl/", "https://wp.pl/abc"},
		NormalizeUrlTestCase{"/abc?q=s", "https://wp.pl/", "https://wp.pl/abc"},
		NormalizeUrlTestCase{"https://onet.pl/abc?q=s", "https://onet.pl", "https://onet.pl/abc"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.url, func(t *testing.T) {
			got := normalizeUrl(testCase.base, testCase.url)
			want := testCase.normalized

			if got != want {
				t.Fatalf("got %v, wanted %v", got, want)
			}
		})
	}
}

func TestSameDomain(t *testing.T) {

	testCases := []TestCase{
		TestCase{"path links", "/test", "/xxx/yyy", true},
		TestCase{"same domain", "https://wp.pl", "https://wp.pl/pogoda", true},
		TestCase{"different domain", "https://wp.pl", "https://onet.pl", false},
		TestCase{"no scheme", "wp.pl", "https://wp.pl", true},
		TestCase{"path link and domain", "https://wp.pl", "/abc", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			got := sameDomain(testCase.A, testCase.B)
			want := testCase.Result

			if got != want {
				t.Fatalf("got %v, wanted %v", got, want)
			}
		})
	}
}
