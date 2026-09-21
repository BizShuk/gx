package applepodcast

import "testing"

const testCollectionID = "1702409419"

func TestParseTargetCollectionID(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"bare id", testCollectionID},
		{"padded id", "  " + testCollectionID + "  "},
		{"id prefix", "id" + testCollectionID},
		{"show url", "https://podcasts.apple.com/tw/podcast/%E7%A7%91%E6%8A%80%E6%B5%AA-tech-wav/id" + testCollectionID + "?l=en-GB"},
		{"no region", "https://podcasts.apple.com/podcast/id" + testCollectionID},
		{"episode url", "https://podcasts.apple.com/us/podcast/x/id" + testCollectionID + "?i=1000712345678"},
		{"itunes host", "https://itunes.apple.com/us/podcast/id" + testCollectionID},
		{"no scheme", "podcasts.apple.com/us/podcast/x/id" + testCollectionID},
		{"id query", "https://itunes.apple.com/lookup?id=" + testCollectionID},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTarget(tc.input)
			if err != nil {
				t.Fatalf("ParseTarget(%q) returned error: %v", tc.input, err)
			}
			if got.ID != testCollectionID {
				t.Errorf("ParseTarget(%q).ID = %q, want %q", tc.input, got.ID, testCollectionID)
			}
		})
	}
}

func TestParseTargetRejectsUnusableInput(t *testing.T) {
	inputs := []string{
		"",
		"   ",
		"0123",
		"@YouTube",
		"https://www.youtube.com/@YouTube",
		"https://podcasts.apple.com/us/podcast/no-id",
		"https://example.com/podcast/id" + testCollectionID,
		"https://open.spotify.com/show/abc",
	}
	for _, input := range inputs {
		if got, err := ParseTarget(input); err == nil {
			t.Errorf("ParseTarget(%q) = %+v, nil error; want error", input, got)
		}
	}
}
