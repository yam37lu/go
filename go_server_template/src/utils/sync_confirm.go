package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const checksumFile = ".checksum"

type ChecksumMeta struct {
	Md5Value string `json:"md5"`
	Size     int64  `json:"size"`
}

// gt=1生成size checksum，gt=2生成md5 checksum，gt=3生成size/md5 checksum
func generateChecksum(path string, gt int) (map[string]*ChecksumMeta, error) {
	checksums := make(map[string]*ChecksumMeta)
	if paths, err := ListDirector(path); err != nil {
		return nil, err
	} else {
		for i := 0; i < len(paths); i++ {
			if paths[i] == "." || paths[i] == ".." || endOf(paths[i], checksumFile) {
				continue
			}
			currentPath := path + "/" + paths[i]
			if IsDirector(currentPath) {
				if s, err := generateChecksum(currentPath, gt); err != nil {
					return nil, err
				} else {
					for k, v := range s {
						checksums[paths[i]+"/"+k] = v
					}
				}
			} else {
				meta := ChecksumMeta{}
				if (gt & 0x01) == 1 {
					if l, err := FileSize(currentPath); err != nil {
						return nil, err
					} else {
						if l == 0 { // 文件大小为0不处理
							continue
						}
						meta.Size = l
					}
				}
				if (gt & 0x02) == 2 {
					if s, err := MD5File(currentPath); err != nil {
						return nil, err
					} else {
						meta.Md5Value = s
					}
				}
				if (gt & 0x03) != 0 {
					checksums[paths[i]] = &meta
				}
			}
		}
	}
	return checksums, nil
}

func GenerateChecksumContent(path string) (string, error) {
	if checksums, err := generateChecksum(path, 1); err != nil {
		return "", err
	} else {
		if content, err := json.Marshal(checksums); err != nil {
			return "", err
		} else {
			return string(content), nil
		}
	}
}

func GenerateChecksumFile(path string) error {
	if content, err := GenerateChecksumContent(path); err != nil {
		return err
	} else {
		if size, err := GetDirectorSize(path); err != nil {
			return err
		} else {
			if err = WriteContent(fmt.Sprintf(
				"%s/%d%s", path, size, checksumFile), content); err != nil {
				return err
			} else {
				return nil
			}
		}
	}
}

func ChecksumFile(path string) (bool, error) {
	filename, err := GetChecksumFile(path)
	if err != nil {
		return false, err
	}
	checksums := make(map[string]*ChecksumMeta)
	content, err := GetContent(path + "/" + filename)
	if err != nil {
		return false, err
	}
	if err = json.Unmarshal([]byte(content), &checksums); err != nil {
		return false, err
	}
	checksumSize, err := generateChecksum(path, 1)
	if err != nil {
		return false, err
	}
	if len(checksumSize) != len(checksums) {
		return false, nil
	}
	for k, v := range checksums {
		if l, ok := checksumSize[k]; !ok || v.Size != l.Size {
			return false, nil
		}
	}
	// md5校验
	//checksumValue, err := generateChecksum(path, 2)
	//if err != nil {
	//	return false, err
	//}
	//if len(checksumValue) != len(checksums) {
	//	return false, nil
	//}
	//for k, v := range checksums {
	//	if s, ok := checksumValue[k]; !ok || v.Md5Value != s.Md5Value {
	//		return false, nil
	//	}
	//}
	return true, nil
}

func endOf(filename string, sub string) bool {
	idx := strings.Index(filename, sub)
	if idx > 0 && idx == (len(filename)-len(sub)) {
		return true
	}
	return false
}

func GetChecksumFile(path string) (string, error) {
	if subs, err := ListDirector(path); err == nil {
		for _, sub := range subs {
			if endOf(sub, checksumFile) {
				return sub, nil
			}
		}
		return "", errors.New("can't find checksum file")
	} else {
		return "", err
	}
}

func GetChecksumFilePrefixSize(filename string) int64 {
	idx := strings.Index(filename, checksumFile)
	if idx <= 0 {
		return 0
	}
	if size, err := strconv.ParseInt(filename[0:idx], 10, 64); err != nil {
		return 0
	} else {
		return size
	}
}
