package youtube

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const testChannelID = "UCBR8-60-B28hp2BmDPdntcQ"

func TestGetChannelResolvesHandle(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Write([]byte(`<link rel="canonical" href="https://www.youtube.com/channel/` + testChannelID + `">`))
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	channel, err := client.GetChannel(context.Background(), "@YouTube")
	if err != nil {
		t.Fatalf("GetChannel() returned error: %v", err)
	}

	if gotPath != "/@YouTube" {
		t.Errorf("requested path = %q, want %q", gotPath, "/@YouTube")
	}
	if channel.ID != testChannelID {
		t.Errorf("channel.ID = %q, want %q", channel.ID, testChannelID)
	}
	if want := srv.URL + "/channel/" + testChannelID; channel.URL != want {
		t.Errorf("channel.URL = %q, want %q", channel.URL, want)
	}
}

// 輸入已經是頻道 ID 時不該發出任何請求。
func TestGetChannelSkipsFetchForChannelID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request to %s", r.URL.Path)
	}))
	defer srv.Close()

	channel, err := NewClient(WithBaseURL(srv.URL)).GetChannel(context.Background(), testChannelID)
	if err != nil {
		t.Fatalf("GetChannel() returned error: %v", err)
	}
	if channel.ID != testChannelID {
		t.Errorf("channel.ID = %q, want %q", channel.ID, testChannelID)
	}
}

func TestGetChannelErrorsWhenPageHasNoChannelID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>nothing here</body></html>"))
	}))
	defer srv.Close()

	if _, err := NewClient(WithBaseURL(srv.URL)).GetChannel(context.Background(), "@ghost"); err == nil {
		t.Fatal("GetChannel() = nil error, want error")
	}
}

// 404 是呼叫端的錯，不該消耗重試預算。
func TestGetChannelDoesNotRetryNotFound(t *testing.T) {
	var attempts int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := NewClient(WithBaseURL(srv.URL)).GetChannel(context.Background(), "@missing")
	if err == nil {
		t.Fatal("GetChannel() = nil error, want error")
	}
	if want := "channel @missing not found"; err.Error() != want {
		t.Errorf("error = %q, want %q", err, want)
	}
	if attempts != 1 {
		t.Errorf("attempts = %d, want 1", attempts)
	}
}

// 5xx 是暫時性失敗，重試後應該成功。
func TestGetChannelRetriesServerError(t *testing.T) {
	var attempts int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte(`href="https://www.youtube.com/channel/` + testChannelID + `"`))
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	channel, err := NewClient(WithBaseURL(srv.URL)).GetChannel(ctx, "@flaky")
	if err != nil {
		t.Fatalf("GetChannel() returned error: %v", err)
	}
	if attempts != 2 {
		t.Errorf("attempts = %d, want 2", attempts)
	}
	if channel.ID != testChannelID {
		t.Errorf("channel.ID = %q, want %q", channel.ID, testChannelID)
	}
}
