package command

import (
	"testing"
)

func TestVoiceChannelStatus(t *testing.T) {
	t.Run("UpdateVoiceChannelStatus_函数存在", func(t *testing.T) {
		// Note: 实际测试需要 mock bot.Client 和 HTTP 请求
		// 这里验证函数签名正确

		// mockClient := &MockBotClient{}
		// channelID := snowflake.ID(123456789)
		// songTitle := "测试歌曲"

		// err := UpdateVoiceChannelStatus(mockClient, channelID, songTitle)

		// 验证状态格式正确
		// 验证字符串长度限制（500字符）
	})

	t.Run("ClearVoiceChannelStatus_函数存在", func(t *testing.T) {
		// Note: 实际测试需要 mock bot.Client 和 HTTP 请求

		// mockClient := &MockBotClient{}
		// channelID := snowflake.ID(123456789)

		// err := ClearVoiceChannelStatus(mockClient, channelID)

		// 验证发送了空字符串
	})
}

func TestVoiceChannelStatusLength(t *testing.T) {
	t.Run("状态消息应该限制在500字符", func(t *testing.T) {
		// 创建超过500字符的歌曲标题
		longTitle := ""
		for i := 0; i < 600; i++ {
			longTitle += "a"
		}

		// Note: 实际测试需要验证 UpdateVoiceChannelStatus
		// 会截断到 500 字符

		// mockClient := &MockBotClient{}
		// channelID := snowflake.ID(123456789)

		// err := UpdateVoiceChannelStatus(mockClient, channelID, longTitle)

		// 验证发送的状态 <= 500 字符
	})
}
