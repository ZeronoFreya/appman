package main

import (
	"os"
	"path/filepath"
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
