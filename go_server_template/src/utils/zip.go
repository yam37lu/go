package utils

import (
	"archive/zip"
	"bytes"
	"strconv"
)

func CreateZip(files map[string][]byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	for name, content := range files {
		ioWriter, err := zipWriter.Create(name)
		if err != nil {
			return nil, err
		}
		_, err = ioWriter.Write(content)
		if err != nil {
			return nil, err
		}
	}
	zipWriter.Close()
	return buf.Bytes(), nil

}

func ReadZip(file []byte) (map[string]string, []string, error) {
	reader := bytes.NewReader(file)
	zipReader, err := zip.NewReader(reader, reader.Size())
	if err != nil {
		return nil, nil, err
	}
	files := make([]string, 0)
	result := make(map[string]string)
	for _, zipFile := range zipReader.File {
		readCloser, err := zipFile.Open()
		if err != nil {
			return nil, nil, err
		}
		defer readCloser.Close()
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(readCloser)
		if err != nil {
			return nil, nil, err
		}
		content := buf.String()
		if content == "" {
			continue
		}
		result[zipFile.Name] = content
		files = append(files, zipFile.Name)
	}
	return result, files, err
}

func EncodeBytes(content, key []byte) []byte {
	segALen := len(content) / 3
	offsetNum := len(content) % 3
	offset := []byte(strconv.Itoa(offsetNum))
	segBLen := segALen
	segA := append(key[0:8], content[0:segALen]...)
	segB := append(key[8:16], content[segALen:segALen+segBLen]...)
	segC := append(key[16:24], content[segALen+segBLen:]...)
	result := append(offset, segC...)
	result = append(result, segB...)
	result = append(result, segA...)
	return result
}

func DecodeBytes(src []byte) []byte {
	offsetNum, _ := strconv.Atoi(string(src[0]))
	segALen := (len(src) - 1 - 24) / 3
	segCLen := 9 + segALen + offsetNum
	segC := src[9:segCLen]
	segB := src[segCLen+8 : segCLen+8+segALen]
	segA := src[segCLen+8+segALen+8:]
	result := append(segA, segB...)
	result = append(result, segC...)
	return result
}
