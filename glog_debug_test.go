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
	"testing"
)

// go test -v -test.run TestDebug ...glog
func TestDebug(t *testing.T) {
	logging.toStderr = true
	logging.verbosity = DEBUG
	logit()
}

func TestTrace(t *testing.T) {
	logging.toStderr = true
	logging.verbosity = TRACE
	logit()
}

func TestInfoOnly(t *testing.T) {
	logging.toStderr = true
	logging.verbosity = 0
	logit()
}

func logit() {
	Info("this is info level")
	Debug("this is debug level")
	if DebugEnabled() {
		Info("info is debug")
	}
	Trace("this is trace level")
	if TraceEnabled() {
		Info("info is trace")
	}
}
