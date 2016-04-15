package lib

import (
	"fmt"
	"github.com/sillydong/goczd/golog"
	"gopkg.in/gographics/imagick.v2/imagick"
	"math"
	"net/url"
	"strconv"
)

const (
	_                        = iota
	MODE_CENTER              //无缩放取中间
	MODE_SCALE_CENTER_CROP   //缩放剪切中间
	MODE_SCALE_CENTER_INSIDE //缩放补全中间
	MODE_FIT_XY              //缩放到指定尺寸
)

type ScaleConf struct {
	ScaleMode int
	Width     uint
	Height    uint
	Rotate    int
	Format    string
	Quality   uint
}

func NewScaleConf(params url.Values) *ScaleConf {
	conf := &ScaleConf{}
	conf.ScaleMode, _ = strconv.Atoi(params.Get("m"))
	w, _ := strconv.Atoi(params.Get("w"))
	conf.Width = uint(w)
	h, _ := strconv.Atoi(params.Get("h"))
	conf.Height = uint(h)
	conf.Rotate, _ = strconv.Atoi(params.Get("r"))
	conf.Format = params.Get("f")
	q, _ := strconv.Atoi(params.Get("q"))
	conf.Quality = uint(q)

	return conf
}

func convert(mw *imagick.MagickWand, conf *ScaleConf) error {
	golog.Debugf("call convert function %+v", conf)

	if err := mw.SetImageOrientation(imagick.ORIENTATION_TOP_LEFT); err != nil {
		return err
	}

	//清理数据
	if err := mw.StripImage(); err != nil {
		return err
	}

	//旋转
	if conf.Rotate > 0 {
		background := imagick.NewPixelWand()
		if background == nil {
			return fmt.Errorf("init new pixelwand faile.")
		}
		defer background.Destroy()
		isOk := background.SetColor("#ffffff")
		if !isOk {
			return fmt.Errorf("set background color faile.")
		}

		if err := mw.RotateImage(background, float64(conf.Rotate)); err != nil {
			return err
		}
	}

	//裁剪
	imagewidth := mw.GetImageWidth()
	imageheight := mw.GetImageHeight()
	if (conf.Width > 0 || conf.Height > 0) && (conf.Width != imagewidth || conf.Height != imageheight) {
		switch conf.ScaleMode {
		case MODE_CENTER:
			if conf.Width == 0 && conf.Height > 0 {
				conf.Width = uint(round((float64(conf.Height) * float64(imagewidth) / float64(imageheight))))
			} else if conf.Width > 0 && conf.Height == 0 {
				conf.Height = uint(round(float64(conf.Width) * float64(imageheight) / float64(imagewidth)))
			}

			x := int(round((float64(imagewidth) - float64(conf.Width)) / 2))
			y := int(round((float64(imageheight) - float64(conf.Height)) / 2))
			golog.Debugf("1 crop: w:%v\th:%v\tx:%v\ty:%v", conf.Width, conf.Height, x, y)
			if err := mw.CropImage(conf.Width, conf.Height, x, y); err != nil {
				return err
			}
		case MODE_SCALE_CENTER_CROP:
			var scale float64
			if conf.Width > 0 && conf.Height > 0 {
				scale = math.Min(float64(imagewidth)/float64(conf.Width), float64(imageheight)/float64(conf.Height))
			} else if conf.Width == 0 && conf.Height > 0 {
				scale = float64(imageheight) / float64(conf.Height)
			} else if conf.Width > 0 && conf.Height == 0 {
				scale = float64(imagewidth) / float64(conf.Width)
			}

			scalewidth := uint(round(float64(imagewidth) / scale))
			if conf.Width == 0 {
				conf.Width = scalewidth
			}
			scaleheight := uint(round(float64(imageheight) / scale))
			if conf.Height == 0 {
				conf.Height = scaleheight
			}
			golog.Debugf("2 thumb: w:%v\th:%v", scalewidth, scaleheight)
			if err := mw.ThumbnailImage(scalewidth, scaleheight); err != nil {
				return err
			}

			x := int(round((float64(scalewidth) - float64(conf.Width)) / 2))
			y := int(round((float64(scaleheight) - float64(conf.Height)) / 2))
			golog.Debugf("2 crop: w:%v\th:%v\tx:%v\ty:%v", conf.Width, conf.Height, x, y)
			if err := mw.CropImage(conf.Width, conf.Height, x, y); err != nil {
				return err
			}
		case MODE_SCALE_CENTER_INSIDE:
			var scale float64
			if conf.Width > 0 && conf.Height > 0 {
				scale = math.Max(float64(imagewidth)/float64(conf.Width), float64(imageheight)/float64(conf.Height))
			} else if conf.Width == 0 && conf.Height > 0 {
				scale = float64(imageheight) / float64(conf.Height)
			} else if conf.Width > 0 && conf.Height == 0 {
				scale = float64(imagewidth) / float64(conf.Width)
			}

			scalewidth := uint(round(float64(imagewidth) / scale))
			if conf.Width == 0 {
				conf.Width = scalewidth
			}
			scaleheight := uint(round(float64(imageheight) / scale))
			if conf.Height == 0 {
				conf.Height = scaleheight
			}
			golog.Debugf("3 thumb: w:%v\th:%v", scalewidth, scaleheight)
			if err := mw.ThumbnailImage(scalewidth, scaleheight); err != nil {
				return err
			}

			x := int(round((float64(conf.Width) - float64(scalewidth)) / 2))
			y := int(round((float64(scaleheight) - float64(conf.Height)) / 2))
			golog.Debugf("3 crop: w:%v\th:%v\tx:%v\ty:%v", conf.Width, conf.Height, x, y)
			if err := mw.ExtentImage(conf.Width, conf.Height, -x, y); err != nil {
				return err
			}
		case MODE_FIT_XY:
			if conf.Width == 0 && conf.Height > 0 {
				conf.Width = uint(round((float64(conf.Height) * float64(imagewidth) / float64(imageheight))))
			} else if conf.Width > 0 && conf.Height == 0 {
				conf.Height = uint(round(float64(conf.Width) * float64(imageheight) / float64(imagewidth)))
			}
			golog.Debugf("4 thumb: w:%v\th:%v", conf.Width, conf.Height)
			if err := mw.ThumbnailImage(uint(conf.Width), uint(conf.Height)); err != nil {
				return err
			}
		}
	}

	//压缩
	if conf.Quality != 100 {
		if err := mw.SetCompressionQuality(conf.Quality); err != nil {
			return err
		}
	}

	//调整格式
	if conf.Format != mw.GetImageFormat() {
		if err := mw.SetImageFormat(conf.Format); err != nil {
			return err
		}
	}

	return nil
}

func round(val float64) float64 {
	if val > 0.0 {
		return math.Floor(val + 0.5)
	} else {
		return math.Ceil(val - 0.5)
	}
}
