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

const (
	DEBUG = 10 // severity levels
	TRACE = 100
)

// DebugEnabled returns true if the severity level is set to DEBUG or higher.
func DebugEnabled() bool {
	return bool(V(DEBUG))
}

// TraceEnabled returns true if the severity level is set to TRACE or higher.
func TraceEnabled() bool {
	return bool(V(TRACE))
}

// Debug prints the message if the severity level is set to DEBUG or higher.
func Debug(message string) {
	V(DEBUG).Info(message)
}

// Debug prints the formatted message if the severity level is set to DEBUG or higher.
func Debugf(format string, args ...interface{}) {
	V(DEBUG).Infof(format, args...)
}

// Trace prints the message if the severity level is set to TRACE or higher.
func Trace(message string) {
	V(TRACE).Info(message)
}

// Tracef prints the formatted message if the severity level is set to TRACE or higher.
func Tracef(format string, args ...interface{}) {
	V(TRACE).Infof(format, args...)
}