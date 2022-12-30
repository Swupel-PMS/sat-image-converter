package converter

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"sync"
)

const (
	contentType  = "contentType"
	zipType      = "application/zip"
	infoFileName = "info.json"
)

type Result struct {
	ContentType string
	Images      []Image
}

type jsonImage struct {
	Name string `json:"name"`
	Image
	ErrorMsg jsonError `json:"error"`
}

type jsonResult struct {
	Images []jsonImage `json:"images"`
}

type jsonError string

func (r Result) ParseResponse(ctx context.Context, writer http.ResponseWriter) error {
	if len(r.Images) == 0 {
		return nil
	}
	if len(r.Images) == 1 {
		return r.asFile(writer)
	}
	return r.asZip(ctx, writer)
}

func (r Result) getFileExtension() string {
	mimeTypes, err := mime.ExtensionsByType(r.ContentType)
	if err != nil {
		return ""
	}
	return mimeTypes[0]

}

func (r Result) asZip(ctx context.Context, writer http.ResponseWriter) error {
	// create in memory zip builder
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	mimeType := r.getFileExtension()
	wg := &sync.WaitGroup{}
	mutexMetaData := &sync.Mutex{}
	mutexWriting := &sync.Mutex{}
	metaData := make([]jsonImage, 0)
	for i, im := range r.Images {
		select {
		case <-ctx.Done():
			zipWriter.Close()
			return ctx.Err()
		default:
			wg.Add(1)
			go func(index int, img Image) {
				defer wg.Done()
				errMsg := ""
				if img.Error != nil {
					errMsg = img.Error.Error()
				}
				jI := jsonImage{
					Name:     fmt.Sprintf("%d%s", index, mimeType),
					Image:    img,
					ErrorMsg: jsonError(errMsg),
				}
				defer func() {
					mutexMetaData.Lock()
					metaData = append(metaData, jI)
					mutexMetaData.Unlock()
				}()
				if jI.Error != nil {
					return
				}
				mutexWriting.Lock()
				defer mutexWriting.Unlock()
				f, err := zipWriter.Create(jI.Name)
				if err != nil {
					img.Error = err
					return
				}
				_, err = io.Copy(f, jI.Data)
				if err != nil {
					img.Error = err
					return
				}
			}(i, im)
		}
	}
	wg.Wait()
	// create metadata file and safe
	info, err := zipWriter.Create(infoFileName)
	if err == nil {
		// if info file can not be create do not even try to create the other ones
		file, _ := json.MarshalIndent(&jsonResult{Images: metaData}, "", " ")
		_, err = info.Write(file)
		if err != nil {
			return err
		}
	}
	err = zipWriter.Close()
	if err != nil {
		return err
	}
	writer.Header().Set(contentType, zipType)
	_, err = io.Copy(writer, buf)
	return err
}

func (r Result) asFile(writer http.ResponseWriter) error {
	img := r.Images[0]
	if img.Error != nil {
		return img.Error
	}
	writer.Header().Add(contentType, r.ContentType)
	_, err := io.Copy(writer, img.Data)
	if err != nil {
		return err
	}
	return nil
}
