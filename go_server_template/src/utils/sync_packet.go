package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var fileSizeThreshold int64 = 1024 * 1024 * 1024 //2G文件大小
var fileSizeStep int64 = 100 * 1024 * 1024       //为了少占用内存，一个分片文件分多次读取和写入，每次的大小
var tmpFileSuffix string = "_TMPFILE"
var finalSuffix string = "_FINAL"

//判断srcStr的结尾是否为endStr
func endWithStr(srcStr, endStr string) bool {
	lengthSrc := len(srcStr)
	lengthEnd := len(endStr)
	if lengthSrc < lengthEnd {
		return false
	}
	newStr := srcStr[lengthSrc-lengthEnd:]
	if newStr == endStr {
		return true
	} else {
		return false
	}
}

//从FINAL文件名中获取原文件名
func getFileName(fileName string) string {
	site := strings.LastIndex(fileName, tmpFileSuffix)
	if site == -1 {
		return fileName
	}
	return fileName[0:site]
}

//拆分文件，拆分成file_TMPFILE0、file_TMPFILE1、、、file_TMPFILE10_FINAL的格式
func splitFile(fileName string) error {
	fileInfo, err := os.Lstat(fileName)
	if err != nil {
		return err
	}
	thisFileSize := fileInfo.Size()
	if thisFileSize <= fileSizeThreshold {
		return nil
	}

	splitNum := 0
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		//	fmt.Println("Open file error!", err)
		return err
	}
	for ; thisFileSize > 0; thisFileSize -= fileSizeThreshold {
		var splitSize int64 = 0
		splitFileName := fileName
		if thisFileSize < fileSizeThreshold {
			splitSize = thisFileSize
			splitFileName = splitFileName + tmpFileSuffix + strconv.Itoa(splitNum) + finalSuffix
			splitNum++
		} else {
			splitSize = fileSizeThreshold
			splitFileName = splitFileName + tmpFileSuffix + strconv.Itoa(splitNum)
			splitNum++
		}
		err = writeBigBufToFile(splitSize, file, splitFileName) //把分片内容写入分片文件
		if err != nil {
			file.Close()
			return err
		}
	}
	file.Close()
	os.Remove(fileName) //删除原大文件
	return nil
}
func writeBigBufToFile(splitSize int64, file *os.File, splitFileName string) error {
	//把分片内容写入分片文件，因为一个分片设计成2G，内存占用还是太大，所以还是要拆小分多次写入文件
	f, err := os.OpenFile(splitFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	thiSplitSize := fileSizeStep
	for splitSize > 0 {
		if thiSplitSize > splitSize {
			thiSplitSize = splitSize
		}
		buf := make([]byte, thiSplitSize)
		_, err := file.Read(buf)
		if err != nil {
			return err
		}

		_, err = f.Write(buf)
		if err != nil {
			return err
		}
		splitSize -= thiSplitSize
	}
	return err
}

//合成文件，外界识别到文件名file_TMPFILE10_FINAL时认为需要合成，作为参数传入，全路径
func unionFile(fileName string) error {
	if !endWithStr(fileName, finalSuffix) {
		return fmt.Errorf("文件名非FINAL文件")
	}
	sourceFileName := getFileName(fileName) //获得原始文件名
	splitNum := 0
	endFlag := false
	sourceF, err := os.OpenFile(sourceFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666) //append方式添加
	if err != nil {
		return err
	}
	defer sourceF.Close()
	for !endFlag {
		splitName := sourceFileName + tmpFileSuffix + strconv.Itoa(splitNum)
		_, err := os.Lstat(splitName)
		if err != nil { //文件还不存在，则可能它是最终文件
			splitName = sourceFileName + tmpFileSuffix + strconv.Itoa(splitNum) + finalSuffix
			_, err = os.Lstat(splitName)
			if err != nil { //则确认此拆分文件还未到达，继续等待,这里后续可以添加退出逻辑，不能无限等待
				time.Sleep(5 * time.Second)
				continue
			}
			endFlag = true
		}
		if err := unionStepBufFile(splitName, sourceF); err != nil {
			return err
		}
		splitNum++
		os.Remove(splitName)
	}
	return nil
}
func unionStepBufFile(fileName string, sourceF *os.File) error {
	splitF, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer splitF.Close()
	var offset int64 = 0
	fileInfo, _ := splitF.Stat()
	length := fileInfo.Size()
	for offset < length {
		var readSize int64
		if (length - offset) < fileSizeStep {
			readSize = length - offset
		} else {
			readSize = fileSizeStep
		}
		buf := make([]byte, readSize)
		n, err := splitF.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			return err
		}
		offset += (int64)(n)

		_, err = sourceF.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

//判断文件的全部分片是否都完整具备，符合合成文件的前置条件了,输入参数为file_TMPFILE10_FINAL
func isAllFileComplete(fileName string) (bool, error) {
	if !endWithStr(fileName, finalSuffix) {
		return false, fmt.Errorf("文件名非FINAL文件")
	}
	sourceFileName := getFileName(fileName) //获得原始文件名
	splitNum := 0
	endFlag := false
	for !endFlag {
		splitName := sourceFileName + tmpFileSuffix + strconv.Itoa(splitNum)
		fileInfo, err := os.Lstat(splitName)
		if err != nil { //文件还不存在，则可能它是最终文件
			splitName = sourceFileName + tmpFileSuffix + strconv.Itoa(splitNum) + finalSuffix
			fileInfo, err = os.Lstat(splitName)
			if err != nil { //所有文件分片还未具备完全
				return false, nil
			}
			endFlag = true
		}
		if endFlag { //最后一个文件，判断它的大小不为0即可
			if fileInfo.Size() == 0 {
				return false, nil
			}
		} else {
			if fileInfo.Size() != fileSizeThreshold { //非最终文件的，它的大小一定要等于固定的切片大小
				return false, nil
			}
		}
		splitNum++
	}
	return true, nil
}

func SplitAllFile(strRootDir string) error {
	dir, err := ioutil.ReadDir(strRootDir)
	if err != nil {
		return err
	}
	for _, v := range dir {
		if v.Name() == "." || v.Name() == ".." {
			continue
		}
		if v.IsDir() {
			SplitAllFile(strRootDir + "/" + v.Name())
		} else {
			splitFile(strRootDir + "/" + v.Name()) //分拆大文件
		}
	}
	return nil
}
func UnionAllFile(strRootDir string) error {
	dir, err := ioutil.ReadDir(strRootDir)
	if err != nil {
		return err
	}
	for _, v := range dir {
		if v.Name() == "." || v.Name() == ".." {
			continue
		}
		if v.IsDir() {
			if err := UnionAllFile(strRootDir + "/" + v.Name()); err != nil {
				return err
			}
		} else {
			if !endWithStr(strRootDir+"/"+v.Name(), finalSuffix) {
				continue
			}
			bFlag, err := isAllFileComplete(strRootDir + "/" + v.Name()) //是否它的全部分片文件都具备了
			if err != nil || !bFlag {
				continue
			}
			if err := unionFile(strRootDir + "/" + v.Name()); err != nil { //组合文件
				return err
			}
		}
	}
	return nil
}
