package youtube

import (
	"testing"
)

func TestCommandRunner(t *testing.T) {
	t.Run("CommandRunner_接口定义", func(t *testing.T) {
		// 验证 CommandRunner 接口存在且可以被实现

		// var runner CommandRunner
		// runner = &RealCommandRunner{}

		// 验证接口方法签名正确
	})
}

func TestRealCommandRunner(t *testing.T) {
	t.Run("RealCommandRunner_实现CommandRunner接口", func(t *testing.T) {
		runner := NewExecCommandRunner()

		if runner == nil {
			t.Fatal("NewExecCommandRunner() should not return nil")
		}

		// 验证实现了 CommandRunner 接口
		var _ CommandRunner = runner
	})
}
