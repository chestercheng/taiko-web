package main

import (
	"flag"
	"os"
	"taiko-web/config"
	"taiko-web/db"
	"taiko-web/web"
)

var (
	h    bool
	mode string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help.")
	flag.StringVar(
		&mode, "mode", "debug", "sets gin mode according to input string.")
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		os.Exit(1)
	}

	config.Init(mode)
	conf := config.GetConfig()
	db.Init(conf)
	web.Init(conf)
}
