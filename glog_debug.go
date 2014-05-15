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

import "strconv"

const (
	DEBUG = 10 // severity levels
	TRACE = 100
)

// DebugEnabled returns true if the severity level is set to DEBUG or higher.
func DebugEnabled() bool {
	return logging.verbosity >= DEBUG
}

// TraceEnabled returns true if the severity level is set to TRACE or higher.
func TraceEnabled() bool {
	return logging.verbosity >= TRACE
}

// Debug prints the message if the severity level is set to DEBUG or higher.
func Debug(args ...interface{}) {
	if logging.verbosity >= DEBUG {
		logging.print(infoLog, args)
	}
}

// Debug prints the formatted message if the severity level is set to DEBUG or higher.
func Debugf(format string, args ...interface{}) {
	if logging.verbosity >= DEBUG {
		logging.printf(infoLog, format, args...)
	}
}

// Trace prints the message if the severity level is set to TRACE or higher.
func Trace(args ...interface{}) {
	if logging.verbosity >= TRACE {
		logging.print(infoLog, args)
	}
}

// Tracef prints the formatted message if the severity level is set to TRACE or higher.
func Tracef(format string, args ...interface{}) {
	if logging.verbosity >= TRACE {
		logging.printf(infoLog, format, args...)
	}
}

// SetVerbosity changes the current verbosity level to v.
func SetVerbosity(v int) {
	logging.verbosity.Set(strconv.Itoa(v))
}

// SetLogToStdErr enables write to stderr.
func SetLoggingToStdErr() {
	logging.toStderr = true
}
