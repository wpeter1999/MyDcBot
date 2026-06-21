package command

// CommandRegistry 存放所有要註冊的指令
// 新增指令時，在此加入即可
var CommandRegistry = []*BotCommand{
	PlayCommand,
	SkipCommand,
	PauseCommand,
	QueueCommand,
	StopCommand,
	NowPlayingCommand,
	DownloadCommand,
	LoopCommand,    // loop 指令
	ShuffleCommand, // shuffle 指令
	HelpCommand,    // help 指令
}
