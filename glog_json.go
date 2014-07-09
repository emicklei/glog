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
	"io"
)

// glogJSON can decode a data slice generated by glog and encode it in logstash json format.
// https://gist.github.com/jordansissel/2996677
type glogJSON struct {
	writer  *bytes.Buffer // for the composition of one message.
	encoder *json.Encoder // used to encode string parameters.
}

// WriteWithStack decodes the data and writes a logstash json event
func (d glogJSON) WriteWithStack(data []byte, stack []byte) {
	d.openEvent()
	// peek for normal logline
	sev := data[0]
	switch sev {
	case 73, 87, 69, 70: // IWEF
		d.iwef(sev, data, stack)
	default:
		d.message(string(data))
	}
	d.closeHash()
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

// message adds a JSON field with the JSON encoded message.
func (d glogJSON) message(msg string) {
	io.WriteString(d.writer, `,"@message":`)
	d.encoder.Encode(msg)
}

// stack adds a JSON field with the JSON encoded stack trace of all goroutines.
func (d glogJSON) stacktrace(stacktrace []byte) {
	io.WriteString(d.writer, `,"stack":`)
	d.encoder.Encode(string(stacktrace))
}

// iwef decodes a glog data packet and write the JSON representation.
// [IWEF]mmdd hh:mm:ss.uuuuuu threadid file:line] msg
func (d glogJSON) iwef(sev byte, data []byte, trace []byte) {
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
	if trace != nil && len(trace) > 0 {
		d.stacktrace(trace)
	}
	// extras?
	for k, v := range ExtraFields {
		io.WriteString(d.writer, `,"`)
		io.WriteString(d.writer, k)
		io.WriteString(d.writer, `":`)
		d.encoder.Encode(v)
	}
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

// stringUpToLineEnd returns the string part from the data up to not-including the line end.
func (i iwefreader) stringUpToLineEnd() string {
	return string(i.data[i.position : len(i.data)-1]) // without the line delimiter
}

// stringUpTo returns the string part from the data up to not-including a delimiter.
func (i *iwefreader) stringUpTo(delim byte) string {
	start := i.position
	for i.data[i.position] != delim {
		i.position++
	}
	return string(i.data[start:i.position])
}
