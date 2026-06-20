package command

import (
	"testing"
)

func TestIsPlaylistURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "播放清單_list參數",
			url:  "https://www.youtube.com/watch?v=xxx&list=PLxxx",
			want: true,
		},
		{
			name: "播放清單_playlist路徑",
			url:  "https://www.youtube.com/playlist?list=PLxxx",
			want: true,
		},
		{
			name: "單個影片",
			url:  "https://www.youtube.com/watch?v=xxx",
			want: false,
		},
		{
			name: "搜尋關鍵字",
			url:  "周杰倫 晴天",
			want: false,
		},
		{
			name: "空字串",
			url:  "",
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

func TestPlaylistEntry(t *testing.T) {
	t.Run("PlaylistEntry_結構定義", func(t *testing.T) {
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
}
