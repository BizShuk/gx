package youtube

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bizshuk/gx/model"
	"github.com/spf13/viper"
)

const (
	testChannelID = "UCBR8-60-B28hp2BmDPdntcQ"
	testTitle     = "YouTube &amp; Friends"
)

func setupTestServer(t *testing.T) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<link rel="canonical" href="https://www.youtube.com/channel/` + testChannelID + `">` +
			`<meta property="og:title" content="` + testTitle + `">`))
	}))
	t.Cleanup(srv.Close)
	viper.Set("youtube_base_url", srv.URL)
	return srv
}

func executeCmd(args []string, in string) (string, error) {
	channelAsJSON = false
	_ = channelCmd.Flags().Set("json", "false")

	var out bytes.Buffer
	Cmd.SetOut(&out)
	Cmd.SetErr(&out)
	if in != "" {
		Cmd.SetIn(strings.NewReader(in))
	}
	Cmd.SetArgs(append([]string{"get", "channel"}, args...))

	err := Cmd.Execute()
	return out.String(), err
}

func wantLines(base string) string {
	return "platform: youtube\n" +
		"id: " + testChannelID + "\n" +
		"handle: @YouTube\n" +
		"title: YouTube & Friends\n" +
		"url: " + base + "/channel/" + testChannelID + "\n" +
		"rss: " + base + "/feeds/videos.xml?channel_id=" + testChannelID + "\n"
}

func TestChannelCmd_SingleArgLines(t *testing.T) {
	srv := setupTestServer(t)

	got, err := executeCmd([]string{"@YouTube"}, "")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if want := wantLines(srv.URL); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// --json 單筆也輸出陣列：呼叫端不必依目標數量切換解析方式。
func TestChannelCmd_SingleArgJSONIsArray(t *testing.T) {
	srv := setupTestServer(t)

	got, err := executeCmd([]string{"@YouTube", "--json"}, "")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	var channels []model.Channel
	if err := json.Unmarshal([]byte(got), &channels); err != nil {
		t.Fatalf("output is not a JSON array: %v\n%s", err, got)
	}
	want := model.Channel{
		Platform: model.PLATFORM_YOUTUBE,
		ID:       testChannelID,
		Handle:   "@YouTube",
		Title:    "YouTube & Friends",
		URL:      srv.URL + "/channel/" + testChannelID,
		RSS:      srv.URL + "/feeds/videos.xml?channel_id=" + testChannelID,
	}
	if len(channels) != 1 || channels[0] != want {
		t.Errorf("got %+v, want [%+v]", channels, want)
	}
}

func TestChannelCmd_MultipleArgsLines(t *testing.T) {
	srv := setupTestServer(t)

	got, err := executeCmd([]string{"@YouTube", "@YouTube"}, "")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if want := wantLines(srv.URL) + "\n" + wantLines(srv.URL); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestChannelCmd_Stdin(t *testing.T) {
	srv := setupTestServer(t)

	got, err := executeCmd([]string{}, "@YouTube\n@YouTube\n")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if want := wantLines(srv.URL) + "\n" + wantLines(srv.URL); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestChannelCmd_EmptyInputError(t *testing.T) {
	setupTestServer(t)

	_, err := executeCmd([]string{}, "\n   \n")
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}
