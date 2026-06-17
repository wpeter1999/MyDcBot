package command

// CommandRegistry 存放所有要註冊的指令
// 新增指令時，在此加入即可
var CommandRegistry = []*BotCommand{
	PingCommand,
	WeatherCommand,
}