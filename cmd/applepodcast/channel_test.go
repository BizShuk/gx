package applepodcast

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	svc "github.com/bizshuk/gx/svc/applepodcast"
	"github.com/spf13/viper"
)

const testCollectionID = "1702409419"

func executeCmd(t *testing.T, args []string, in string) (string, error) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"resultCount":1,"results":[{"kind":"podcast","collectionId":1702409419,
"collectionName":"科技浪 Tech.wav","feedUrl":"https://feed.example/rss"}]}`))
	}))
	t.Cleanup(srv.Close)
	viper.Set(svc.KEY_BASE_URL, srv.URL)
	t.Cleanup(func() { viper.Set(svc.KEY_BASE_URL, "") })

	channelAsJSON = false
	_ = channelCmd.Flags().Set("json", "false")

	var out bytes.Buffer
	Cmd.SetOut(&out)
	Cmd.SetErr(&out)
	Cmd.SetIn(strings.NewReader(in))
	Cmd.SetArgs(append([]string{"get", "channel"}, args...))
	err := Cmd.Execute()
	return out.String(), err
}

const wantLines = "platform: apple-podcast\n" +
	"id: " + testCollectionID + "\n" +
	"title: 科技浪 Tech.wav\n" +
	"url: https://podcasts.apple.com/podcast/id" + testCollectionID + "\n" +
	"rss: https://feed.example/rss\n"

func TestChannelCmd_SingleArgLines(t *testing.T) {
	got, err := executeCmd(t, []string{"https://podcasts.apple.com/tw/podcast/x/id" + testCollectionID + "?l=en-GB"}, "")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if got != wantLines {
		t.Errorf("got %q, want %q", got, wantLines)
	}
}

func TestChannelCmd_JSONIsArray(t *testing.T) {
	got, err := executeCmd(t, []string{testCollectionID, "--json"}, "")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(got), "[") || !strings.Contains(got, `"platform": "apple-podcast"`) {
		t.Errorf("unexpected JSON output: %s", got)
	}
}

func TestChannelCmd_Stdin(t *testing.T) {
	got, err := executeCmd(t, nil, testCollectionID+"\nid"+testCollectionID+"\n")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if want := wantLines + "\n" + wantLines; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestChannelCmd_RejectsNonAppleURL(t *testing.T) {
	if _, err := executeCmd(t, []string{"https://open.spotify.com/show/abc"}, ""); err == nil {
		t.Fatal("expected error for non-Apple url, got nil")
	}
}
