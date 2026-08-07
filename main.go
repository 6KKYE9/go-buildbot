package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"go-buildbot/internal/bot"
)

var dbPath = "buildbot.json"

func run() error {
	runCfg := flag.String("run", "", "运行构建（JSON 配置字符串）")
	list := flag.Bool("list", false, "列出构建历史")
	watch := flag.Int("watch", 0, "持续监控间隔（秒）")
	cf := flag.String("f", "", "配置文件路径")
	flag.Parse()

	s, err := bot.Load(dbPath)
	if err != nil {
		return err
	}

	switch {
	case *runCfg != "":
		var cfg bot.Config
		if err := json.Unmarshal([]byte(*runCfg), &cfg); err != nil {
			return fmt.Errorf("配置 JSON 解析失败: %w", err)
		}
		r := bot.Run(cfg, s.NextID())
		printResult(r)
	case *list:
		for _, r := range s.List() {
			status := "成功"
			if !r.Success {
				status = fmt.Sprintf("失败(%d)", r.ExitCode)
			}
			fmt.Printf("#%d [%s] %s %s\n", r.ID, r.Time.Format("01-02 15:04"), r.Command, status)
		}
	case *watch > 0:
		if *cf == "" {
			return fmt.Errorf("-watch 需要 -f 指定配置文件")
		}
		data, err := os.ReadFile(*cf)
		if err != nil {
			return err
		}
		var cfg bot.Config
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		fmt.Printf("监控 %s，每 %d 秒一次\n", cfg.RepoURL, *watch)
		ticker := time.NewTicker(time.Duration(*watch) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			r := bot.Run(cfg, s.NextID())
			printResult(r)
		}
	default:
		return fmt.Errorf("用法: -run/-list/-watch")
	}
	return nil
}

func printResult(r bot.BuildResult) {
	if r.Success {
		fmt.Printf("#%d 成功 (%s)\n", r.ID, r.Duration)
	} else {
		fmt.Printf("#%d 失败 退出码 %d (%s)\n", r.ID, r.ExitCode, r.Duration)
		if r.Output != "" {
			fmt.Println(r.Output)
		}
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}
