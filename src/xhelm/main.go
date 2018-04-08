package main

import (
	_ "charts"
	_ "repository"
	"xlog"
)

func main() {
	err := xlog.Init()
	if err != nil {
		panic(err.Error())
	}
	xlog.Logger.Info("start to run xhelm")
}
