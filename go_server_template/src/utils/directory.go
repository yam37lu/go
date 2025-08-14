package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"template/global"

	"go.uber.org/zap"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDirector(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			global.SYS_LOG.Debug("create directory" + v)
			err = os.MkdirAll(v, os.ModePerm)
			if err != nil {
				global.SYS_LOG.Error("create directory"+v, zap.Any(" error:", err))
			}
		}
	}
	return err
}

func ListDirector(path string) ([]string, error) {
	if rd, err := ioutil.ReadDir(path); err != nil {
		return nil, err
	} else {
		paths := make([]string, 0)
		for _, fi := range rd {
			// if fi.IsDir() {
			paths = append(paths, fi.Name())
			// }
		}
		return paths, nil
	}
}

func IsDirector(path string) bool {
	fi, e := os.Stat(path)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

func CopyDirector(srcPath, desPath string) error {
	//检查目录是否正确
	if srcInfo, err := os.Stat(srcPath); err != nil {
		return err
	} else {
		if !srcInfo.IsDir() {
			return errors.New("source director is error！")
		}
	}

	if desInfo, err := os.Stat(desPath); err != nil {
		return err
	} else {
		if !desInfo.IsDir() {
			return errors.New("dest director is error！")
		}
	}

	if strings.TrimSpace(srcPath) == strings.TrimSpace(desPath) {
		return errors.New("source director must be not same to dest director！")
	}

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		//复制目录是将源目录中的子目录复制到目标路径中，不包含源目录本身
		if path == srcPath {
			return nil
		}

		//生成新路径
		destNewPath := strings.Replace(path, srcPath, desPath, -1)

		if !f.IsDir() {
			CopyFile(path, destNewPath)
		} else {
			if !FileExists(destNewPath) {
				return CreateDirector(destNewPath)
			}
		}

		return nil
	})
	return err
}

func RemoveDirector(dir string) error {
	return os.RemoveAll(dir)
}

func GetDirectorSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func IsDirectorEmpty(path string) (bool, error) {
	subs, err := ListDirector(path)
	if err != nil {
		return true, err
	}
	for _, sub := range subs {
		if sub == "." || sub == ".." {
			continue
		}
		if IsDirector(path + "/" + sub) {
			if flag, _ := IsDirectorEmpty(path + "/" + sub); !flag {
				return false, nil
			}
		} else {
			return false, nil
		}
	}
	return true, nil
}
