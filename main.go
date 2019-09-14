package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	from := os.Args[1]
	to := os.Args[2]

	// 遍历输入文件夹
	err := filepath.Walk(from, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			//return filepath.SkipDir // 忽略子目录
			return nil
		}
		if !isPic(path) && !isMov(path) {
			fmt.Println("is not pic or movie, ignore it", path)
			return nil
		}

		destPath := getDestAbsPath(to, path)
		destDir := filepath.Dir(getDestAbsPath(to, path))
		os.MkdirAll(destDir, 0755)

		if IsExist(destPath) {
			return nil
		}

		// 移动文件
		if _, err = CopyFile(path, destPath); err != nil {
			fmt.Printf("copy %s failed; %v\n", path, err)
			return nil
		}
		// 删除源文件
		os.Remove(path)

		return nil
	})
	if err != nil {
		log.Fatal("filepath.Walk failed; detail: ", err)
	}
}

// 根据照片拍摄日期决定存储目录
func getPlacePath(tm time.Time) string {
	return filepath.Join(strconv.Itoa(tm.Year()), fmt.Sprintf("%02d", tm.Month()))
}

// 决定照片的存储位置的绝对路径
func getDestAbsPath(dest string, src string) string {
	name := filepath.Base(src)
	tm := FetchTokenTime(src)
	placePath := getPlacePath(tm)
	destPath := filepath.Join(dest, placePath, name)
	absPath, _ := filepath.Abs(destPath)
	if absPath == "" {
		absPath = destPath
	}
	return absPath
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
	// 或者
	//return err == nil || !os.IsNotExist(err)
	// 或者
	//return !os.IsNotExist(err)
}

func CopyFile(src, des string) (written int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	//获取源文件的权限
	fi, _ := srcFile.Stat()
	perm := fi.Mode()

	//desFile, err := os.Create(des)  //无法复制源文件的所有权限
	desFile, err := os.OpenFile(des, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm) //复制源文件的所有权限
	if err != nil {
		return 0, err
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
}
