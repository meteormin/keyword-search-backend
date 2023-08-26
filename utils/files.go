package utils

import (
	"mime/multipart"
	"net/http"
	"strings"
)

type FileInfo struct {
	MimeType  string
	Size      int64
	Extension string
}

type sizer interface {
	Size() int64
}

func GetFileInfo(file *multipart.FileHeader) (*FileInfo, error) {
	split := strings.Split(file.Filename, ".")
	ext := ""
	if len(split) > 1 {
		ext = split[1]
	}

	f, err := file.Open()
	if err != nil {
		return nil, err
	}

	fh := make([]byte, 512)
	_, err = f.Read(fh)
	if err != nil {
		return nil, err
	}

	mimType := http.DetectContentType(fh)

	return &FileInfo{
		MimeType:  mimType,
		Size:      f.(sizer).Size(),
		Extension: ext,
	}, nil
}
