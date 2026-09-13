package bilibili

import "testing"

const testSpaceID = "2267573"

func TestParseTargetSpaceID(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"bare id", testSpaceID},
		{"padded id", "  " + testSpaceID + "  "},
		{"space url", "https://space.bilibili.com/" + testSpaceID},
		{"trailing slash", "https://space.bilibili.com/" + testSpaceID + "/"},
		{"sub page", "https://space.bilibili.com/" + testSpaceID + "/video"},
		{"collection page", "https://space.bilibili.com/" + testSpaceID + "/lists/598034?type=season"},
		{"query string", "https://space.bilibili.com/" + testSpaceID + "?spm_id_from=abc"},
		{"fragment", "https://space.bilibili.com/" + testSpaceID + "#/"},
		{"no scheme", "space.bilibili.com/" + testSpaceID},
		{"mobile space path", "https://m.bilibili.com/space/" + testSpaceID},
		{"www space path", "https://www.bilibili.com/space/" + testSpaceID},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTarget(tc.input)
			if err != nil {
				t.Fatalf("ParseTarget(%q) returned error: %v", tc.input, err)
			}
			if got.ID != testSpaceID {
				t.Errorf("ParseTarget(%q).ID = %q, want %q", tc.input, got.ID, testSpaceID)
			}
		})
	}
}

func TestParseTargetRejectsUnusableInput(t *testing.T) {
	inputs := []string{
		"",
		"   ",
		"/",
		"https://space.bilibili.com/",
		"https://www.bilibili.com/",
		"https://www.bilibili.com/video/BV1hzqrBtEMP",
		"BV1hzqrBtEMP",
		"DIYgod",
		"@DIYgod",
		"https://b23.tv/abcdef",
		"0",
		"0123",
	}

	for _, input := range inputs {
		if got, err := ParseTarget(input); err == nil {
			t.Errorf("ParseTarget(%q) = %+v, nil error; want error", input, got)
		}
	}
}

func TestIsSpaceID(t *testing.T) {
	cases := map[string]bool{
		testSpaceID:       true,
		" " + testSpaceID: true,
		"2":               true,
		"0":               false,
		"0123":            false,
		"DIYgod":          false,
		"https://space.bilibili.com/" + testSpaceID: false,
	}

	for input, want := range cases {
		if got := IsSpaceID(input); got != want {
			t.Errorf("IsSpaceID(%q) = %v, want %v", input, got, want)
		}
	}
}
