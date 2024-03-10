package main

import (
	"bytes"
	"errors"
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
	logFilePath := filepath.Join(rootPath, "mksymlink.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("无法打开日志文件: %v", err)
	}
	log.SetOutput(logFile)
}

func mkSymlink(dirpath, symlink string) error {
	// 目标路径存在
	if Exists(symlink) {
		// 判断是否是软链接
		isLink, err := IsSymlink(symlink)
		if err != nil {
			return err
		}
		// 不是软链接
		if !isLink {
			// 原始目录不存在
			if !Exists(dirpath) {
				if err := os.MkdirAll(dirpath, 0755); err != nil {
					return err
				}
				var stderr bytes.Buffer
				// 复制子文件到原始目录
				// robocopy 不遵守标准返回格式，需特殊处理返回值
				cmd := exec.Command("robocopy", symlink, dirpath, "/e", "/mir")
				cmd.SysProcAttr = &syscall.SysProcAttr{
					HideWindow:    true,
					CreationFlags: 0x08000000,
				}
				cmd.Stderr = &stderr
				if err := cmd.Run(); err != nil {
					if stderr.Len() > 0 {
						return errors.New(stderr.String())
					}
				}

			}
			// 删除目录
			if err = os.RemoveAll(symlink); err != nil {
				return err
			}
			// 创建软链接
			if err = os.Symlink(dirpath, symlink); err != nil {
				return err
			}
		} else {
			// 检查软链接的指向是否是dirpath
			targetPath, err := GetSymlinkRealPath(symlink)
			if err != nil {
				return err
			}
			if targetPath != dirpath {
				// 删除目录
				if err = os.RemoveAll(symlink); err != nil {
					return err
				}
				// 创建软链接
				if err = os.Symlink(dirpath, symlink); err != nil {
					return err
				}
			}
		}

	} else {
		// 原始目录不存在
		if !Exists(dirpath) {
			if err := os.MkdirAll(dirpath, 0755); err != nil {
				return err
			}
		}
		// 软链接父级目录
		parentDir := filepath.Dir(symlink)
		if !Exists(parentDir) {
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return err
			}
		}
		// 创建软链接
		if err := os.Symlink(dirpath, symlink); err != nil {
			return err
		}
	}
	return nil
}

func callback(result bool, msg string) {
	content := ""
	if !result {
		content = msg
	}

	tempFilePath := filepath.Join(os.TempDir(), "run_app_temp", "rs_temp_symlink.txt")

	err := os.WriteFile(tempFilePath, []byte(content), 0644)
	if err != nil {
		log.Fatal("无法写入 rs_temp_symlink.txt:" + err.Error())
		return
	}
	targetFilePath := filepath.Join(os.TempDir(), "run_app_temp", "rs_symlink.txt")
	err = os.Rename(tempFilePath, targetFilePath)
	if err != nil {
		log.Fatal("无法重命名 rs_temp_symlink.txt:" + err.Error())
		return
	}
}

func main() {
	filepath := filepath.Join(os.TempDir(), "run_app_temp", "temp_symlink.txt")
	if !Exists(filepath) {
		return
	}
	content, err := os.ReadFile(filepath)
	if err != nil {
		callback(false, err.Error())
		return
	}
	rs := strings.Split(string(content), ";")
	for _, v := range rs {
		d := strings.Split(v, ",")
		if len(d) < 2 {
			callback(false, "temp_symlink格式错误:"+v)
			return
		}
		err := mkSymlink(d[0], d[1])
		if err != nil {
			callback(false, err.Error())
			return
		}
	}
	if err = os.Remove(filepath); err != nil {
		log.Info(err)
	}
	callback(true, "")
}
