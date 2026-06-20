package command

import (
	"testing"
)

func TestYtDlpInfo(t *testing.T) {
	t.Run("YtDlpInfo_結構定義", func(t *testing.T) {
		info := YtDlpInfo{
			URL:   "https://youtube.com/watch?v=test",
			Title: "Test Video",
		}

		if info.URL == "" {
			t.Error("YtDlpInfo.URL should not be empty")
		}

		if info.Title == "" {
			t.Error("YtDlpInfo.Title should not be empty")
		}
	})

	t.Run("YtDlpInfo_Formats字段", func(t *testing.T) {
		info := YtDlpInfo{
			Formats: []struct {
				URL      string `json:"url"`
				AudioExt string `json:"aext"`
				VideoExt string `json:"vext"`
			}{
				{
					URL:      "https://example.com/audio.m4a",
					AudioExt: "m4a",
					VideoExt: "none",
				},
			},
		}

		if len(info.Formats) != 1 {
			t.Errorf("YtDlpInfo.Formats length = %d, want 1", len(info.Formats))
		}

		if info.Formats[0].AudioExt != "m4a" {
			t.Errorf("Format.AudioExt = %s, want 'm4a'", info.Formats[0].AudioExt)
		}
	})
}
