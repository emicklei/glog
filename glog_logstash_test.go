// Go support for leveled logs, analogous to https://code.google.com/p/google-glog/
//
// Modifications copyright 2013 Ernest Micklei. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package glog

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

// go test -v -test.run TestInfoLogstash ...glog
func TestInfoLogstash(t *testing.T) {
	defer func(previous func() time.Time) { timeNow = previous }(timeNow)
	timeNow = func() time.Time {
		return time.Date(2006, 1, 2, 15, 4, 5, .678901e9, time.Local)
	}
	logstash.toLogstash = true // simulate -logstash=true
	host = "unknownhost"
	capture := new(bytes.Buffer)
	SetLogstashWriter(capture)
	Info("hello")
	actual := capture.String()
	if strings.HasPrefix(actual, jsonBegin) && strings.HasSuffix(actual, jsonEnd) {
		println(actual)
		t.Fatalf("mismatch in json")
	}
	Flush()
}

var jsonBegin = `{"@source_host":"unknownhost"
,"@timestamp":"2006-01-02T15:04:05.678901+01:00"
,"@fields":{"level":"INFO","threadid":`

var jsonEnd = `"file":"glog_logstash_test.go","line":18}
,"@message":"hello"
}
`

// go test -v -test.run TestEnabledLogstashNoWriter ...glog
func TestEnabledLogstashNoWriter(t *testing.T) {
	logstash.toLogstash = true
	SetLogstashWriter(os.Stdout)
	ExtraFields["instance"] = "ps34"
	ExtraFields["role"] = "webservice"
	Info("hello")
	Info("world")
	Flush()
	logstash.toLogstash = false
}

type failingWriter struct{}

func (f failingWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("simulated fail")
}

// go test -v -test.run TestErrorWritingLogstashNoWriter ...glog
func TestErrorWritingLogstashNoWriter(t *testing.T) {
	logstash.toLogstash = true
	SetLogstashWriter(failingWriter{})
	ExtraFields["instance"] = "ps34"
	ExtraFields["role"] = "webservice"
	Info("hello")
	Info("world")
	Flush()
	logstash.toLogstash = false
}

func ExampleSetLogstashWriter() {
	// TODO write the UDP example
	logstash, err := os.Create("logstash.log")
	if err == nil {
		defer logstash.Close()
		SetLogstashWriter(logstash)
	} else {
		os.Stderr.Write([]byte("unable to create logstash.log:" + err.Error()))
	}
}
