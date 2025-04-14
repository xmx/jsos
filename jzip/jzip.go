package jzip

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/fs"
)

func OpenFile(filepath string) (*JZip, error) {
	zrc, err := zip.OpenReader(filepath)
	if err != nil {
		return nil, err
	}

	return Open(zrc)
}

func Open(zrc *zip.ReadCloser) (*JZip, error) {
	manifestFile, err := zrc.Open("manifest.json")
	if err != nil {
		_ = zrc.Close()
		return nil, fmt.Errorf("缺少 manifest.json: %w", err)
	}
	defer manifestFile.Close()

	var manifest Manifest
	if err = json.NewDecoder(manifestFile).Decode(&manifest); err != nil {
		_ = zrc.Close()
		return nil, fmt.Errorf("解析 manifest.json 错误: %w", err)
	}
	app := manifest.Application
	if app.ID == "" {
		_ = zrc.Close()
		return nil, fmt.Errorf("软件ID不能为空")
	}
	if app.Name == "" {
		_ = zrc.Close()
		return nil, fmt.Errorf("软件名不能为空")
	}

	return &JZip{Manifest: manifest, zrc: zrc}, nil
}

type Manifest struct {
	Version     int         `json:"version"`
	Application Application `json:"application"`
}

type Application struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Main    string `json:"main"`
	Version string `json:"version"`
}

type JZip struct {
	Manifest Manifest
	zrc      *zip.ReadCloser
}

func (jz JZip) Open(name string) (fs.File, error) {
	return jz.zrc.Open(name)
}

func (jz JZip) Close() error {
	return jz.zrc.Close()
}
