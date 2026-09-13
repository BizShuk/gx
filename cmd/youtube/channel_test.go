package youtube

import (
	"bytes"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

const testChannelID = "UCBR8-60-B28hp2BmDPdntcQ"

func setupTestServer(t *testing.T) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<link rel="canonical" href="https://www.youtube.com/channel/` + testChannelID + `">`))
	}))
	t.Cleanup(srv.Close)
	viper.Set("youtube_base_url", srv.URL)
	return srv
}

func executeCmd(args []string, in string) (string, error) {
	channelAsJSON = false
	channelIDOnly = false
	channelRSSOnly = false
	_ = channelCmd.Flags().Set("json", "false")
	_ = channelCmd.Flags().Set("id", "false")
	_ = channelCmd.Flags().Set("rss", "false")

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

func TestChannelCmd_SingleArg(t *testing.T) {
	srv := setupTestServer(t)

	got, err := executeCmd([]string{"@YouTube"}, "")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	want := srv.URL + "/channel/" + testChannelID + "\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestChannelCmd_SingleArgJSON(t *testing.T) {
	setupTestServer(t)

	got, err := executeCmd([]string{"@YouTube", "--json"}, "")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !strings.Contains(got, `"id": "`+testChannelID+`"`) {
		t.Errorf("JSON output does not contain expected channel ID: %s", got)
	}
	if !strings.Contains(got, `/feeds/videos.xml?channel_id=`+testChannelID) {
		t.Errorf("JSON output does not contain official RSS: %s", got)
	}
}

func TestChannelCmd_RSSFlag(t *testing.T) {
	srv := setupTestServer(t)

	got, err := executeCmd([]string{testChannelID, "--rss"}, "")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	want := srv.URL + "/feeds/videos.xml?channel_id=" + testChannelID + "\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestChannelCmd_MultipleArgs(t *testing.T) {
	srv := setupTestServer(t)

	got, err := executeCmd([]string{"@YouTube", "@NASA"}, "")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	expectedURL := srv.URL + "/channel/" + testChannelID
	want := expectedURL + "\n" + expectedURL + "\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestChannelCmd_MultipleArgsJSON(t *testing.T) {
	setupTestServer(t)

	got, err := executeCmd([]string{"@YouTube", "@NASA", "--json"}, "")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !strings.HasPrefix(strings.TrimSpace(got), "[") {
		t.Errorf("expected JSON array output, got: %s", got)
	}
}

func TestChannelCmd_Stdin(t *testing.T) {
	srv := setupTestServer(t)

	got, err := executeCmd([]string{}, "@YouTube\n@NASA\n")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	expectedURL := srv.URL + "/channel/" + testChannelID
	want := expectedURL + "\n" + expectedURL + "\n"
	if got != want {
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
