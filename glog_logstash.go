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
	"encoding/json"
	"flag"
	"io"
	"os"
)

/*
{
   "@source":"test.here.com",
   "@type":"glog",
   "@timestamp":"2013-10-24T09:30:46.947024155+02:00",
   "@fields":{
      "level":"INFO",
      "threadid":400004,
      "file":"file.go",
      "line":10
   },
   "@message":"hello"
}
*/

// ExtraFields contains a set of @fields elements that can be used by the application
// to pass appliction and/or environment specific information.
var ExtraFields = map[string]string{}

// logstash is a logstashPublisher that decodes each glog data,
// encodes it into JSON and writes it asynchronuously to an io.Writer.
var logstash logstashPublisher

// Set the io.Writer to write JSON. This is required if -logstash=true
func SetLogstashWriter(writer io.Writer) {
	logstash.writer = newBufferedWriter(writer)
}

func init() {
	flag.BoolVar(&logstash.toLogstash, "logstash", false, "log also in JSON using the Logstash writer")
	// Write to Stderr until SetLogstashWriter is called so we do not loose events.
	SetLogstashWriter(os.Stderr)
}

// logstashPublisher holds global state for publishing messages in JSON.
type logstashPublisher struct {
	toLogstash bool            // The -logstash flag.
	writer     *bufferedWriter // Buffered target writer for JSON messages.
}

// WriteWithStack decodes the data and writes a logstash json event
func (p logstashPublisher) WriteWithStack(data []byte, stack []byte) {
	buffer := new(bytes.Buffer)
	glogJSON{writer: buffer, encoder: json.NewEncoder(buffer)}.WriteWithStack(data, stack)
	p.writer.Write(buffer.Bytes())
}

// flush waits until all pending messages are written by the asyncWriter.
func (p logstashPublisher) flush() {
	if p.writer != nil { // be robust
		p.writer.flush()
	}
}

// bufferedWriter collects []byte until a flush.
type bufferedWriter struct {
	buffer [][]byte
	writer io.Writer
}

// newBufferedWriter decorates the underlyingWriter.
func newBufferedWriter(underlyingWriter io.Writer) *bufferedWriter {
	bw := new(bufferedWriter)
	bw.buffer = [][]byte{}
	bw.writer = underlyingWriter
	return bw
}

// Write is for implementing io.Writer
func (b *bufferedWriter) Write(data []byte) (n int, err error) {
	b.buffer = append(b.buffer, data)
	return len(data), nil
}

// flush drains the buffer. it is called from the daemon goroutine.
func (b *bufferedWriter) flush() {
	for _, each := range b.buffer {
		_, err := b.writer.Write(each)
		if err != nil {
			os.Stderr.WriteString("[glog error] unable to flush buffered logstash message:\n")
			os.Stderr.WriteString(string(each))
		}
	}
	b.buffer = [][]byte{}
}
