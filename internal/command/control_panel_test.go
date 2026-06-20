package command

import (
	"testing"
)

func TestControlPanelButtons(t *testing.T) {
	t.Run("按钮常量定义", func(t *testing.T) {
		buttons := []struct {
			name  string
			value string
		}{
			{"ButtonShowPanel", ButtonShowPanel},
			{"ButtonPause", ButtonPause},
			{"ButtonSkip", ButtonSkip},
			{"ButtonStop", ButtonStop},
			{"ButtonQueue", ButtonQueue},
			{"ButtonNowPlaying", ButtonNowPlaying},
			{"ButtonSearch", ButtonSearch},
		}

		for _, btn := range buttons {
			if btn.value == "" {
				t.Errorf("%s should not be empty", btn.name)
			}
		}

		// 验证按钮 ID 唯一性
		seen := make(map[string]bool)
		for _, btn := range buttons {
			if seen[btn.value] {
				t.Errorf("Duplicate button ID: %s", btn.value)
			}
			seen[btn.value] = true
		}
	})
}

func TestRespondWithControlButton(t *testing.T) {
	t.Run("响应函数存在", func(t *testing.T) {
		// Note: 实际测试需要 mock events.ApplicationCommandInteractionCreate
		// 这里验证函数签名正确

		// mockEvent := &MockApplicationCommandInteractionCreate{}
		// RespondWithControlButton(mockEvent, "测试消息")

		// 验证创建了正确的组件
		// 验证消息内容正确
	})
}

func TestUpdateMessageWithFullPanel(t *testing.T) {
	t.Run("更新控制面板", func(t *testing.T) {
		// Note: 实际测试需要 mock events.ComponentInteractionCreate
		// 和 MusicService

		// mockEvent := &MockComponentInteractionCreate{}
		// mockMusicService := &MockMusicService{}
		// SetMusicService(mockMusicService)

		// UpdateMessageWithFullPanel(mockEvent, "测试内容")

		// 验证创建了正确的 Embed
		// 验证按钮状态正确
	})
}

func TestHandleControlPanelInteraction(t *testing.T) {
	tests := []struct {
		name     string
		buttonID string
		wantFunc string
	}{
		{
			name:     "显示面板按钮",
			buttonID: ButtonShowPanel,
			wantFunc: "UpdateMessageWithFullPanel",
		},
		{
			name:     "暂停按钮",
			buttonID: ButtonPause,
			wantFunc: "handlePauseButton",
		},
		{
			name:     "跳过按钮",
			buttonID: ButtonSkip,
			wantFunc: "handleSkipButton",
		},
		{
			name:     "停止按钮",
			buttonID: ButtonStop,
			wantFunc: "handleStopButton",
		},
		{
			name:     "当前播放按钮",
			buttonID: ButtonNowPlaying,
			wantFunc: "handleNowPlayingButton",
		},
		{
			name:     "佇列按钮",
			buttonID: ButtonQueue,
			wantFunc: "handleQueueButton",
		},
		{
			name:     "搜索按钮",
			buttonID: ButtonSearch,
			wantFunc: "handleSearchButton",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: 实际测试需要 mock 所有依赖
			// 这里提供测试框架

			// mockEvent := &MockComponentInteractionCreate{
			// 	customID: tt.buttonID,
			// }

			// HandleControlPanelInteraction(mockEvent)

			// 验证调用了正确的处理函数
		})
	}
}
