/*
 * Copyright © 2022 photowey (photowey@gmail.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logger

type Level = int

const (
	TraceLevel Level = 1 << iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

var _ Logger = (*loggerx)(nil)

// Logger - the logger root API
type Logger interface {
	Trace(template string, args ...any)
	Debug(template string, args ...any)
	Info(template string, args ...any) // default
	Warn(template string, args ...any)
	Error(template string, args ...any)
	IsTraceEnabled() bool
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarnEnabled() bool
	IsErrorEnabled() bool
}

// loggerx - default implementation of Logger
type loggerx struct {
	level Level
}

func NewLogger() Logger {
	return &loggerx{}
}

func (lx loggerx) Trace(template string, args ...any) {

}
func (lx loggerx) Debug(template string, args ...any) {

}
func (lx loggerx) Info(template string, args ...any) {

}
func (lx loggerx) Warn(template string, args ...any) {

}
func (lx loggerx) Error(template string, args ...any) {

}

func (lx loggerx) IsTraceEnabled() bool {
	return false
}

func (lx loggerx) IsDebugEnabled() bool {
	return false
}

func (lx loggerx) IsInfoEnabled() bool {
	return false
}

func (lx loggerx) IsWarnEnabled() bool {
	return false
}

func (lx loggerx) IsErrorEnabled() bool {
	return false
}
