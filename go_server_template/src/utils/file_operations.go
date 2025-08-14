package utils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func FileMove(src string, dst string) (err error) {
	if dst == "" {
		return nil
	}
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}
	var revoke = false
	dir := filepath.Dir(dst)
Redirect:
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
		if !revoke {
			revoke = true
			goto Redirect
		}
	}
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	if flag, _ := PathExists(dst); !flag {
		if err = CreateDirector(dst); err != nil {
			return err
		}
	}
	return CopyDirector(src, dst)
}

func TrimSpace(target interface{}) {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Ptr {
		return
	}
	t = t.Elem()
	v := reflect.ValueOf(target).Elem()
	for i := 0; i < t.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).SetString(strings.TrimSpace(v.Field(i).String()))
		}
	}
	// return
}

func GetContent(filename string) (string, error) {
	if content, err := ioutil.ReadFile(filename); err != nil {
		return "", err
	} else {
		return string(content), nil
	}
}

func WriteContent(filename string, content string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if file, err := os.Create(filename); err == nil {
			defer file.Close()
			if _, err := io.WriteString(file, content); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if file, err := os.Open(filename); err == nil {
			defer file.Close()
			if _, err := io.WriteString(file, content); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
func IsFile(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}

func FileSize(filename string) (int64, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return -1, err
	}
	return fi.Size(), nil
}

func FileExists(filename string) bool {
	existed := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		existed = false
	}
	return existed
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
