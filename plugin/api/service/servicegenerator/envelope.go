package servicegenerator

import (
	"fmt"
	"github.com/thriftrw/thriftrw-go/plugin/api"
	"github.com/thriftrw/thriftrw-go/wire"
)

// Client implements a ServiceGenerator client.
type client struct {
	send func(wire.Envelope) (wire.Envelope, error)
}

// NewClient builds a new ServiceGenerator client.
func NewClient(t func(wire.Envelope) (wire.Envelope, error)) api.ServiceGenerator {
	return &client{
		send: t,
	}
}

func (c *client) Generate(
	_Request *api.GenerateServiceRequest,
) (success *api.GenerateServiceResponse, err error) {
	args := GenerateHelper.Args(_Request)

	var body wire.Value
	body, err = args.ToWire()
	if err != nil {
		return
	}

	var envelope wire.Envelope
	envelope, err = c.send(wire.Envelope{
		Name:  "generate",
		Type:  wire.Call,
		Value: body,
	})
	if err != nil {
		return
	}

	switch {
	case envelope.Type == wire.Exception:
		// TODO(abg): use envelope exceptions
		err = fmt.Errorf("envelope error: %v", envelope.Value)
		return
	case envelope.Type != wire.Reply:
		err = fmt.Errorf("unknown envelope type for reply, got %v", envelope.Type)
		return
	}

	var result GenerateResult
	if err = result.FromWire(envelope.Value); err != nil {
		return
	}

	success, err = GenerateHelper.UnwrapResponse(&result)
	return
}

// Handler serves an implementation of the ServiceGenerator service.
type Handler struct {
	impl api.ServiceGenerator
}

// NewHandler builds a new ServiceGenerator handler.
func NewHandler(service api.ServiceGenerator) Handler {
	return Handler{
		impl: service,
	}
}

// Handle receives an enveloped request for ServiceGenerator service and
// returns an enveloped response.
func (h Handler) Handle(envelope wire.Envelope) (wire.Envelope, error) {
	responseEnvelope := wire.Envelope{
		Name:  envelope.Name,
		Type:  wire.Reply,
		SeqID: envelope.SeqID,
	}

	switch envelope.Name {

	case "generate":
		var args GenerateArgs
		if err := args.FromWire(envelope.Value); err != nil {
			return responseEnvelope, err
		}

		result, err := GenerateHelper.WrapResponse(
			h.impl.Generate(args.Request),
		)
		if err != nil {
			return responseEnvelope, err
		}

		responseEnvelope.Value, err = result.ToWire()
		if err != nil {
			return responseEnvelope, err
		}

	default:

		// TODO(abg): Use TException
		return responseEnvelope, fmt.Errorf("unknown method %q", envelope.Name)

	}

	return responseEnvelope, nil
}
