package applepodcast

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

const lookupBody = `{"resultCount":1,"results":[{"wrapperType":"track","kind":"podcast",
"collectionId":1702409419,"collectionName":"科技浪 Tech.wav",
"feedUrl":"https://feed.firstory.me/rss/user/abc"}]}`

func TestGetChannelLooksUpCollection(t *testing.T) {
	var gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Path + "?" + r.URL.RawQuery
		w.Write([]byte(lookupBody))
	}))
	defer srv.Close()

	ch, err := NewClient(WithBaseURL(srv.URL)).GetChannel(context.Background(),
		"https://podcasts.apple.com/tw/podcast/x/id"+testCollectionID+"?l=en-GB")
	if err != nil {
		t.Fatalf("GetChannel() returned error: %v", err)
	}
	if want := "/lookup?entity=podcast&id=" + testCollectionID; gotQuery != want {
		t.Errorf("requested = %q, want %q", gotQuery, want)
	}
	if ch.Platform != "apple-podcast" || ch.ID != testCollectionID || ch.Title != "科技浪 Tech.wav" {
		t.Errorf("channel = %+v", ch)
	}
	if want := "https://podcasts.apple.com/podcast/id" + testCollectionID; ch.URL != want {
		t.Errorf("channel.URL = %q, want %q", ch.URL, want)
	}
	if ch.RSS != "https://feed.firstory.me/rss/user/abc" {
		t.Errorf("channel.RSS = %q", ch.RSS)
	}
}

// lookup 對不存在的 ID 回 200 + 空結果，不是 404。
func TestGetChannelEmptyResultIsNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"resultCount":0,"results":[]}`))
	}))
	defer srv.Close()

	if _, err := NewClient(WithBaseURL(srv.URL)).GetChannel(context.Background(), "https://podcasts.apple.com/podcast/id"+testCollectionID); err == nil {
		t.Fatal("GetChannel() = nil error, want not found")
	}
}
