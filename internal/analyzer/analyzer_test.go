package analyzer

import "testing"

func TestAnalyzeSentiment(t *testing.T) {
	// TODO: Add tests with mock OpenAI client/responses
	t.Skip("Skipping analyzer test until actual analysis logic is implemented")
}

func TestStripHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple HTML", "<p>Hello, <b>world</b>!</p>", "Hello, world!"},
		{"No HTML", "Just plain text", "Just plain text"},
		{"Empty String", "", ""},
		{"Nested Tags", "<div><span>Text</span></div>", "Text"},
		{"Incomplete Tag", "<p>Incomplete", "Incomplete"},
		{"Tag Soup", "<p>Text <b>Bold</p> More", "Text Bold More"}, // Simple strip works ok here
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := stripHTML(tc.input)
			if actual != tc.expected {
				t.Errorf("stripHTML(%q) = %q; want %q", tc.input, actual, tc.expected)
			}
		})
	}
}
