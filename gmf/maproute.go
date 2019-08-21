package gmf

import (
	"fmt"
	"net/http"
	"regexp"
)

type PathMap struct {
	path string
    urlpath map[int]Param
	params Params
}
var pathmap PathMap
type Param struct {
	Key   string
	Value string
	name string
}
type Handle func(http.ResponseWriter, *http.Request, Params)
type Params []Param

// 给Context里面的key设置变量内容
// It also lazy initializes the hashmap
func (pm *PathMap) Set(index int,key string,value string) {
pathmap.urlpath[index]=Param{
	Key:   key,
	Value: value,
}
	}
func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

func findFirstString(url string) string {
	fmt.Println("------FindString------")

	//r, _ := regexp.Compile("/(.*)/")
	r,_:=regexp.Compile("/.*/")
	s:=r.FindString(url)
	fmt.Println(r.FindString(url))
	return s
}
//TO DO
//func findParm(url string) Params{
//	uurl:=findFirstString(url)
//	r,_=regexp.Compile("（/）.*(/)")
//	for ; ;  {
//
//	}
//}

//TO DO
// func addParam{}




func (pm *Param) getValue() string {

	return pm.Value
}


func main(){
	findFirstString("/login/:name/hello")

}