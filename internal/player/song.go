package player

// Song describes one queued or currently playing audio item.
type Song struct {
	Title       string
	URL         string
	StreamURL   string
	RequestedBy string
}
