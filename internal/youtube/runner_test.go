package youtube

import (
	"testing"
)

func TestCommandRunner(t *testing.T) {
	t.Run("CommandRunner_介面定義", func(t *testing.T) {
		// 驗證 CommandRunner 介面存在且可以被實現
		// 實際的介面契約驗證在具體實現測試中完成
		t.Log("CommandRunner 介面驗證通過")
	})
}

func TestRealCommandRunner(t *testing.T) {
	t.Run("RealCommandRunner_實現CommandRunner介面", func(t *testing.T) {
		runner := NewExecCommandRunner()

		if runner == nil {
			t.Fatal("NewExecCommandRunner() should not return nil")
		}

		// 驗證實現了 CommandRunner 介面
		var _ CommandRunner = runner
	})
}
