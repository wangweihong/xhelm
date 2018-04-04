package main

import "xlog"

func main() {
	err := xlog.Init()
	if err != nil {
		panic(err.Error())
	}
	xlog.Logger.Info("start to run xhelm")
}
