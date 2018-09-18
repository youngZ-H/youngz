package main

import (
	_ "class3/routers"
	"github.com/astaxie/beego"
	_"class3/models"

	"strconv"
)

func main() {
	beego.AddFuncMap("ShowPrepage",HandlePrepage)
	beego.AddFuncMap("ShowNextpage",HandleNextpage)
	beego.Run()
}

func HandlePrepage(data int)string  {
	pageIndex := data -1
	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}
func HandleNextpage(data int)string  {
	pageIndex := data +1
	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}