// Copyright (c) 2015 Uber Technologies, Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package testutils

import (
	"fmt"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/uber/tchannel-go"
)

type errorLoggerState struct {
	matchCount []uint32
}

type errorLogger struct {
	tchannel.Logger
	t testing.TB
	v *LogVerification
	s *errorLoggerState
}

func (l errorLogger) checkErr(msg string, args ...interface{}) {
	if len(l.v.Filters) == 0 {
		l.t.Errorf(msg, args...)
		return
	}

	formatted := fmt.Sprintf(msg, args...)
	match := -1
	for i, filter := range l.v.Filters {
		if strings.Contains(formatted, filter.Filter) {
			match = i
		}
	}

	if match >= 0 {
		matchCount := atomic.AddUint32(&l.s.matchCount[match], 1)
		if uint(matchCount) <= l.v.Filters[match].Count {
			return
		}
	}

	l.t.Errorf(msg, args...)
}

func (l errorLogger) Fatalf(msg string, args ...interface{}) {
	l.checkErr("[Fatal] "+msg, args...)
	l.Logger.Fatalf(msg, args...)
}

func (l errorLogger) Errorf(msg string, args ...interface{}) {
	l.checkErr("[Error] "+msg, args...)
	l.Logger.Errorf(msg, args...)
}

func (l errorLogger) Warnf(msg string, args ...interface{}) {
	l.checkErr("[Warn] "+msg, args...)
	l.Logger.Warnf(msg, args...)
}

func (l errorLogger) WithFields(fields ...tchannel.LogField) tchannel.Logger {
	return errorLogger{l.Logger.WithFields(fields...), l.t, l.v, l.s}
}
