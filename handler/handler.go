package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/mongoose-os/mos/common/fwbundle"

	"github.com/mognoose-os/ota-server/flushing_writer"
)

const (
	deviceIDHeader  = "X-MGOS-Device-ID"
	fwVersionheader = "X-MGOS-FW-Version"
)

type RequestInfo struct {
	DeviceID   string
	DeviceMAC  string
	DeviceArch string
	FWVersion  string
	FWBuildID  string
}

type OTAFileHandler struct {
	root      string
	chunkSize int
	delay     time.Duration
}

func (h *OTAFileHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasSuffix(req.URL.Path, ".zip") {
		h.serveFirmware(w, req)
	} else {
		h.serveOther(w, req)
	}
}

func (h *OTAFileHandler) serveFirmware(w http.ResponseWriter, req *http.Request) {
	ri := &RequestInfo{}
	parts := strings.Split(req.Header.Get(deviceIDHeader), " ")
	ri.DeviceID = parts[0]
	if len(parts) > 1 {
		ri.DeviceMAC = parts[1]
	}
	parts = strings.Split(req.Header.Get(fwVersionheader), " ")
	ri.DeviceArch = parts[0]
	if len(parts) > 1 {
		ri.FWVersion = parts[1]
	}
	if len(parts) > 2 {
		ri.FWBuildID = parts[2]
	}
	// Path is sanitized already, so it's fine.
	fname, err := filepath.Abs(filepath.Join(h.root, req.URL.Path))
	if err != nil {
		sendError(w, req, ri, "What is this?", http.StatusBadRequest)
		return
	}
	fwVersion, fwBuildId := "", ""
	if ri.FWVersion != "" {
		fwb, err := fwbundle.ReadZipFirmwareBundle(fname)
		if err != nil {
			sendError(w, req, ri, "Not there", http.StatusNotFound)
			return
		}
		fwVersion, fwBuildId = fwb.Version, fwb.BuildID
		if fwVersion == ri.FWVersion {
			sendError(w, req, ri, "Not Modified", http.StatusNotModified)
			return
		}
	}
	glog.Infof("%s %s -> [%s %s]", req.RemoteAddr, ri, fwVersion, fwBuildId)
	http.ServeFile(flushing_writer.NewFlushingResponseWriter(w.(flushing_writer.FlushingResponseWriter), h.chunkSize, h.delay), req, fname)
}

func (h *OTAFileHandler) serveOther(w http.ResponseWriter, req *http.Request) {
	// Path is sanitized already, so it's fine.
	fname := filepath.Join(h.root, req.URL.Path)
	http.ServeFile(w, req, fname)
}

func sendError(w http.ResponseWriter, req *http.Request, ri *RequestInfo, message string, code int) {
	glog.Infof("%s %s -> %d %s", req.RemoteAddr, ri, code, message)
	http.Error(w, message, code)
}

func (ri *RequestInfo) String() string {
	parts := []string{}
	if ri.DeviceID != "" {
		parts = append(parts, ri.DeviceID)
	}
	if ri.FWVersion != "" {
		parts = append(parts, ri.FWVersion)
	}
	if ri.FWBuildID != "" {
		parts = append(parts, ri.FWBuildID)
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, " "))
}

func NewOTAFileHandler(root string, chunkSize int, delay time.Duration) *OTAFileHandler {
	return &OTAFileHandler{root: root, chunkSize: chunkSize, delay: delay}
}
