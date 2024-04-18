package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	log      = logrus.New()
	rootPath string
)

func init() {
	ex, err := os.Executable()
	rootPath = filepath.Dir(ex)

	// 设置日志级别
	log.SetLevel(logrus.InfoLevel)

	// 设置日志格式
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置日志输出到文件
	logFilePath := filepath.Join(rootPath, "runadmin.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("无法打开日志文件: %v", err)
	}
	log.SetOutput(logFile)
}

func main() {
	filepath := filepath.Join(os.TempDir(), "run_app_temp", "temp_app_runas_admin.txt")
	if !Exists(filepath) {
		log.Fatal("temp_app_runas_admin.txt不存在")
		return
	}
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("无法读取temp_app_runas_admin.txt")
		return
	}
	rs := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
	for _, v := range rs {
		// cmd := exec.Command("cmd.exe", "/c", "start "+v)
		cmd := exec.Command("cmd.exe")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CmdLine:       `/c chcp 65001 && start "" ` + v,
			HideWindow:    true,
			CreationFlags: 0x08000000,
		}
		// cmd.Start()
		if err := cmd.Run(); err != nil {
			log.Info(err.Error())
		}
	}

	os.Remove(filepath)
}
