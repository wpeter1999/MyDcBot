package player

import (
	"testing"
)

func TestSong(t *testing.T) {
	t.Run("Song_结构定义", func(t *testing.T) {
		song := Song{
			Title:       "测试歌曲",
			URL:         "https://youtube.com/watch?v=test",
			StreamURL:   "https://stream.example.com/audio.m4a",
			RequestedBy: "user123",
		}

		if song.Title == "" {
			t.Error("Song.Title should not be empty")
		}

		if song.URL == "" {
			t.Error("Song.URL should not be empty")
		}

		if song.RequestedBy == "" {
			t.Error("Song.RequestedBy should not be empty")
		}
	})

	t.Run("Song_零值", func(t *testing.T) {
		var song Song

		if song.Title != "" {
			t.Error("Zero value Song.Title should be empty")
		}

		if song.URL != "" {
			t.Error("Zero value Song.URL should be empty")
		}

		if song.StreamURL != "" {
			t.Error("Zero value Song.StreamURL should be empty")
		}

		if song.RequestedBy != "" {
			t.Error("Zero value Song.RequestedBy should be empty")
		}
	})

	t.Run("Song_复制", func(t *testing.T) {
		original := Song{
			Title:       "原始歌曲",
			URL:         "https://youtube.com/watch?v=original",
			StreamURL:   "https://stream.example.com/original.m4a",
			RequestedBy: "user123",
		}

		copied := original

		if copied.Title != original.Title {
			t.Error("Copied song should have same Title")
		}

		if copied.URL != original.URL {
			t.Error("Copied song should have same URL")
		}

		if copied.StreamURL != original.StreamURL {
			t.Error("Copied song should have same StreamURL")
		}

		if copied.RequestedBy != original.RequestedBy {
			t.Error("Copied song should have same RequestedBy")
		}

		// 修改复制的歌曲不应影响原始歌曲
		copied.Title = "修改后的歌曲"
		if original.Title == copied.Title {
			t.Error("Modifying copied song should not affect original")
		}
	})
}
