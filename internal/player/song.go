package player

// Song 描述一首已加入佇列或正在播放的歌曲。
type Song struct {
	Title       string
	URL         string
	StreamURL   string
	RequestedBy string
}
