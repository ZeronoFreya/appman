package main

import (
	"errors"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	log      = logrus.New()
	rootPath string
	usr      *user.User
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
	logFilePath := filepath.Join(rootPath, "appman.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("无法打开日志文件: %v", err)
	}
	log.SetOutput(logFile)

	usr, err = user.Current()
	if err != nil {
		log.Fatalf("无法获取当前用户: %v", err)
		return
	}
}

func isMkSymlink(dirpath, symlink string) (bool, error) {
	// 目标路径不存在
	if !Exists(symlink) {
		return true, nil
	}
	// 判断是否是软链接
	isLink, err := IsSymlink(symlink)
	if err != nil {
		return false, err
	}
	if !isLink {
		return true, nil
	} else {
		// 检查软链接的指向是否是dirpath
		targetPath, err := GetSymlinkRealPath(symlink)
		if err != nil {
			return false, err
		}
		if targetPath != dirpath {
			return true, nil
		}
	}
	return false, nil
}

func checkSymlink(fileName, parentDir string, cfg *Config) error {
	parentDir = filepath.Clean(parentDir)
	var d []SL
	for k, v := range cfg.Symlink {
		if k == fileName {
			for _, v := range v {
				dir_v := aliasPath(v.Dir)
				if dir_v == "" || len(v.Link) == 0 {
					continue
				}
				if dir_v == parentDir {
					d = v.Link
					break
				}
			}
			break
		}
	}
	var links []string
	for _, v := range d {
		dir_v := aliasPath(v.Dirpath)
		lnk_v := aliasPath(v.Symlink)
		if dir_v == "" || lnk_v == "" {
			continue
		}
		rs, err := isMkSymlink(dir_v, lnk_v)
		if err != nil {
			return err
		}
		if rs {
			links = append(links, dir_v+","+lnk_v)
		}
	}
	if len(links) > 0 {

		tempDir := filepath.Join(os.TempDir(), "run_app_temp")
		if !Exists(tempDir) {
			err := os.MkdirAll(tempDir, 0700)
			if err != nil {
				return err
			}
		}

		content := []byte(strings.Join(links[:], ";"))

		err := os.WriteFile(filepath.Join(tempDir, "temp_symlink.txt"), content, 0644)
		if err != nil {
			return err
		}

		cmd := exec.Command("schtasks", "/Run", "/TN", "noUAC.mksymlink")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000,
		}
		// var out bytes.Buffer
		// var stderr bytes.Buffer
		// cmd.Stdout = &out
		// cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return errors.New("需要运行 安装.bat")
		}

		// 设置超时时间
		timeoutDuration := 10 * time.Second
		// 记录开始时间
		startTime := time.Now()

		errMsg := ""
		rs_symlink := filepath.Join(os.TempDir(), "run_app_temp", "rs_symlink.txt")
		for {
			// 检查是否已经超时
			if time.Since(startTime) >= timeoutDuration {
				errMsg = "超时，文件读取未完成。"
				break
			}
			if Exists(rs_symlink) {
				content, err := os.ReadFile(rs_symlink)
				if err != nil {
					return err
				}
				errMsg = string(content)
				if err = os.Remove(rs_symlink); err != nil {
					log.Info(err.Error())
				}
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		if errMsg != "" {
			return errors.New(errMsg)
		}
	}

	return nil
}

func runas(exepath, exeargs string) {
	target := exepath + " " + exeargs
	content := []byte(target)
	err := os.WriteFile(filepath.Join(os.TempDir(), "run_app_temp", "temp_app_runas_admin.txt"), content, 0644)
	if err != nil {
		log.Fatal("无法写入 temp_app_runas_admin: " + err.Error())
		return
	}
	cmd := exec.Command("schtasks", "/Run", "/TN", "noUAC.runadmin")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	if err := cmd.Run(); err != nil {
		log.Fatal("需要运行 安装.bat")
	}
}

func main() {
	args := os.Args[1:]
	len_args := len(args)
	if len_args == 0 {
		log.Fatal("未传递参数")
		return
	}
	exepath := strings.ToLower(filepath.Clean(args[0]))
	var exeargs []string

	if len_args > 1 {
		exeargs = args[1:]
	}

	if !Exists(exepath) {
		log.Fatal("目标程序不存在")
		return
	}

	fileName := strings.ToLower(filepath.Base(exepath))
	parentDir := strings.ToLower(filepath.Dir(exepath))

	var cfg Config

	err := GetConfig(&cfg)
	if err != nil {
		log.Fatal("无法获取配置: " + err.Error())
		return
	}

	inSymlink := false
	for k, _ := range cfg.Symlink {
		if k == fileName {
			inSymlink = true
			break
		}
	}
	if inSymlink {
		err = checkSymlink(fileName, parentDir, &cfg)
		if err != nil {
			log.Fatal("无法检查软链接: " + err.Error())
			return
		}
	}

	isRoot := false
	for _, v := range cfg.Nouac {
		temp := strings.ToLower(filepath.Clean(strings.TrimSpace(v)))
		if temp == exepath {
			isRoot = true
			break
		}
	}
	if isRoot {
		runas(exepath, strings.Join(exeargs, " "))
		return
	} else if !inSymlink {
		// 此程序路径没有记录
		nircmd := filepath.Join(rootPath, "tools", "nircmd.exe")
		if !Exists(nircmd) {
			log.Fatal("nircmd.exe不存在")
			return
		}
		filename, _ := os.Executable()
		folder := filepath.Dir(exepath)
		baseName := filepath.Base(exepath)
		dotIndex := strings.LastIndex(baseName, ".")
		shortcutname := baseName
		if dotIndex > 0 {
			shortcutname = baseName[:dotIndex]
		}
		// cl := fmt.Sprintf(`%s shortcut "%s" "%s" "%s_as_appman" "%s" "%s" 0 "min" "%s"`, nircmd, filename, folder, baseName, exepath, exepath, folder)

		cmd := exec.Command(nircmd, "shortcut", filename, folder, shortcutname+"_as_appman", exepath, exepath, "0", "min", folder)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		cfg.Symlink[strings.ToLower(baseName)] = []SymlinkData{
			{
				Dir:  strings.ToLower(filepath.Clean(folder)),
				Link: []SL{},
			},
		}
		SetConfig(&cfg)
		return
	}

	cmd := exec.Command(exepath, exeargs...)
	if err := cmd.Run(); err != nil {
		runas(exepath, strings.Join(exeargs, " "))
	}
}
