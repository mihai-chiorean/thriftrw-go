// Code generated by thriftrw
// @generated

// Copyright (c) 2016 Uber Technologies, Inc.
//
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

package plugin

import (
	"fmt"
	"github.com/thriftrw/thriftrw-go/wire"
	"strings"
)

type GoodbyeArgs struct{}

func (v *GoodbyeArgs) ToWire() (wire.Value, error) {
	var (
		fields [0]wire.Field
		i      int = 0
	)
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *GoodbyeArgs) FromWire(w wire.Value) error {
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		}
	}
	return nil
}

func (v *GoodbyeArgs) String() string {
	var fields [0]string
	i := 0
	return fmt.Sprintf("GoodbyeArgs{%v}", strings.Join(fields[:i], ", "))
}

func (v *GoodbyeArgs) MethodName() string {
	return "goodbye"
}

func (v *GoodbyeArgs) EnvelopeType() wire.EnvelopeType {
	return wire.Call
}

type GoodbyeResult struct{}

func (v *GoodbyeResult) ToWire() (wire.Value, error) {
	var (
		fields [0]wire.Field
		i      int = 0
	)
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *GoodbyeResult) FromWire(w wire.Value) error {
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		}
	}
	return nil
}

func (v *GoodbyeResult) String() string {
	var fields [0]string
	i := 0
	return fmt.Sprintf("GoodbyeResult{%v}", strings.Join(fields[:i], ", "))
}

func (v *GoodbyeResult) MethodName() string {
	return "goodbye"
}

func (v *GoodbyeResult) EnvelopeType() wire.EnvelopeType {
	return wire.Reply
}

var GoodbyeHelper = struct {
	IsException    func(error) bool
	Args           func() *GoodbyeArgs
	WrapResponse   func(error) (*GoodbyeResult, error)
	UnwrapResponse func(*GoodbyeResult) error
}{}

func init() {
	GoodbyeHelper.IsException = func(err error) bool {
		switch err.(type) {
		default:
			return false
		}
	}
	GoodbyeHelper.Args = func() *GoodbyeArgs {
		return &GoodbyeArgs{}
	}
	GoodbyeHelper.WrapResponse = func(err error) (*GoodbyeResult, error) {
		if err == nil {
			return &GoodbyeResult{}, nil
		}
		return nil, err
	}
	GoodbyeHelper.UnwrapResponse = func(result *GoodbyeResult) (err error) {
		return
	}
}
