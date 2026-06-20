package command

import (
	"testing"
)

func TestYtDlpInfo_Structure(t *testing.T) {
	t.Run("YtDlpInfo_完整結構", func(t *testing.T) {
		info := YtDlpInfo{
			URL:   "https://youtube.com/watch?v=test",
			Title: "Test Video",
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
				{
					URL:      "https://example.com/video.mp4",
					AudioExt: "none",
					VideoExt: "mp4",
				},
			},
		}

		if info.URL == "" {
			t.Error("YtDlpInfo.URL should not be empty")
		}

		if info.Title == "" {
			t.Error("YtDlpInfo.Title should not be empty")
		}

		if len(info.Formats) != 2 {
			t.Errorf("YtDlpInfo.Formats length = %d, want 2", len(info.Formats))
		}

		// 驗證音頻格式
		if info.Formats[0].AudioExt != "m4a" {
			t.Errorf("Format[0].AudioExt = %s, want 'm4a'", info.Formats[0].AudioExt)
		}

		if info.Formats[0].VideoExt != "none" {
			t.Errorf("Format[0].VideoExt = %s, want 'none'", info.Formats[0].VideoExt)
		}
	})

	t.Run("YtDlpInfo_零值", func(t *testing.T) {
		var info YtDlpInfo

		if info.URL != "" {
			t.Error("Zero value URL should be empty")
		}

		if info.Title != "" {
			t.Error("Zero value Title should be empty")
		}

		if len(info.Formats) != 0 {
			t.Error("Zero value Formats should be empty slice")
		}
	})
}
