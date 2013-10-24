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

// logstashAdapter is a glogJSON that decodes glog data and encodes it into JSON.
var logstashAdapter *glogJSON

// Set the io.Writer to write JSON. This is required if -logstash=true
func SetLogstashWriter(writer io.Writer) {
	logstashAdapter = &glogJSON{writer: writer, encoder: json.NewEncoder(writer)}
}

func init() {
	flag.BoolVar(&logging.toLogstash, "logstash", false, "log also in JSON using the Logstash writer")
	// Write to Stderr until SetLogstashWriter is called so we do not loose events.
	SetLogstashWriter(os.Stderr)
}

// glogJSON can decode the data generated by glog and encode it in logstash json format.
// https://gist.github.com/jordansissel/2996677
// implements io.Writer
type glogJSON struct {
	writer  io.Writer
	encoder *json.Encoder
}

// Write decodes the data and writes a logstash json event
func (d glogJSON) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}
	d.openEvent()
	// peek for normal logline
	sev := data[0]
	switch sev {
	case 73, 87, 69, 70: // IWEF
		d.iwef(sev, data)
	default:
		d.message(string(data))
	}
	d.closeHash()
	// My Write always succeeds
	return len(data), nil
}

// openEvent writes the "header" part of the JSON message.
func (d glogJSON) openEvent() {
	io.WriteString(d.writer, `{"@source":`)
	d.encoder.Encode(host) // uses glog package var
	io.WriteString(d.writer, `,"@type":"glog"`)
	io.WriteString(d.writer, `,"@timestamp":`)
	// ignore time information given, take new snapshot
	d.encoder.Encode(timeNow()) // use testable function stored in var
}

// closeAll writes the closing brackets for the main and fields hash.
func (d glogJSON) closeHash() {
	io.WriteString(d.writer, "}\n")
}

// message add a JSON field with the JSON encoded message.
func (d glogJSON) message(msg string) {
	io.WriteString(d.writer, `,"@message":`)
	d.encoder.Encode(msg)
}

// iwef decodes a glog data packet and write the JSON representation.
// [IWEF]mmdd hh:mm:ss.uuuuuu threadid file:line] msg
func (d glogJSON) iwef(sev byte, data []byte) {
	io.WriteString(d.writer, `,"@fields":{"level":"`)
	switch sev {
	case 73:
		io.WriteString(d.writer, "INFO")
	case 87:
		io.WriteString(d.writer, "ERROR")
	case 69:
		io.WriteString(d.writer, "WARNING")
	case 70:
		io.WriteString(d.writer, "FATAL")
	}
	r := &iwefreader{data, 22} // past last u
	io.WriteString(d.writer, `","threadid":`)
	io.WriteString(d.writer, r.stringUpTo(32)) // space
	r.skip()                                   // space
	io.WriteString(d.writer, `,"file":"`)
	io.WriteString(d.writer, r.stringUpTo(58)) // :
	r.skip()                                   // :
	io.WriteString(d.writer, `","line":`)
	io.WriteString(d.writer, r.stringUpTo(93)) // ]
	// ]
	r.skip()
	// space
	r.skip()
	// fields
	d.closeHash()
	d.message(r.stringUpToLineEnd())
}

// iwefreader is a small helper object to parse a glog IWEF entry
type iwefreader struct {
	data     []byte
	position int // read offset in data
}

// skip advances the position in data
func (i *iwefreader) skip() {
	i.position++
}

func (i iwefreader) stringUpToLineEnd() string {
	return string(i.data[i.position : len(i.data)-1]) // without the line delimiter
}

func (i *iwefreader) stringUpTo(delim byte) string {
	start := i.position
	for i.data[i.position] != delim {
		i.position++
	}
	return string(i.data[start:i.position])
}

func (i *iwefreader) intUpTo(delim byte) int {
	val := int(0)
	for i.data[i.position] != delim {
		val = (val * 10) + int(i.data[i.position]) - 48 // 0
		i.position++
	}
	return val
}
