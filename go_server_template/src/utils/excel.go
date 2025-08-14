package utils

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"io"
	"mime/multipart"
	"os"
)

func GetHeadAz(length int) ([]string, error) {
	if length > 26*26 {
		return []string{}, errors.New("列数超过限制")
	}
	var az []string
	var i int
	for {
		az = append(az, string(rune('A'+i)))
		i++
		if i > 25 {
			break
		}
	}
	if length <= 26 {
		return az[0:length], nil
	}
	arr := az
	for _, val := range az {
		for _, v := range az {
			arr = append(arr, val+v)
			if len(arr) == length {
				goto ret
			}
		}
	}
ret:
	return arr, nil
}

func UploadExcel(file *multipart.FileHeader) ([][]string, error) {
	openFile, err := file.Open()
	if err != nil {
		return nil, InvaildParams.Wrap(err, "打开文件失败")
	}
	tmpFile, err := os.CreateTemp("", "upload.xlsx")
	if err != nil {
		return nil, InvaildParams.Wrap(err, "创建临时文件失败")
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err = io.Copy(tmpFile, openFile); err != nil {
		return nil, InvaildParams.Wrap(err, "写入临时文件失败")
	}
	f, err := excelize.OpenFile(tmpFile.Name())
	if err != nil {
		return nil, InvaildParams.Wrap(err, "解析失败")
	}
	sheet := f.GetSheetName(0)
	return f.GetRows(sheet)
}
