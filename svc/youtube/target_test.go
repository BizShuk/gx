package youtube

import "testing"

func TestParseTargetHandle(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"bare name", "YouTube", "@YouTube"},
		{"at prefixed", "@YouTube", "@YouTube"},
		{"full url", "https://www.youtube.com/@YouTube", "@YouTube"},
		{"trailing slash", "https://www.youtube.com/@YouTube/", "@YouTube"},
		{"sub page", "https://www.youtube.com/@YouTube/videos", "@YouTube"},
		{"query string", "https://www.youtube.com/@YouTube?si=abc", "@YouTube"},
		{"no scheme", "youtube.com/@YouTube", "@YouTube"},
		{"mobile host", "https://m.youtube.com/@YouTube", "@YouTube"},
		{"legacy c path", "youtube.com/c/YouTube", "@YouTube"},
		{"legacy user path", "https://www.youtube.com/user/YouTube", "@YouTube"},
		{"padded input", "  @YouTube  ", "@YouTube"},
		{"dotted handle", "@some.name", "@some.name"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTarget(tc.input)
			if err != nil {
				t.Fatalf("ParseTarget(%q) returned error: %v", tc.input, err)
			}
			if got.Handle != tc.want {
				t.Errorf("ParseTarget(%q).Handle = %q, want %q", tc.input, got.Handle, tc.want)
			}
			if got.ID != "" {
				t.Errorf("ParseTarget(%q).ID = %q, want empty", tc.input, got.ID)
			}
		})
	}
}

func TestParseTargetChannelID(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"bare id", testChannelID},
		{"padded id", "  " + testChannelID + "  "},
		{"channel url", "https://www.youtube.com/channel/" + testChannelID},
		{"channel url sub page", "https://www.youtube.com/channel/" + testChannelID + "/videos"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTarget(tc.input)
			if err != nil {
				t.Fatalf("ParseTarget(%q) returned error: %v", tc.input, err)
			}
			if got.ID != testChannelID {
				t.Errorf("ParseTarget(%q).ID = %q, want %q", tc.input, got.ID, testChannelID)
			}
		})
	}
}

func TestParseTargetRejectsUnusableInput(t *testing.T) {
	inputs := []string{
		"",
		"   ",
		"/",
		"https://www.youtube.com/",
		"https://www.youtube.com",
		"@",
		"https://www.youtube.com/channel/not-an-id",
	}

	for _, input := range inputs {
		if got, err := ParseTarget(input); err == nil {
			t.Errorf("ParseTarget(%q) = %+v, nil error; want error", input, got)
		}
	}
}

func TestIsChannelID(t *testing.T) {
	cases := map[string]bool{
		testChannelID:       true,
		" " + testChannelID: true, // 前後空白會被修掉
		"@YouTube":          false,
		"UCshort":           false,
		"https://www.youtube.com/channel/" + testChannelID: false,
	}

	for input, want := range cases {
		if got := IsChannelID(input); got != want {
			t.Errorf("IsChannelID(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestExtractChannelID(t *testing.T) {
	cases := []struct {
		name string
		html string
	}{
		{"canonical link", `<link rel="canonical" href="https://www.youtube.com/channel/` + testChannelID + `">`},
		{"channelId field", `{"channelId":"` + testChannelID + `","title":"YouTube"}`},
		{"external id field", `{"externalId": "` + testChannelID + `"}`},
		{"plain channel url", `<a href="https://www.youtube.com/channel/` + testChannelID + `/videos">`},
		{"itemprop identifier", `<meta itemprop="identifier" content="` + testChannelID + `">`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := ExtractChannelID(tc.html)
			if !ok {
				t.Fatal("ExtractChannelID() found nothing, want a channel id")
			}
			if got != testChannelID {
				t.Errorf("ExtractChannelID() = %q, want %q", got, testChannelID)
			}
		})
	}
}

// 頁面裡與頻道 ID 同形狀的隨機 token 不該被誤認。
func TestExtractChannelIDIgnoresLookalikeTokens(t *testing.T) {
	html := `{"visitorData":"UCzzzzzzzzzzzzzzzzzzzzzz","channelId":"` + testChannelID + `"}`
	got, ok := ExtractChannelID(html)
	if !ok {
		t.Fatal("ExtractChannelID() found nothing, want a channel id")
	}
	if got != testChannelID {
		t.Errorf("ExtractChannelID() = %q, want %q", got, testChannelID)
	}

	if _, ok := ExtractChannelID(`{"visitorData":"UCzzzzzzzzzzzzzzzzzzzzzz"}`); ok {
		t.Error("ExtractChannelID() matched a token outside any channel context")
	}
}

func TestExtractChannelTitle(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
		ok   bool
	}{
		{"og title unescaped", `<meta property="og:title" content="A &amp; B">`, "A & B", true},
		{"metadata renderer json", `"channelMetadataRenderer":{"title":"財經 \"皓角\"","description":""}`, `財經 "皓角"`, true},
		{"og preferred", `<meta property="og:title" content="OG">"channelMetadataRenderer":{"title":"JSON"}`, "OG", true},
		{"blank og falls back", `<meta property="og:title" content=" ">"channelMetadataRenderer":{"title":"JSON"}`, "JSON", true},
		{"none", `<html></html>`, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ExtractChannelTitle(tt.html)
			if got != tt.want || ok != tt.ok {
				t.Errorf("ExtractChannelTitle() = (%q, %v), want (%q, %v)", got, ok, tt.want, tt.ok)
			}
		})
	}
}
