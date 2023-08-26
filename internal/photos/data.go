package photos

import (
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/keyword-search-backend/entity"
)

type FileInfo struct {
	MimType   string
	Size      int64
	Path      string
	Extension string
}

type Request struct {
	restapi.RequestDTO[entity.Photo]
	Caption   string   `json:"caption"`
	FileInfo  FileInfo `json:"file_info"`
	GalleryId uint
}

func (pr *Request) ToEntity(ent *entity.Photo) error {
	ent.Caption = pr.Caption
	ent.GalleryId = pr.GalleryId
	ent.File = entity.File{
		MimeType:  pr.FileInfo.MimType,
		Size:      pr.FileInfo.Size,
		Path:      pr.FileInfo.Path,
		Extension: pr.FileInfo.Extension,
	}
	return nil
}

type Response struct {
	restapi.ResponseDTO[entity.Photo]
	Caption   string   `json:"caption"`
	GalleryId uint     `json:"gallery_id"`
	FileInfo  FileInfo `json:"file_info"`
}

func (pr *Response) FromEntity(ent entity.Photo) error {
	pr.FileInfo = FileInfo{
		Path:      ent.File.Path,
		MimType:   ent.File.MimeType,
		Extension: ent.File.Extension,
		Size:      ent.File.Size,
	}
	pr.Caption = ent.Caption
	pr.GalleryId = ent.GalleryId

	return nil
}
