package main

import (
	"encoding/json"
	"fmt"
	"git.sillydong.com/chenzhidong/thumbservice/lib"
	"github.com/codegangsta/cli"
	"github.com/sillydong/goczd/goconf"
	"github.com/sillydong/goczd/gofile"
	"github.com/sillydong/goczd/golog"
	"gopkg.in/gographics/imagick.v2/imagick"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const VERSION = "1.0"

var Service ServiceConf

type ServiceConf struct {
	Root    string
	Cacher  *lib.Cacher
	Thumber *lib.Thumber
}

func main() {
	app := cli.NewApp()
	app.Name = "ThubmService"
	app.Usage = "Service to resize images"
	app.Author = "Chen.Zhidong"
	app.Copyright = "http://sillydong.com"
	app.Version = VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c",
			Value: "",
			Usage: "set conf file for current service",
		},
	}
	app.Action = run
	app.Run(os.Args)
}

func run(ctx *cli.Context) {
	if ctx.String("c") == "" {
		cli.ShowAppHelp(ctx)
	} else {
		confpath := ctx.String("c")
		if !gofile.FileExists(confpath) {
			fmt.Printf("%s not exists\n", confpath)
		} else {
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
			cache, _ := conf.GetString("default", "memcache")
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
			defaultmode, _ := conf.GetInt("scale", "mode")
			if defaultmode == 0 {
				defaultmode = lib.MODE_CENTER
			}
			defaultformat, _ := conf.GetString("scale", "format")
			if defaultformat == "" {
				defaultformat = "jpg"
			}
			defaultquality, _ := conf.GetInt("scale", "quality")
			if defaultquality == 0 {
				defaultquality = 90
			}

			logfilename, _ := conf.GetString("log", "filename")
			if logfilename == "" {
				logfilename = "thumb.log"
			}
			loglevel, _ := conf.GetInt("log", "level")
			if loglevel == 0 {
				loglevel = golog.LevelError
			}
			logmaxdays, _ := conf.GetInt("log", "maxdays")
			if logmaxdays == 0 {
				logmaxdays = 7
			}

			logconf := map[string]interface{}{
				"filename": logfilename,
				"daily":    true,
				"maxdays":  logmaxdays,
				"rotate":   true,
				"level":    loglevel,
			}
			logstr, _ := json.Marshal(logconf)
			golog.SetLogger("file", string(logstr))

			cacher := lib.NewCacher(cache, prefix, int32(expires))

			thumber := lib.NewThumber(defaultmode, defaultformat, defaultquality, cacher)

			Service = ServiceConf{
				Root:    root,
				Cacher:  cacher,
				Thumber: thumber,
			}

			http.HandleFunc("/", handler)
			if err := http.ListenAndServe(listen, nil); err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("service running...")
			}
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	imagick.Initialize()
	defer imagick.Terminate()

	orig_path := r.URL.Path
	thumb_path := r.RequestURI
	scale := r.URL.Query()

	orig_path = path.Join(Service.Root, orig_path)
	thumb_path = path.Join(Service.Root, thumb_path)

	var err error
	var route []string
	var status int
	var data []byte
	start := time.Now().UnixNano()
	if len(orig_path) <= 1 || len(thumb_path) <= 1 {
		route = append(route, "err")
		status = http.StatusForbidden
	} else if orig_path == "/favicon.ico" {
		route = append(route, "fav")
		status = http.StatusNotFound
	} else if len(scale) == 0 {
		route = append(route, "or")
		//取原图
		data, _ = Service.Cacher.Get(orig_path)
		if data == nil || len(data) == 0 {
			route = append(route, "go")
			//将原图读取到memcache
			if gofile.FileExists(orig_path) {
				data, err = Service.Thumber.ReadFile(orig_path)
				if err != nil {
					//有错误
					golog.Errorf("or-go [%v]%v", err, orig_path)
					status = http.StatusInternalServerError
				} else {
					status = http.StatusOK
				}
			} else {
				status = http.StatusNotFound
			}
		} else {
			route = append(route, "gc")
			//操作原图缓存
			status = http.StatusOK
		}
	} else {
		route = append(route, "th")
		//缓存的缩略图
		data, _ = Service.Cacher.Get(thumb_path)
		if data == nil || len(data) == 0 {
			route = append(route, "gto")

			//缓存的原图
			scaleconf := lib.NewScaleConf(scale)
			if scaleconf.ScaleMode <= 0 || scaleconf.ScaleMode > 4 {
				scaleconf.ScaleMode = Service.Thumber.DefaultMode
			}
			if scaleconf.Format == "" {
				scaleconf.Format = Service.Thumber.DefaultFormat
			}
			if scaleconf.Quality <= 0 || scaleconf.Quality > 90 {
				scaleconf.Quality = Service.Thumber.DefaultQuality
			}
			origin, _ := Service.Cacher.Get(orig_path)
			if origin == nil || len(origin) == 0 {
				route = append(route, "go")
				//操作原图文件

				if gofile.FileExists(orig_path) {
					data, err = Service.Thumber.ParseFile(orig_path, thumb_path, scaleconf)
					if err != nil {
						//有错误
						golog.Errorf("th-go [%v]%v", err, orig_path)
						status = http.StatusInternalServerError
					} else {
						status = http.StatusOK
					}
				} else {
					status = http.StatusNotFound
				}
			} else {
				route = append(route, "gc")
				//操作原图缓存
				data, err = Service.Thumber.ParseBlob(origin, thumb_path, scaleconf)
				if err != nil {
					//有错误
					golog.Errorf("th-gc [%v]%v", err, orig_path)
					status = http.StatusInternalServerError
				} else {
					status = http.StatusOK
				}
			}
		} else {
			route = append(route, "gtc")
			status = http.StatusOK
		}
	}
	route = append(route, strconv.FormatInt(time.Now().UnixNano()-start, 10))
	w.Header().Set("Service-Log", strings.Join(route, "-"))
	w.WriteHeader(status)
	w.Write(data)
}
