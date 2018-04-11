package main

import (
	_ "xhelm/charts"
	_ "xhelm/repository"
	_ "xhelm/routers"
	"xhelm/xlog"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	err := xlog.Init()
	if err != nil {
		panic(err.Error())
	}
	xlog.Logger.Info("start to run xhelm")

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	logs.EnableFuncCallDepth(true)
	/*
		else {
			if err := os.Mkdir(filepath.Dir("./logs/"), 0755); err != nil {
				if !os.IsExist(err) {
					panic(err.Error())
				}
			}
			logs.EnableFuncCallDepth(true)
			logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/deploy.log", "perm":"0664"}`)
		}
	*/
	beego.Run()
}
