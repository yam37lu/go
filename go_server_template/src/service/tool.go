package service

import (
	"template/dao"
	"template/global"
	"template/model/request"
	"template/utils"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func ExportTool(param *request.ExportToolRequest) (*os.File, error) {
	if param.Table == "" {
		param.Table = "tb_addr_std"
	}
	if param.FileName == "" {
		param.FileName = param.Table
	}
	if param.PageSize <= 0 {
		param.PageSize = 10000
	}
	if param.Page <= 0 {
		param.Page = 1
	}
	files := make(map[string][]byte)
	for {
		rows, err := dao.ExportTool(param)
		if err != nil || len(rows) == 0 {
			break
		}
		files[fmt.Sprintf("%d.sg", param.Page)] = []byte(strings.Join(rows, "\n"))
		param.Page = param.Page + 1
	}
	if len(files) == 0 {
		return nil, utils.InvaildParams.New("无数据")
	}
	if param.Truncate {
		files[fmt.Sprintf("%d.sg", 0)] = []byte(fmt.Sprintf("truncate table %s", param.Table))
	}
	path := getFileName(param.FileName, param.FilePath)
	zip, err := utils.CreateZip(files)
	if err != nil {
		return nil, err
	}
	bytes := utils.EncodeBytes(zip, []byte(utils.MD5V([]byte(param.FileName))))
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	_, err = file.Write(bytes)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return file, err
}

func ImportTool(param *request.ImportToolRequest) (int64, error) {
	if param.FileName == "" {
		param.FileName = "tb_addr_std"
	}
	var file []byte
	if param.File == nil {
		fileName := getFileName(param.FileName, param.FilePath)
		fileTmp, err := os.ReadFile(fileName)
		if err != nil {
			return 0, utils.InvaildParams.Wrap(err, "打开文件失败")
		}
		file = fileTmp
		defer os.Remove(fileName)
	} else {
		openFile, err := param.File.Open()
		if err != nil {
			return 0, utils.InvaildParams.Wrap(err, "打开文件失败")
		}
		buf := new(bytes.Buffer)
		if _, err = io.Copy(buf, openFile); err != nil {
			return 0, utils.InvaildParams.Wrap(err, "写入临时文件失败")
		}
		file = buf.Bytes()
	}
	if len(file) == 0 {
		return 0, utils.InvaildParams.New("文件格式错误")
	}
	file = utils.DecodeBytes(file)
	fileMap, files, err := utils.ReadZip(file)
	if err != nil || files == nil || len(files) == 0 {
		return 0, err
	}
	sort.Strings(files)
	var contens []string
	for _, fileName := range files {
		contens = append(contens, fileMap[fileName])
	}
	num, err := dao.ImportTool(contens)
	return num, err
}

func getFileName(fileName, filePath string) string {
	if filePath == "" {
		filePath = global.SYS_CONFIG.System.FilePath
	}
	path := fmt.Sprintf("%s.sgb", fileName)
	if filePath != "" {
		path = fmt.Sprintf("%s/%s", filePath, path)
	}
	return path
}
