package bot

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// Config 是 bot 的配置。
type Config struct {
	RepoURL string `json:"repo_url"`
	Branch  string `json:"branch"`
	WorkDir string `json:"work_dir"`
	Command string `json:"command"` // 构建命令
}

// BuildResult 记录一次构建结果。
type BuildResult struct {
	ID       int64     `json:"id"`
	Repo     string    `json:"repo"`
	Branch   string    `json:"branch"`
	Command  string    `json:"command"`
	ExitCode int       `json:"exit_code"`
	Output   string    `json:"output"`
	Duration string    `json:"duration"`
	Time     time.Time `json:"time"`
	Success  bool      `json:"success"`
}

type Store struct {
	path    string
	results []BuildResult
	next    int64
}

func Load(path string) (*Store, error) {
	s := &Store{path: path, next: 1}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return s, nil
	}
	if err := json.Unmarshal(data, &s.results); err != nil {
		return nil, err
	}
	for _, r := range s.results {
		if r.ID >= s.next {
			s.next = r.ID + 1
		}
	}
	return s, nil
}

func (s *Store) Save() error {
	data, err := json.MarshalIndent(s.results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Run 执行一次构建。
func Run(cfg Config, id int64) BuildResult {
	r := BuildResult{
		ID:      id,
		Repo:    cfg.RepoURL,
		Branch:  cfg.Branch,
		Command: cfg.Command,
		Time:    time.Now(),
	}

	start := time.Now()
	shell, shellArg := "sh", "-c"
	if runtime.GOOS == "windows" {
		shell, shellArg = "cmd", "/c"
	}
	var cmd *exec.Cmd
	if cfg.WorkDir != "" {
		cmd = exec.Command(shell, shellArg, fmt.Sprintf("cd %s && %s", cfg.WorkDir, cfg.Command))
	} else {
		cmd = exec.Command(shell, shellArg, cfg.Command)
	}

	out, err := cmd.CombinedOutput()
	r.Output = string(out)
	r.Duration = time.Since(start).Truncate(time.Second).String()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			r.ExitCode = exitErr.ExitCode()
		} else {
			r.ExitCode = -1
		}
	} else {
		r.Success = true
	}

	s, _ := Load("buildbot.json")
	s.results = append(s.results, r)
	s.next = id + 1
	s.Save()
	return r
}

func (s *Store) NextID() int64 {
	return s.next
}

func (s *Store) List() []BuildResult {
	out := make([]BuildResult, len(s.results))
	copy(out, s.results)
	return out
}

var _ = fmt.Println
