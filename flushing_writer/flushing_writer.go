package flushing_writer

import (
	"net/http"
	"time"
)

type FlushingResponseWriter interface {
	http.ResponseWriter
	http.Flusher
}

type flushingResponseWriter struct {
	w         FlushingResponseWriter
	chunkSize int
	delay     time.Duration

	numPending int
}

func NewFlushingResponseWriter(w FlushingResponseWriter, chunkSize int, delay time.Duration) FlushingResponseWriter {
	return &flushingResponseWriter{w: w, chunkSize: chunkSize, delay: delay}
}

func (fw *flushingResponseWriter) Header() http.Header {
	return fw.w.Header()
}

func (fw *flushingResponseWriter) Write(data []byte) (int, error) {
	totalNumWritten := 0
	for len(data) > 0 {
		toWrite := fw.chunkSize - fw.numPending
		if len(data) < toWrite {
			toWrite = len(data)
		}
		numWritten, err := fw.w.Write(data[:toWrite])
		if numWritten > 0 {
			totalNumWritten += numWritten
			fw.numPending += numWritten
			if fw.numPending >= fw.chunkSize {
				fw.Flush()
			}
		}
		if err != nil || numWritten == 0 {
			return totalNumWritten, err
		}
		data = data[numWritten:]
	}
	return totalNumWritten, nil
}

func (fw *flushingResponseWriter) WriteHeader(statusCode int) {
	fw.w.WriteHeader(statusCode)
}

func (fw *flushingResponseWriter) Flush() {
	fw.w.Flush()
	if fw.delay > 0 && fw.numPending > 0 {
		time.Sleep(fw.delay)
	}
	fw.numPending = 0
}
