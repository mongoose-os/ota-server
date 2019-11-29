package main

import (
	"flag"
	"net/http"
	"path/filepath"
	"time"

	"github.com/golang/glog"

	"github.com/mognoose-os/ota-server/handler"
)

var (
	flagRoot          = flag.String("root", "", "directory to serve")
	flagListenAddr    = flag.String("listen-addr", "", "address:port to listen on")
	flagTLSListenAddr = flag.String("tls-listen-addr", "", "address:port to listen on")
	flagTLSCert       = flag.String("tls-cert", "", "TLS certificate file")
	flagTLSKey        = flag.String("tls-key", "", "TLS key file")
	flagChunkSize     = flag.Int("chunk-size", 1024, "Serving chunk size, bytes")
	flagDelay         = flag.Duration("delay", 0, "Inter-chunk delay")
)

func main() {
	flag.Parse()
	if *flagRoot == "" {
		glog.Exitf("--root is required")
	}
	if *flagListenAddr == "" && *flagTLSListenAddr == "" {
		glog.Exitf("at least one of --listen-addr, --tls-listen-addr is required")
	}
	var err error
	var absRoot string
	if absRoot, err = filepath.Abs(*flagRoot); err != nil {
		glog.Exitf("Invalid --root %q: %s", *flagRoot, err)
	}
	glog.Infof("Root: %s", absRoot)
	handler := handler.NewOTAFileHandler(absRoot, *flagChunkSize, *flagDelay)
	if *flagListenAddr != "" {
		go func() {
			glog.Infof("Plain listener: %s", *flagListenAddr)
			if err := http.ListenAndServe(*flagListenAddr, handler); err != nil {
				glog.Exitf("ListenAndServe: %s", err)
			}
		}()
	}
	if *flagTLSListenAddr != "" {
		if *flagTLSCert == "" {
			glog.Exitf("at least one of --listen-addr, --tls-listen-addr is required")
		}
		go func() {
			time.Sleep(10 * time.Millisecond)
			glog.Infof("TLS listener: %s, cert: %s, key: %s", *flagTLSListenAddr, *flagTLSCert, *flagTLSKey)
			if err := http.ListenAndServeTLS(*flagTLSListenAddr, *flagTLSCert, *flagTLSKey, handler); err != nil {
				glog.Exitf("ListenAndServeTLS: %s", err)
			}
		}()
	}
	<-make(chan struct{}) // sleep forever
}
