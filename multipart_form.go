package webeh

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
)

var kMaxMemoryParseMultipartForm int64 = 32 << 20

type FileInfo struct {
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
	Path     string
	Ext      string
}

type MultipartFormReaderOption struct {
	MaxMemoryMB int64
}

func NewDefaultMultipartFormReaderOption() MultipartFormReaderOption {
	return MultipartFormReaderOption{
		MaxMemoryMB: 32,
	}
}

type MultipartFormReader struct {
	R *http.Request
}

func NewMultipartFormReader(r *http.Request, options ...MultipartFormReaderOption) (mr *MultipartFormReader) {
	defaultOption := NewDefaultMultipartFormReaderOption()
	for _, opt := range options {
		if opt.MaxMemoryMB > 0 {
			defaultOption.MaxMemoryMB = opt.MaxMemoryMB
		}
	}
	r.ParseMultipartForm(defaultOption.MaxMemoryMB << 20)
	mr = &MultipartFormReader{
		R: r,
	}
	return
}

func (mr *MultipartFormReader) HandleFile(field string, f func(multipart.File, *multipart.FileHeader) error) (err error) {
	file, header, e := mr.R.FormFile(field)
	if e != nil {
		err = e
		return
	} else {
		defer file.Close()
		err = f(file, header)
	}
	return
}

func (mr *MultipartFormReader) SaveFile(field string, path, name string) (fileInfo FileInfo, err error) {
	err = mr.HandleFile(field, func(mf multipart.File, mfh *multipart.FileHeader) (err error) {
		targetName := name
		ext := filepath.Ext(mfh.Filename)
		if len(targetName) > 0 {
			WebLog(mr.R, "ext:", ext)
			targetName += ext
		} else {
			targetName = mfh.Filename
		}

		uri := filepath.Join(path, targetName)
		WebLog(mr.R, "path", uri)
		os.MkdirAll(path, 0755)
		if f, e := os.Create(uri); e != nil {
			err = e
			return
		} else {
			defer f.Close()
			io.Copy(f, mf)
			fileInfo.Filename = targetName
			fileInfo.Ext = ext
			fileInfo.Path = path
			fileInfo.Header = mfh.Header
			fileInfo.Size = mfh.Size
			f.Sync()
			return
		}

	})
	return
}

func (mr *MultipartFormReader) GetString(field string) string {
	return mr.R.FormValue(field)
}
