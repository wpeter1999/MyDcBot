package command

import (
	"testing"
)

func TestPlaylistEntry_Structure(t *testing.T) {
	t.Run("PlaylistEntry_基本結構", func(t *testing.T) {
		entry := PlaylistEntry{
			ID:    "test123",
			Title: "Test Song",
			URL:   "https://youtube.com/watch?v=test123",
		}

		if entry.ID != "test123" {
			t.Errorf("entry.ID = %v, want 'test123'", entry.ID)
		}

		if entry.Title != "Test Song" {
			t.Errorf("entry.Title = %v, want 'Test Song'", entry.Title)
		}

		if entry.URL != "https://youtube.com/watch?v=test123" {
			t.Errorf("entry.URL = %v, want correct URL", entry.URL)
		}
	})

	t.Run("PlaylistEntry_零值", func(t *testing.T) {
		var entry PlaylistEntry

		if entry.ID != "" {
			t.Error("Zero value entry.ID should be empty")
		}

		if entry.Title != "" {
			t.Error("Zero value entry.Title should be empty")
		}

		if entry.URL != "" {
			t.Error("Zero value entry.URL should be empty")
		}
	})
}

func TestIsPlaylistURL_ExtendedCases(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "播放清單_標準格式",
			url:  "https://www.youtube.com/playlist?list=PLxxx",
			want: true,
		},
		{
			name: "播放清單_帶視頻ID",
			url:  "https://www.youtube.com/watch?v=xxx&list=PLxxx",
			want: true,
		},
		{
			name: "播放清單_短連結",
			url:  "https://youtu.be/xxx?list=PLxxx",
			want: true,
		},
		{
			name: "單個影片_無list參數",
			url:  "https://www.youtube.com/watch?v=xxx",
			want: false,
		},
		{
			name: "單個影片_短連結",
			url:  "https://youtu.be/xxx",
			want: false,
		},
		{
			name: "搜索關鍵字",
			url:  "周杰倫 晴天",
			want: false,
		},
		{
			name: "空字符串",
			url:  "",
			want: false,
		},
		{
			name: "非YouTube網址",
			url:  "https://example.com/video",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPlaylistURL(tt.url)
			if got != tt.want {
				t.Errorf("IsPlaylistURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}
