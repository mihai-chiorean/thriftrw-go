// This package provides a plugin that generates code used by the plugin
// system itself.

// Note that this plugin doesn't support imports.

package main

import (
	"path/filepath"
	"strings"

	"github.com/thriftrw/thriftrw-go/plugin"
	"github.com/thriftrw/thriftrw-go/plugin/api"
)

func main() {
	plugin.Main(&plugin.Plugin{
		Name:      "envelope",
		Generator: generator{},
	})
}

type generator struct{}

func (generator) Generate(req *api.GenerateRequest) (*api.GenerateResponse, error) {
	files := make(map[string][]byte)
	for _, serviceID := range req.RootServices {
		service := req.Services[serviceID]
		module := req.Modules[service.ModuleID]

		templateData := struct {
			Service *api.Service
			Request *api.GenerateRequest
		}{Service: service, Request: req}

		var (
			err       error
			ifaceOpts []plugin.TemplateOption
		)

		ifacePath := filepath.Join(module.Directory, strings.ToLower(service.Name)+".go")
		ifaceOpts = append(ifaceOpts, templateOptions...)
		ifaceOpts = append(ifaceOpts, plugin.GoFileImportPath(module.Package))
		files[ifacePath], err = plugin.GoFileFromTemplate(
			ifacePath, interfaceTemplate, templateData, ifaceOpts...)
		if err != nil {
			return nil, err
		}

		envPath := filepath.Join(service.Directory, "envelope.go")
		files[envPath], err = plugin.GoFileFromTemplate(
			envPath, serviceTemplate, templateData, templateOptions...)
		if err != nil {
			return nil, err
		}
	}
	return &api.GenerateResponse{Files: files}, nil
}

// convenience function because "index .Request.Services .Service.ParentID"
// doesn't work. Index expects a ServiceID, not *ServiceID.
func getService(req *api.GenerateRequest, id api.ServiceID) *api.Service {
	return req.Services[id]
}

var templateOptions = []plugin.TemplateOption{
	plugin.TemplateFunc("basename", filepath.Base),
	plugin.TemplateFunc("getService", getService),
}

const interfaceTemplate = `
<$module := index .Request.Modules .Service.ModuleID>
package <basename $module.Package>

type <.Service.Name> interface {
	<if .Service.ParentID>
		<$parent := getService .Request .Service.ParentID>
		<if eq $parent.ModuleID .Service.ModuleID>
			<$parent.Name>
		<else>
			<$parentModule := index .Request.Modules $parent.ModuleID>
			<import $parentModule.Package>.<$parent.Name>
		<end>
	<end>

	<range .Service.Functions>
		<.Name>(<range .Arguments>
			<.Name> <formatType .Type>,<end>
		) <if .ReturnType>(<formatType .ReturnType>, error)<else>error<end>
	<end>
}
`

const serviceTemplate = `
package <basename .Service.Package>

<$bytes    := import "bytes">
<$protocol := import "github.com/thriftrw/thriftrw-go/protocol">
<$envelope := import "github.com/thriftrw/thriftrw-go/envelope">
<$wire     := import "github.com/thriftrw/thriftrw-go/wire">
<$module   := import (index .Request.Modules .Service.ModuleID).Package>

// Protocol is the Thrift protocol used to deserialize and serialize requests.
var Protocol = <$protocol>.Binary

// Client implements a <.Service.Name> client.
type client struct {
	<if .Service.ParentID>
		<$parent := getService .Request .Service.ParentID>
		<$parentModule := index .Request.Modules $parent.ModuleID>
		<import $parentModule.Package>.<$parent.Name>
	<end>
	send func([]byte) ([]byte, error)
}

// NewClient builds a new <.Service.Name> client.
func NewClient(t func([]byte) ([]byte, error)) <$module>.<.Service.Name> {
	return &client{
		send: t,
		<if .Service.ParentID>
			<$parent := getService .Request .Service.ParentID>
			<$parent.Name>: <import $parent.Package>.NewClient(t),
		<end>
	}
}

<range .Service.Functions>
func (c *client) <.Name>(<range .Arguments>
	_<.Name> <formatType .Type>,<end>
) (<if .ReturnType>success <formatType .ReturnType>,<end> err error) {
	args := <.Name>Helper.Args(<range .Arguments>_<.Name>, <end>)

	var buff <$bytes>.Buffer
	if err = <$envelope>.Write(Protocol, &buff, 1, args); err != nil {
		return
	}

	var resBody []byte
	resBody, err = c.send(buff.Bytes())
	if err != nil {
		return
	}

	var body <$wire>.Value
	body, _, err = <$envelope>.ReadReply(Protocol, <$bytes>.NewReader(resBody))
	if err != nil {
		return
	}

	var result <.Name>Result
	if err = result.FromWire(body); err != nil {
		return
	}

	<if .ReturnType>success, <end>err = <.Name>Helper.UnwrapResponse(&result)
	return
}
<end>

// Handler serves an implementation of the <.Service.Name> service.
type Handler struct {
	impl <$module>.<.Service.Name>

	<if .Service.ParentID>
		<$parent := getService .Request .Service.ParentID>
		parent <import $parent.Package>.Handler
	<end>
}

// NewHandler builds a new <.Service.Name> handler.
func NewHandler(service <$module>.<.Service.Name>) Handler {
	return Handler{
		impl: service,
		<if .Service.ParentID>
			<$parent := getService .Request .Service.ParentID>
			parent: <import $parent.Package>.NewHandler(service),
		<end>
	}
}

// Handle handles the given request and returns a response.
func (h Handler) Handle(data []byte) ([]byte, error) {
	envelope, err := Protocol.DecodeEnveloped(<$bytes>.NewReader(data))
	if err != nil {
		return nil, err
	}

	responseEnvelope, err := h.HandleEnvelope(envelope)
	if err != nil {
		return nil, err
	}

	var buff <$bytes>.Buffer
	if err := Protocol.EncodeEnveloped(responseEnvelope, &buff); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// HandleEnvelope receives an enveloped request for <.Service.Name> service
// and returns an enveloped response.
func (h Handler) HandleEnvelope(envelope <$wire>.Envelope) (<$wire>.Envelope, error) {
	responseEnvelope := <$wire>.Envelope{
		Name: envelope.Name,
		Type: <$wire>.Reply,
		SeqID: envelope.SeqID,
	}

	switch envelope.Name {
		<range .Service.Functions>
			case "<.ThriftName>":
				var args <.Name>Args
				if err := args.FromWire(envelope.Value); err != nil {
					return responseEnvelope, err
				}

				result, err := <.Name>Helper.WrapResponse(
					h.impl.<.Name>(<range .Arguments>args.<.Name>, <end>),
				)
				if err != nil {
					return responseEnvelope, err
				}

				responseEnvelope.Value, err = result.ToWire()
				if err != nil {
					return responseEnvelope, err
				}
		<end>
		default:
			<if .Service.ParentID>
				return h.parent.HandleEnvelope(envelope)
			<else>
				// TODO(abg): Use TException
				return responseEnvelope, <import "fmt">.Errorf("unknown method %q", envelope.Name)
			<end>
	}

	return responseEnvelope, nil
}
`
