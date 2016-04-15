package lib

import (
	"fmt"
	"github.com/sillydong/goczd/gofile"
	"gopkg.in/gographics/imagick.v2/imagick"
	"net/url"
	"testing"
)

func TestScaleConf(t *testing.T) {
	patha := "/linux.jpg"
	urla, _ := url.ParseRequestURI(patha)
	fmt.Printf("%+v\n", NewScaleConf(urla.Query()))
	pathb := "/linux.jpg?"
	urlb, _ := url.ParseRequestURI(pathb)
	fmt.Printf("%+v\n", NewScaleConf(urlb.Query()))
	pathc := "/linux.jpg?m=1&w=100&h=100"
	urlc, _ := url.ParseRequestURI(pathc)
	fmt.Printf("%+v\n", NewScaleConf(urlc.Query()))
	pathd := "/linux.jpg?m=1&w=100"
	urld, _ := url.ParseRequestURI(pathd)
	fmt.Printf("%+v\n", NewScaleConf(urld.Query()))
	pathe := "/linux.jpg?m=1&w=100&h=100&r=90"
	urle, _ := url.ParseRequestURI(pathe)
	fmt.Printf("%+v\n", NewScaleConf(urle.Query()))
	pathf := "/linux.jpg?m=1&w=100&h=100&f=png"
	urlf, _ := url.ParseRequestURI(pathf)
	fmt.Printf("%+v\n", NewScaleConf(urlf.Query()))
	pathg := "/linux.jpg?m=1&w=100&h=100&q=80"
	urlg, _ := url.ParseRequestURI(pathg)
	fmt.Printf("%+v\n", NewScaleConf(urlg.Query()))
	pathh := "/linux.jpg?m=1&w=100&h=100&r=&f=&q=80"
	urlh, _ := url.ParseRequestURI(pathh)
	fmt.Printf("%+v\n", NewScaleConf(urlh.Query()))
	pathi := "/linux.jpg?m=1&w=100&h=100&f&q=80"
	urli, _ := url.ParseRequestURI(pathi)
	fmt.Printf("%+v\n", NewScaleConf(urli.Query()))
	pathj := "/linux.jpg?w=100&h=100&r=90&f=png&q=80"
	urlj, _ := url.ParseRequestURI(pathj)
	fmt.Printf("%+v\n", NewScaleConf(urlj.Query()))
}

func TestConvert(t *testing.T) {
	dir, _ := gofile.WorkDir()
	fmt.Println(dir)
	filepath := "/Volumes/DATA/Src/Go/src/git.sillydong.com/chenzhidong/thumbservice/test/linux.jpg"

	uris := []string{
		//"/linux.jpg?m=1&w=100&h=100",
		//"/linux.jpg?m=1&w=100",
		//"/linux.jpg?m=1&h=100",
		//"/linux.jpg?m=1&h=100&r=90",
		//"/linux.jpg?m=1&w=100&h=100&f=png",
		//"/linux.jpg?m=1&w=100&h=100&q=80",
		//"/linux.jpg?m=2&w=100&h=100",
		//"/linux.jpg?m=2&w=100",
		//"/linux.jpg?m=2&h=100",
		//"/linux.jpg?m=2&h=100&r=90",
		//"/linux.jpg?m=2&w=100&h=100&f=png",
		//"/linux.jpg?m=2&w=100&h=100&q=80",
		"/linux.jpg?m=3&w=100&h=100",
		"/linux.jpg?m=3&w=100",
		"/linux.jpg?m=3&h=100",
		"/linux.jpg?m=3&h=100&r=90",
		"/linux.jpg?m=3&w=100&h=100&r=90",
		"/linux.jpg?m=3&w=100&h=100&f=png",
		"/linux.jpg?m=3&w=100&h=100&q=80",
		//"/linux.jpg?m=4&w=100&h=100",
		//"/linux.jpg?m=4&w=100",
		//"/linux.jpg?m=4&h=100",
		//"/linux.jpg?m=4&h=100&r=90",
		//"/linux.jpg?m=4&w=100&h=100&f=png",
		//"/linux.jpg?m=4&w=100&h=100&q=80",
	}
	for index, uri := range uris {
		parse(filepath, index, uri)
	}
}

func parse(filepath string, index int, uri string) {
	fmt.Println(index)
	fmt.Println(uri)
	params, _ := url.ParseRequestURI(uri)
	fmt.Printf("%+v\n", params)
	conf := NewScaleConf(params.Query())

	if conf.ScaleMode <= 0 || conf.ScaleMode > 4 {
		conf.ScaleMode = MODE_CENTER
	}
	if conf.Format == "" {
		conf.Format = "jpg"
	}
	if conf.Quality <= 0 || conf.Quality > 90 {
		conf.Quality = 90
	}

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	mw.ReadImage(filepath)
	origin := mw.GetImageBlob()
	if origin == nil {
		fmt.Printf("fail get image content: %s", filepath)
	}
	T.cacher.Put(filepath, origin)
	fmt.Printf("w:%v\th:%v\n", mw.GetImageWidth(), mw.GetImageHeight())
	if err := convert(mw, conf); err != nil {
		fmt.Printf("%+v\n", err)
	} else {
		fmt.Printf("%v\n", mw.GetImageBlob())
		//err:=mw.WriteImage("/Volumes/DATA/Src/Go/src/git.sillydong.com/chenzhidong/thumbservice/lib/result_" + strconv.Itoa(index) + "." + conf.Format)
		//if err!=nil{
		//	fmt.Printf("%+v\n",err)
		//}else{
		//	fmt.Printf("w:%v\th:%v\n", mw.GetImageWidth(), mw.GetImageHeight())
		//	fmt.Println("convert finish")
		//}
	}
}
