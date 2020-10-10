package wrappers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	ApacheFormatPattern = "%s - - [%s] \"%s %d %d\" %f\n"
)

type ApacheLogRecord struct {
	ip                    string
	time                  time.Time
	method, uri, protocol string
	status                int
	responseBytes         int64
	elapsedTime           time.Duration
}

func (r *ApacheLogRecord) Log(out io.Writer) {
	timeFormatted := r.time.Format("02/Jan/2006 03:04:05")
	requestLine := fmt.Sprintf("%s %s %s", r.method, r.uri, r.protocol)
	fmt.Fprintf(out, ApacheFormatPattern, r.ip, timeFormatted, requestLine, r.status, r.responseBytes,
		r.elapsedTime.Seconds())
}

func WithLog(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
			clientIP = clientIP[:colon]
		}

		record := &ApacheLogRecord{
			ip:          clientIP,
			time:        time.Time{},
			method:      r.Method,
			uri:         r.RequestURI,
			protocol:    r.Proto,
			status:      http.StatusOK,
			elapsedTime: time.Duration(0),
		}

		startTime := time.Now()

		f(w, r)
		finishTime := time.Now()

		record.time = finishTime.UTC()
		record.elapsedTime = finishTime.Sub(startTime)

		record.Log(os.Stderr)
	}
}
