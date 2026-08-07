package bot

import (
	"os"
	"testing"
)

func TestRunSimple(t *testing.T) {
	cfg := Config{
		RepoURL: "test",
		Branch:  "main",
		Command: "echo ok",
	}
	r := Run(cfg, 1)
	if !r.Success {
		t.Fatalf("echo 应成功，输出: %s", r.Output)
	}
	if r.ExitCode != 0 {
		t.Fatalf("退出码应为 0，实际 %d", r.ExitCode)
	}
	defer os.Remove("buildbot.json")
}

func TestRunFailure(t *testing.T) {
	cfg := Config{
		RepoURL: "test",
		Command: "exit 1",
	}
	r := Run(cfg, 2)
	if r.Success {
		t.Fatal("exit 1 应失败")
	}
	defer os.Remove("buildbot.json")
}

func TestSaveReload(t *testing.T) {
	p := t.TempDir() + "/buildbot.json"
	s, _ := Load(p)
	s.results = append(s.results, BuildResult{ID: 1, Success: true})
	s.Save()

	s2, _ := Load(p)
	if len(s2.List()) != 1 {
		t.Fatal("重新加载应有 1 个结果")
	}
}

func testStore(t *testing.T) *Store {
	t.Helper()
	s, _ := Load(t.TempDir() + "/buildbot.json")
	return s
}
