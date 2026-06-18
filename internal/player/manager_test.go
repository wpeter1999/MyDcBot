package player

import "testing"

func TestManager_GetOrCreateReturnsSamePlayerForSameGuild(t *testing.T) {
	manager := NewManager(50)

	first := manager.GetOrCreate("guild-1")
	second := manager.GetOrCreate("guild-1")

	if first != second {
		t.Fatal("相同 GuildID 應取得同一個 GuildPlayer")
	}
}

func TestManager_GetOrCreateReturnsDifferentPlayersForDifferentGuilds(t *testing.T) {
	manager := NewManager(50)

	guildA := manager.GetOrCreate("guild-a")
	guildB := manager.GetOrCreate("guild-b")

	if guildA == guildB {
		t.Fatal("不同 GuildID 應取得不同 GuildPlayer")
	}
	if guildA.GuildID() != "guild-a" {
		t.Fatalf("guildA GuildID 應為 guild-a，實際為 %q", guildA.GuildID())
	}
	if guildB.GuildID() != "guild-b" {
		t.Fatalf("guildB GuildID 應為 guild-b，實際為 %q", guildB.GuildID())
	}
}

func TestManager_RemoveStopsAndDeletesPlayer(t *testing.T) {
	manager := NewManager(50)
	player := manager.GetOrCreate("guild-1")

	removed := manager.Remove("guild-1")
	if !removed {
		t.Fatal("Remove 已存在的 player 應回傳 true")
	}
	if !player.IsStopped() {
		t.Fatal("Remove 應停止被移除的 player")
	}
	if _, ok := manager.Get("guild-1"); ok {
		t.Fatal("Remove 後不應再取得 player")
	}
	if manager.Remove("guild-1") {
		t.Fatal("Remove 不存在的 player 應回傳 false")
	}
}
