package main

import (
	"net/http"
	"os"
	"github.com/codegangsta/cli"
	"github.com/sillydong/goczd/gofile"
	"fmt"
	"github.com/sillydong/goczd/goconf"
	"git.sillydong.com/chenzhidong/thumbservice/lib"
)

const VERSION = "1.0"

var Service ServiceConf

type ServiceConf struct {
	Root    string
	Cacher  *lib.Cacher
	Thumber *lib.Thumber
}

func main() {
	app := cli.NewApp();
	app.Name = "图片缩略图服务"
	app.Usage = "不影响原文件系统布局情况下生成图片指定尺寸及清晰度的缩略图"
	app.Author = "陈志东"
	app.Copyright = "http://sillydong.com"
	app.Version = VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:"c",
			Value: "",
			Usage:"set conf file for current service",
		},
	}
	app.Action = run
	app.Run(os.Args)
}

func run(ctx *cli.Context) {
	if ctx.String("c") == "" {
		cli.ShowAppHelp(ctx)
	}else {
		confpath := ctx.String("c")
		if !gofile.FileExists(confpath) {
			fmt.Printf("%s not exists\n", confpath)
		}else {
			conf, err := goconf.ReadConfigFile(confpath)
			if err != nil {
				panic(err)
			}
			listen, _ := conf.GetString("default", "listen")
			if len(listen) == 0 {
				panic("not sure which ip:port to listen")
			}
			root, _ := conf.GetString("default", "root")
			if len(root) == 0 {
				panic("root directory not set")
			}
			cache, _ := conf.GetString("default", "cache")
			if len(cache) == 0 {
				panic("cache not set")
			}
			prefix, _ := conf.GetString("default", "prefix")
			if len(prefix) == 0 {
				prefix = "ts_"
			}
			expires, _ := conf.GetInt("default", "expires")
			if expires == 0 {
				expires = 86400
			}

			cacher := lib.NewCacher(cache, prefix, int32(expires))
			thumber := lib.NewThumber(root, cacher)

			Service = ServiceConf{
				Root:root,
				Cacher:cacher,
				Thumber:thumber,
			}

			http.HandleFunc("/", handler)
			if err := http.ListenAndServe(listen, nil); err != nil {
				panic(err)
			}
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	thumbkey := "aaa"
	//缓存的缩略图
	data, _ := Service.Cacher.Get(thumbkey)
	if data == nil {
		originkey := "bbb"
		var err error
		//缓存的原图
		origin, _ := Service.Cacher.Get(originkey)
		if origin == nil {
			//操作原图文件
			data, err = Service.Thumber.ParseFile("ccc",lib.ScaleConf{})
		}else {
			//操作原图缓存
			data, err = Service.Thumber.ParseBlob(origin,lib.ScaleConf{})
		}
		if err != nil {
			//有错误
			fmt.Printf("%+v\n", err)
		}
	}
	w.Write(data)

}
