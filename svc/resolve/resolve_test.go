package resolve

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bizshuk/gx/model"
	"github.com/bizshuk/gx/svc/applepodcast"
	"github.com/bizshuk/gx/svc/youtube"
)

func TestPlatformByURLDomain(t *testing.T) {
	cases := map[string]string{
		"https://www.youtube.com/@YouTube":                             model.PLATFORM_YOUTUBE,
		"youtube.com/channel/UCBR8-60-B28hp2BmDPdntcQ":                 model.PLATFORM_YOUTUBE,
		"https://m.youtube.com/@YouTube":                               model.PLATFORM_YOUTUBE,
		"https://space.bilibili.com/2267573/video":                     model.PLATFORM_BILIBILI,
		"https://podcasts.apple.com/tw/podcast/x/id1702409419?l=en-GB": model.PLATFORM_APPLE_PODCAST,
		"https://itunes.apple.com/us/podcast/id1702409419":             model.PLATFORM_APPLE_PODCAST,
		// 不是網址：說不出平台。
		"@YouTube":                 "",
		"UCBR8-60-B28hp2BmDPdntcQ": "",
		"1702409419":               "",
		"id1702409419":             "",
		"":                         "",
		// 網域不認得。
		"https://open.spotify.com/show/abc": "",
		"https://example.com/@YouTube":      "",
	}
	for input, want := range cases {
		if got := Platform(input); got != want {
			t.Errorf("Platform(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestChannelDispatchesByDomain(t *testing.T) {
	yt := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<link rel="canonical" href="https://www.youtube.com/channel/UCBR8-60-B28hp2BmDPdntcQ">`))
	}))
	defer yt.Close()
	ap := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"resultCount":1,"results":[{"kind":"podcast","collectionId":1702409419,
"collectionName":"科技浪 Tech.wav","feedUrl":"https://feed.example/rss"}]}`))
	}))
	defer ap.Close()

	r := New(
		WithYouTube(youtube.NewClient(youtube.WithBaseURL(yt.URL))),
		WithApplePodcast(applepodcast.NewClient(applepodcast.WithBaseURL(ap.URL))),
	)
	ctx := context.Background()

	ch, err := r.Channel(ctx, "https://www.youtube.com/@YouTube")
	if err != nil || ch.Platform != model.PLATFORM_YOUTUBE || ch.ID != "UCBR8-60-B28hp2BmDPdntcQ" {
		t.Errorf("youtube: ch=%+v err=%v", ch, err)
	}
	ch, err = r.Channel(ctx, "https://podcasts.apple.com/tw/podcast/x/id1702409419")
	if err != nil || ch.Platform != model.PLATFORM_APPLE_PODCAST || ch.RSS != "https://feed.example/rss" {
		t.Errorf("apple-podcast: ch=%+v err=%v", ch, err)
	}
	ch, err = r.Channel(ctx, "https://space.bilibili.com/2267573")
	if err != nil || ch.Platform != model.PLATFORM_BILIBILI || ch.ID != "2267573" {
		t.Errorf("bilibili: ch=%+v err=%v", ch, err)
	}
	if _, err := r.Channel(ctx, "@YouTube"); !errors.Is(err, ErrUnknownURL) {
		t.Errorf("bare handle: err=%v, want ErrUnknownURL", err)
	}
}
