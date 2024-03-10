package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func IsSymlink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}

	// 在Windows上，我们可以直接检查info.Mode()来确定是否为软链接
	// 软链接的模式位是os.ModeSymlink
	return info.Mode()&os.ModeSymlink != 0, nil
}

func aliasPath(p string) string {
	p = strings.ToLower(strings.TrimSpace(p))
	if strings.HasPrefix(p, "%appdata%") {
		appDataPath := os.Getenv("APPDATA")
		if appDataPath == "" {
			// 如果APPDATA环境变量不存在，尝试使用默认路径
			// Windows系统下，APPDATA环境变量通常指向
			// "C:\Users\%USERNAME%\AppData\Roaming"
			appDataPath = filepath.Join(usr.HomeDir, "AppData", "Roaming")
		}
		p = strings.Replace(p, "%appdata%", appDataPath, 1)
	} else if strings.HasPrefix(p, "%localappdata%") {
		localAppData := os.Getenv("LocalAppData")
		if localAppData == "" {
			localAppData = filepath.Join(usr.HomeDir, "AppData", "Local")
		}
		p = strings.Replace(p, "%localappdata%", localAppData, 1)
	}
	return filepath.Clean(p)
}

func GetSymlinkRealPath(symlink string) (string, error) {
	targetPath, err := os.Readlink(symlink)
	if err != nil {
		return "", err
	}
	absTargetPath, err := filepath.EvalSymlinks(targetPath)
	if err != nil {
		return "", err
	}
	return absTargetPath, nil
}

func GetConfig(cfg *Config) error {

	// json5path := filepath.Join(rootPath, "config.json5")
	// var file *os.File
	// file, err := os.Open(json5path)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// dec := json5.NewDecoder(file)
	// err = dec.Decode(&cfg)

	dataStr, err := os.ReadFile(filepath.Join(rootPath, "config.json"))
	if err == nil {
		err = json.Unmarshal(dataStr, &cfg)
	}
	return err
}

func SetConfig(cfg *Config) {

	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		log.Info(err)
		return
	}

	err = os.WriteFile(filepath.Join(rootPath, "config.json"), data, 0644)
	if err != nil {
		log.Warn(err)
	}
}

// func MkShortcut(exepath string, exeargs []string) error {
// 	baseName := filepath.Base(exepath)
// 	// pPath := filepath.Dir(exepath)
// 	dotIndex := strings.LastIndex(baseName, ".")
// 	if dotIndex > 0 {
// 		baseName = baseName[:dotIndex]
// 	}
// 	baseName += "_runas_appman.lnk"

// 	src, _ := os.Executable()
// 	src = `"` + src + `" "` + exepath + `"`
// 	if len(exeargs) > 0 {
// 		src += ` ` + strings.Join(exeargs, " ")
// 	}

// 	// ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
// 	// oleShellObject, err := oleutil.CreateObject("WScript.Shell")
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// defer oleShellObject.Release()
// 	// wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// defer wshell.Release()

// 	// cs, err := oleutil.CallMethod(wshell, "CreateShortcut", filepath.Join(pPath, baseName))
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// idispatch := cs.ToIDispatch()
// 	// log.Info(src)
// 	// _, err = oleutil.PutProperty(idispatch, "TargetPath", src)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// log.Info(exepath)
// 	// _, err = oleutil.PutProperty(idispatch, "IconLocation", exepath, "0")
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// log.Info("save")
// 	// _, err = oleutil.CallMethod(idispatch, "Save")

// 	// sc := shortcut.Shortcut{
// 	// 	ShortcutPath:     filepath.Join(pPath, baseName),
// 	// 	Target:           src,
// 	// 	IconLocation:     exepath + ",0",
// 	// 	Arguments:        "",
// 	// 	Description:      "",
// 	// 	Hotkey:           "",
// 	// 	WindowStyle:      "1",
// 	// 	WorkingDirectory: "",
// 	// }
// 	// sc := shortcut.Shortcut{
// 	// 	ShortcutPath:     "./google.lnk",
// 	// 	Target:           "https://google.com",
// 	// 	IconLocation:     "%SystemRoot%\\System32\\SHELL32.dll,0",
// 	// 	Arguments:        "",
// 	// 	Description:      "",
// 	// 	Hotkey:           "",
// 	// 	WindowStyle:      "1",
// 	// 	WorkingDirectory: "",
// 	// }
// 	// err := shortcut.Create(sc)

// 	shortcut.CreateDesktopShortcut("google", "https://google.com", "%SystemRoot%\\System32\\SHELL32.dll,0")

// 	// if err != nil {
// 	// 	return err
// 	// }
// 	return nil
// }

// Set oWS = WScript.CreateObject("WScript.Shell")
// sLinkFile = "C:\MyShortcut.LNK"
// Set oLink = oWS.CreateShortcut(sLinkFile)
//     oLink.TargetPath = "C:\Program Files\MyApp\MyProgram.EXE"
//  '  oLink.Arguments = ""
//  '  oLink.Description = "MyProgram"
//  '  oLink.HotKey = "ALT+CTRL+F"
//  '  oLink.IconLocation = "C:\Program Files\MyApp\MyProgram.EXE, 2"
//  '  oLink.WindowStyle = "1"
//  '  oLink.WorkingDirectory = "C:\Program Files\MyApp"
// oLink.Save

// http://nircmd.nirsoft.net/shortcut.html
// shortcut "f:\winnt\system32\calc.exe" "~$folder.desktop$" "Windows Calculator"
// shortcut "f:\winnt\system32\calc.exe" "~$folder.programs$\Calculators" "Windows Calculator"
// shortcut "f:\Program Files\KaZaA\Kazaa.exe" "c:\temp\MyShortcuts" "Kazaa"
// shortcut "f:\Program Files" "c:\temp\MyShortcuts" "Program Files Folder" "" "f:\winnt\system32\shell32.dll" 45
// shortcut "f:\Program Files" "c:\temp\MyShortcuts" "Program Files Folder" "" "" "" "max"

// shortcut "E:\Project\Golang\appman\build\appman.exe" "D:\Software\Other" "Ruler_as_appman" "D:\Software\Other\Ruler.exe" "D:\Software\Other\Ruler.exe" 0 "min" "D:\Software\Other"
