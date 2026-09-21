package bilibili

import (
	"bytes"
	"strings"
	"testing"
)

const testSpaceID = "2267573"

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

const wantLines = "platform: bilibili\n" +
	"id: " + testSpaceID + "\n" +
	"url: https://space.bilibili.com/" + testSpaceID + "\n"

func TestChannelCmd_SingleArgLines(t *testing.T) {
	got, err := executeCmd([]string{"https://space.bilibili.com/" + testSpaceID + "/video"}, "")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if got != wantLines {
		t.Errorf("got %q, want %q", got, wantLines)
	}
}

func TestChannelCmd_SingleArgJSON(t *testing.T) {
	got, err := executeCmd([]string{testSpaceID, "--json"}, "")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !strings.HasPrefix(strings.TrimSpace(got), "[") {
		t.Errorf("expected JSON array output, got: %s", got)
	}
	if !strings.Contains(got, `"id": "`+testSpaceID+`"`) {
		t.Errorf("JSON output does not contain expected uid: %s", got)
	}
	if !strings.Contains(got, `"platform": "bilibili"`) {
		t.Errorf("JSON output does not contain platform: %s", got)
	}
	if strings.Contains(got, `"rss"`) {
		t.Errorf("JSON output must not invent an RSS field: %s", got)
	}
}

func TestChannelCmd_Stdin(t *testing.T) {
	got, err := executeCmd([]string{}, testSpaceID+"\n"+testSpaceID+"\n")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if want := wantLines + "\n" + wantLines; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestChannelCmd_EmptyInputError(t *testing.T) {
	_, err := executeCmd([]string{}, "\n   \n")
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestChannelCmd_RejectsVideoURL(t *testing.T) {
	_, err := executeCmd([]string{"https://www.bilibili.com/video/BV1hzqrBtEMP"}, "")
	if err == nil {
		t.Fatal("expected error for video url, got nil")
	}
}
