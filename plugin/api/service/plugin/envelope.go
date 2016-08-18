package plugin

import (
	"fmt"
	"github.com/thriftrw/thriftrw-go/plugin/api"
	"github.com/thriftrw/thriftrw-go/wire"
)

// Client implements a Plugin client.
type client struct {
	send func(wire.Envelope) (wire.Envelope, error)
}

// NewClient builds a new Plugin client.
func NewClient(t func(wire.Envelope) (wire.Envelope, error)) api.Plugin {
	return &client{
		send: t,
	}
}

func (c *client) Goodbye() (err error) {
	args := GoodbyeHelper.Args()

	var body wire.Value
	body, err = args.ToWire()
	if err != nil {
		return
	}

	var envelope wire.Envelope
	envelope, err = c.send(wire.Envelope{
		Name:  "goodbye",
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

	var result GoodbyeResult
	if err = result.FromWire(envelope.Value); err != nil {
		return
	}

	err = GoodbyeHelper.UnwrapResponse(&result)
	return
}

func (c *client) Handshake(
	_Request *api.HandshakeRequest,
) (success *api.HandshakeResponse, err error) {
	args := HandshakeHelper.Args(_Request)

	var body wire.Value
	body, err = args.ToWire()
	if err != nil {
		return
	}

	var envelope wire.Envelope
	envelope, err = c.send(wire.Envelope{
		Name:  "handshake",
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

	var result HandshakeResult
	if err = result.FromWire(envelope.Value); err != nil {
		return
	}

	success, err = HandshakeHelper.UnwrapResponse(&result)
	return
}

// Handler serves an implementation of the Plugin service.
type Handler struct {
	impl api.Plugin
}

// NewHandler builds a new Plugin handler.
func NewHandler(service api.Plugin) Handler {
	return Handler{
		impl: service,
	}
}

// Handle receives an enveloped request for Plugin service and
// returns an enveloped response.
func (h Handler) Handle(envelope wire.Envelope) (wire.Envelope, error) {
	responseEnvelope := wire.Envelope{
		Name:  envelope.Name,
		Type:  wire.Reply,
		SeqID: envelope.SeqID,
	}

	switch envelope.Name {

	case "goodbye":
		var args GoodbyeArgs
		if err := args.FromWire(envelope.Value); err != nil {
			return responseEnvelope, err
		}

		result, err := GoodbyeHelper.WrapResponse(
			h.impl.Goodbye(),
		)
		if err != nil {
			return responseEnvelope, err
		}

		responseEnvelope.Value, err = result.ToWire()
		if err != nil {
			return responseEnvelope, err
		}

	case "handshake":
		var args HandshakeArgs
		if err := args.FromWire(envelope.Value); err != nil {
			return responseEnvelope, err
		}

		result, err := HandshakeHelper.WrapResponse(
			h.impl.Handshake(args.Request),
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
