package gzip

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"mime"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Gzip interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type GzipConfig struct {
	MinSize     int      `config:"minsize"`
	Level       int      `config:"level"`
	ContentType []string `config:"contenttype"`
}

type gzipWriter struct {
	writer *gzip.Writer
	buffer []byte
}

type gzipImplement struct {
	level       int
	minSize     int
	contentType map[string]bool
	pool        *sync.Pool
}

func NewGzip(config GzipConfig) (Gzip, error) {
	if config.Level == 0 {
		config.Level = gzip.DefaultCompression
	}
	if config.MinSize <= 0 {
		config.MinSize = 1024
	}
	contentType := map[string]bool{}
	for _, singleContentType := range config.ContentType {
		mediaType, _, err := mime.ParseMediaType(singleContentType)
		if err != nil {
			return nil, err
		}
		contentType[mediaType] = true
	}
	pool := &sync.Pool{
		New: func() interface{} {
			w, err := gzip.NewWriterLevel(nil, config.Level)
			if err != nil {
				panic(err)
			}
			buffer := make([]byte, config.MinSize, config.MinSize)
			return &gzipWriter{
				writer: w,
				buffer: buffer,
			}
		},
	}
	return &gzipImplement{
		level:       config.Level,
		minSize:     config.MinSize,
		contentType: contentType,
		pool:        pool,
	}, nil
}

func (this *gzipImplement) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add(vary, acceptEncoding)
	gw := &gzipResponseWriter{
		ResponseWriter: w,
		r:              r,
		gzipImpl:       this,
		hasParse:       false,
		enableGzip:     false,
		writer:         nil,
		totalSize:      0,
		statusCode:     0,
	}
	defer gw.Close()

	if _, ok := w.(http.CloseNotifier); ok {
		gwcn := gzipResponseWriterWithCloseNotify{gw}
		next.ServeHTTP(gwcn, r)
	} else {
		next.ServeHTTP(gw, r)
	}
}

const (
	vary            = "Vary"
	acceptEncoding  = "Accept-Encoding"
	contentEncoding = "Content-Encoding"
	contentType     = "Content-Type"
	contentLength   = "Content-Length"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	r          *http.Request
	gzipImpl   *gzipImplement
	hasParse   bool
	enableGzip bool
	writer     *gzipWriter
	totalSize  int
	statusCode int
}

type gzipResponseWriterWithCloseNotify struct {
	*gzipResponseWriter
}

func (w gzipResponseWriterWithCloseNotify) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w *gzipResponseWriter) parseEncodings(s string) (map[string]float64, error) {
	c := map[string]float64{}
	var e []string

	for _, ss := range strings.Split(s, ",") {
		coding, qvalue, err := w.parseCoding(ss)

		if err != nil {
			e = append(e, err.Error())
		} else {
			c[coding] = qvalue
		}
	}

	if len(e) > 0 {
		return c, fmt.Errorf("errors while parsing encodings: %s", strings.Join(e, ", "))
	}

	return c, nil
}

func (w *gzipResponseWriter) parseCoding(s string) (coding string, qvalue float64, err error) {
	for n, part := range strings.Split(s, ";") {
		part = strings.TrimSpace(part)
		qvalue = 1.0

		if n == 0 {
			coding = strings.ToLower(part)
		} else if strings.HasPrefix(part, "q=") {
			qvalue, err = strconv.ParseFloat(strings.TrimPrefix(part, "q="), 64)

			if qvalue < 0.0 {
				qvalue = 0.0
			} else if qvalue > 1.0 {
				qvalue = 1.0
			}
		}
	}

	if coding == "" {
		err = fmt.Errorf("empty content-coding")
	}

	return
}

func (w *gzipResponseWriter) parseAcceptEncoding() bool {
	acceptedEncodings, err := w.parseEncodings(w.r.Header.Get(acceptEncoding))
	if err != nil {
		return false
	}
	return acceptedEncodings["gzip"] > 0.0
}

func (w *gzipResponseWriter) parseContentType() bool {
	if len(w.gzipImpl.contentType) == 0 {
		return true
	}
	contentType := w.Header().Get(contentType)
	if contentType == "" {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return w.gzipImpl.contentType[mediaType]
}

func (w *gzipResponseWriter) parseContentEncoding() bool {
	return w.Header().Get(contentEncoding) == ""
}

func (w *gzipResponseWriter) shouldGzip() bool {
	if w.hasParse == true {
		return w.enableGzip
	}
	parseHandler := [...]func() bool{
		w.parseContentEncoding,
		w.parseAcceptEncoding,
		w.parseContentType,
	}
	for _, handler := range parseHandler {
		isOk := handler()
		if isOk == false {
			w.hasParse = true
			w.enableGzip = false
			return false
		}
	}
	w.hasParse = true
	w.enableGzip = true
	return true
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	shouldGzip := w.shouldGzip()
	if shouldGzip == false {
		return w.ResponseWriter.Write(b)
	} else {
		err := w.handleGzip(b)
		if err != nil {
			return 0, err
		}
		return len(b), nil
	}
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	shouldGzip := w.shouldGzip()
	if shouldGzip == false {
		w.ResponseWriter.WriteHeader(statusCode)
	} else {
		w.statusCode = statusCode
	}
}

func (w *gzipResponseWriter) handleGzip(b []byte) error {

	if len(b) == 0 {
		return nil
	}

	if w.writer == nil {
		w.writer = w.gzipImpl.pool.Get().(*gzipWriter)
	}

	if w.totalSize+len(b) <= w.gzipImpl.minSize {
		//还未超过minSize
		copy(w.writer.buffer[w.totalSize:], b)
	} else {
		//超过minSize
		if w.totalSize <= w.gzipImpl.minSize {
			w.Header().Set(contentEncoding, "gzip")
			w.Header().Del(contentLength)
			if w.statusCode != 0 {
				w.ResponseWriter.WriteHeader(w.statusCode)
				w.statusCode = 0
			}
			w.writer.writer.Reset(w.ResponseWriter)
			if w.totalSize != 0 {
				_, err := w.writer.writer.Write(w.writer.buffer[0:w.totalSize])
				if err != nil {
					return err
				}
			}
		}
		_, err := w.writer.writer.Write(b)
		if err != nil {
			return err
		}
	}

	w.totalSize += len(b)
	return nil
}

func (w *gzipResponseWriter) Close() error {
	if w.statusCode != 0 {
		w.ResponseWriter.WriteHeader(w.statusCode)
	}

	if w.writer == nil {
		return nil
	}

	if w.totalSize <= w.gzipImpl.minSize {
		//还未超过minSize
		if w.totalSize != 0 {
			_, err := w.ResponseWriter.Write(w.writer.buffer[0:w.totalSize])
			if err != nil {
				return err
			}
		}
	} else {
		//超过minSize
		err := w.writer.writer.Close()
		w.gzipImpl.pool.Put(w.writer)
		if err != nil {
			return err
		}
	}
	w.writer = nil
	return nil
}

func (w *gzipResponseWriter) Flush() {
	//FIXME 如果数据滞留在buffer里面，则没能及时地flush
	if w.writer == nil {
		return
	}

	w.writer.writer.Flush()

	if fw, ok := w.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}

func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("http.Hijacker interface is not supported")
}

var _ http.Hijacker = &gzipResponseWriter{}

var _ http.Flusher = &gzipResponseWriter{}

var _ http.CloseNotifier = &gzipResponseWriterWithCloseNotify{}
