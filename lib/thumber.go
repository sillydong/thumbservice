package lib

import (
	"github.com/sillydong/goczd/gofile"
	"fmt"
	"path"
	"gopkg.in/gographics/imagick.v2/imagick"
	"strings"
	"os"
)

const (
	_ = iota
	TYPE_FOLLOW_WIDTH_SCALE
	TYPE_FOLLOW_HEIGHT_SCALE
	TYPE_CUT_TOP_LEFT
	TYPE_CUT_TOP_RIGHT
	TYPE_CUT_BOTTOM_LEFT
	TYPE_CUT_BOTTOM_RIGHT
	TYPE_CUT_CENTER
)

type Thumber struct {
	root   string
	cacher *Cacher
}

type ScaleConf struct {
	ImageType  string
	Width      int
	Height     int
	Proportion int
	Gary       int
	X          int
	Y          int
	Rotate     int
	Format     string
	Quality    int
}

func NewThumber(root string, cacher *Cacher) *Thumber {
	return &Thumber{root:root, cacher:cacher}
}

func (thumb *Thumber)ParseFile(request string, conf ScaleConf) ([]byte, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	request = path.Join(thumb.root, request)
	//request = strings.Replace(request, "/", os.PathSeparator, -1)
	fmt.Println(request)
	if gofile.FileExists(request) {
		return nil, fmt.Errorf("%s not found", request)
	}
	mw.ReadImage(request)
	origin := mw.GetImageBlob()
	if origin == nil {
		panic(fmt.Errorf("fail get image content: %s", request))
	}
	thumb.cacher.Put("bbb", origin)
	if err := convert(mw, conf); err != nil {
		panic(err)
	}
	result := mw.GetImageBlob()
	if result != nil {
		thumb.cacher.Put("aaa", result)
		return result, nil
	}else {
		return nil, fmt.Errorf("fail convert image: %s", request)
	}
}

func (thumb *Thumber)ParseBlob(data []byte, conf ScaleConf) ([]byte, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	mw.ReadImageBlob(data)
	if err := convert(mw, conf); err != nil {
		panic(err)
	}
	result := mw.GetImageBlob()
	if result != nil {
		thumb.cacher.Put("aaa", result)
		return result, nil
	}else {
		return nil, fmt.Errorf("fail convert image blob")
	}
}
