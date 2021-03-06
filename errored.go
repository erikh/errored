/***
Copyright 2014 Cisco Systems Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package errored

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

type errorStack struct {
	file string
	line int
	fun  string
}

// Error is our custom error with description, file, and line.
type Error struct {
	desc  string
	stack []errorStack
}

// Error() allows *core.Error to present the `error` interface.
func (e *Error) Error() string {
	var ret string

	if os.Getenv("CONTIV_TRACE") != "" {
		ret = e.desc + "\n"

		for _, stack := range e.stack {
			ret += fmt.Sprintf("%s [%s %d]\n", stack.fun, stack.file, stack.line)
		}
	} else {
		ret = fmt.Sprintf("%s [%s %s %d]", e.desc, e.stack[0].fun, e.stack[0].file, e.stack[0].line)
	}

	return ret
}

// Errorf returns an *Error based on the format specification provided.
func Errorf(f string, args ...interface{}) *Error {
	e := &Error{
		stack: []errorStack{},
		desc:  fmt.Sprintf(f, args...),
	}

	i := 1

	for {
		stack := errorStack{}
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fun := runtime.FuncForPC(pc)
		if fun != nil {
			stack.fun = fun.Name()
		}

		stack.file = path.Base(file)
		stack.line = line
		e.stack = append(e.stack, stack)

		i++
	}

	return e
}

// ErrIfKeyExists checks if the error message contains "Key not found".
func ErrIfKeyExists(err error) error {
	if err == nil || strings.Contains(err.Error(), "Key not found") {
		return nil
	}

	return err
}
