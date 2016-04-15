package lib

import (
	"fmt"
	"testing"
)

var T *Thumber

func init() {
	cacher := NewCacher("127.0.0.1:11211", "ts_", 3600)
	T = NewThumber(MODE_CENTER, "jpg", 90, cacher)
}

func TestReadFile(t *testing.T) {
	filepath := "/Volumes/DATA/Src/Go/src/git.sillydong.com/chenzhidong/thumbservice/test/linux.jpg"
	x, err := T.ReadFile(filepath)
	fmt.Printf("%+v\n%+v\n", x, err)
}
