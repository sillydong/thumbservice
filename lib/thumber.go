package lib

import (
	"fmt"
	"github.com/sillydong/goczd/golog"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type Thumber struct {
	DefaultMode    int
	DefaultFormat  string
	DefaultQuality uint
	cacher         *Cacher
}

func NewThumber(defaultmode int, defaultformat string, defaultquality int, cacher *Cacher) *Thumber {
	return &Thumber{DefaultMode: defaultmode, DefaultFormat: defaultformat, DefaultQuality: uint(defaultquality), cacher: cacher}
}

func (thumb *Thumber) ReadFile(orig_path string) ([]byte, error) {
	golog.Debugf("read file %v\n", orig_path)
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	mw.ReadImage(orig_path)
	mw.ResetIterator()
	origin := mw.GetImageBlob()
	if origin == nil || len(origin) == 0 {
		return nil, fmt.Errorf("fail get image content: %s", orig_path)
	}
	thumb.cacher.Put(orig_path, origin)
	return origin, nil
}

func (thumb *Thumber) ParseFile(orig_path, thumb_path string, conf *ScaleConf) ([]byte, error) {
	golog.Debugf("parse file %v", orig_path)
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	mw.ReadImage(orig_path)
	origin := mw.GetImageBlob()
	if origin == nil || len(origin) == 0 {
		return nil, fmt.Errorf("fail get image content: %s", orig_path)
	}
	golog.Debugf("cache origin:%v", len(origin))
	thumb.cacher.Put(orig_path, origin)
	if err := convert(mw, conf); err != nil {
		return nil, err
	} else {
		result := mw.GetImageBlob()
		if result != nil && len(result) > 0 {
			golog.Debugf("cache thumb:%v", len(result))
			thumb.cacher.Put(thumb_path, result)
			return result, nil
		} else {
			return nil, fmt.Errorf("fail convert image: %s", orig_path)
		}
	}
}

func (thumb *Thumber) ParseBlob(data []byte, thumb_path string, conf *ScaleConf) ([]byte, error) {
	golog.Debugf("parse blob %v", len(data))
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	mw.ReadImageBlob(data)
	if err := convert(mw, conf); err != nil {
		return nil, err
	} else {
		mw.ResetIterator()
		result := mw.GetImageBlob()
		if result != nil && len(result) > 0 {
			golog.Debugf("cache thumb:%v", len(result))
			thumb.cacher.Put(thumb_path, result)
			return result, nil
		} else {
			return nil, fmt.Errorf("fail convert image blob")
		}
	}
}
