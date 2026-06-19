package command

import "github.com/disgoorg/disgolink/v3/disgolink"

var lavalinkClient disgolink.Client

// SetLavalinkClient 設定全域 Lavalink client
func SetLavalinkClient(client disgolink.Client) {
	lavalinkClient = client
}

// GetLavalinkClient 取得全域 Lavalink client
func GetLavalinkClient() disgolink.Client {
	return lavalinkClient
}
