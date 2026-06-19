package youtube

import (
	"context"
	"os/exec"
)

// execCommandRunner 使用 os/exec 執行外部指令的 CommandRunner 實作。
type execCommandRunner struct{}

// NewExecCommandRunner 建立使用 os/exec 的 CommandRunner。
func NewExecCommandRunner() CommandRunner {
	return &execCommandRunner{}
}

// Run 實作 CommandRunner 介面，執行外部指令並回傳標準輸出。
func (r *execCommandRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.Output()
}
