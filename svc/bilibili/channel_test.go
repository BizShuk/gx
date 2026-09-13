package bilibili

import (
	"context"
	"testing"
)

func TestGetChannelFromSpaceID(t *testing.T) {
	channel, err := NewClient().GetChannel(context.Background(), testSpaceID)
	if err != nil {
		t.Fatalf("GetChannel() returned error: %v", err)
	}
	if channel.ID != testSpaceID {
		t.Errorf("channel.ID = %q, want %q", channel.ID, testSpaceID)
	}
	want := DEFAULT_BASE_URL + "/" + testSpaceID
	if channel.URL != want {
		t.Errorf("channel.URL = %q, want %q", channel.URL, want)
	}
}

func TestGetChannelFromSpaceURL(t *testing.T) {
	input := "https://space.bilibili.com/" + testSpaceID + "/video"
	channel, err := NewClient().GetChannel(context.Background(), input)
	if err != nil {
		t.Fatalf("GetChannel() returned error: %v", err)
	}
	if channel.ID != testSpaceID {
		t.Errorf("channel.ID = %q, want %q", channel.ID, testSpaceID)
	}
}

func TestGetChannelUsesBaseURL(t *testing.T) {
	channel, err := NewClient(WithBaseURL("https://example.test/space/")).GetChannel(context.Background(), testSpaceID)
	if err != nil {
		t.Fatalf("GetChannel() returned error: %v", err)
	}
	want := "https://example.test/space/" + testSpaceID
	if channel.URL != want {
		t.Errorf("channel.URL = %q, want %q", channel.URL, want)
	}
}

func TestGetChannelRejectsVideoURL(t *testing.T) {
	if _, err := NewClient().GetChannel(context.Background(), "https://www.bilibili.com/video/BV1hzqrBtEMP"); err == nil {
		t.Fatal("GetChannel() = nil error, want error")
	}
}
