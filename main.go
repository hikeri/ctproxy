package main

import (
	"flag"
	_ "github.com/sakirsensoy/genv/dotenv/autoload"
	"gitlab.roskomsvoboda.org/devops/censortracker-proxy/src"
	"log"
	"net/http"
	"strconv"
)

func main() {
	debug := flag.Bool("debug", false, "Run in verbose mode")
	port := flag.Int("port", 8888, "Port to run proxy at")
	flag.Parse()

	src.LoadLuaConfig()

	proxy := src.GetProxy()
	proxy.Verbose = *debug

	log.Println("Running proxy server on port " + strconv.Itoa(*port))

	var err error
	if src.GetConfigBool("SSL") {
		certFile, keyFile := src.GetConfig("PROXY_CERTFILE"), src.GetConfig("PROXY_KEYFILE")
		err = http.ListenAndServeTLS(":"+strconv.Itoa(*port), certFile, keyFile, proxy)
	} else {
		err = http.ListenAndServe(":"+strconv.Itoa(*port), proxy)
	}

	if err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}
