package plugin

import (
	"bytes"
	"fmt"
	"github.com/thriftrw/thriftrw-go/protocol"
	"github.com/thriftrw/thriftrw-go/plugin/api"
	"github.com/thriftrw/thriftrw-go/wire"
	"github.com/thriftrw/thriftrw-go/envelope"
)

// Protocol is the Thrift protocol used to deserialize and serialize requests.
var Protocol = protocol.Binary

// Interface provides the Plugin service.
type Interface interface {
	Generate(
		Request *api.GenerateRequest,
	) (*api.GenerateResponse, error)

	Goodbye() error

	Handshake(
		Request *api.HandshakeRequest,
	) (*api.HandshakeResponse, error)
}

// Client implements a Plugin client.
type client struct {
	send func([]byte) ([]byte, error)
}

// NewClient builds a new Plugin client.
func NewClient(t func([]byte) ([]byte, error)) Interface {
	return &client{
		send: t,
	}
}

func (c *client) Generate(
	_Request *api.GenerateRequest,
) (success *api.GenerateResponse, err error) {
	args := GenerateHelper.Args(_Request)

	var buff bytes.Buffer
	if err = envelope.Write(Protocol, &buff, 1, args); err != nil {
		return
	}

	var resBody []byte
	resBody, err = c.send(buff.Bytes())
	if err != nil {
		return
	}

	var body wire.Value
	body, _, err = envelope.ReadReply(Protocol, bytes.NewReader(resBody))
	if err != nil {
		return
	}

	var result GenerateResult
	if err = result.FromWire(body); err != nil {
		return
	}

	success, err = GenerateHelper.UnwrapResponse(&result)
	return
}

func (c *client) Goodbye() (err error) {
	args := GoodbyeHelper.Args()

	var buff bytes.Buffer
	if err = envelope.Write(Protocol, &buff, 1, args); err != nil {
		return
	}

	var resBody []byte
	resBody, err = c.send(buff.Bytes())
	if err != nil {
		return
	}

	var body wire.Value
	body, _, err = envelope.ReadReply(Protocol, bytes.NewReader(resBody))
	if err != nil {
		return
	}

	var result GoodbyeResult
	if err = result.FromWire(body); err != nil {
		return
	}

	err = GoodbyeHelper.UnwrapResponse(&result)
	return
}

func (c *client) Handshake(
	_Request *api.HandshakeRequest,
) (success *api.HandshakeResponse, err error) {
	args := HandshakeHelper.Args(_Request)

	var buff bytes.Buffer
	if err = envelope.Write(Protocol, &buff, 1, args); err != nil {
		return
	}

	var resBody []byte
	resBody, err = c.send(buff.Bytes())
	if err != nil {
		return
	}

	var body wire.Value
	body, _, err = envelope.ReadReply(Protocol, bytes.NewReader(resBody))
	if err != nil {
		return
	}

	var result HandshakeResult
	if err = result.FromWire(body); err != nil {
		return
	}

	success, err = HandshakeHelper.UnwrapResponse(&result)
	return
}

// Handler serves an implementation of the Plugin service.
type Handler struct {
	impl Interface
}

// NewHandler builds a new Plugin handler.
func NewHandler(service Interface) Handler {
	return Handler{
		impl: service,
	}
}

// Handle handles the given request and returns a response.
func (h Handler) Handle(data []byte) ([]byte, error) {
	envelope, err := Protocol.DecodeEnveloped(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	responseEnvelope, err := h.HandleEnvelope(envelope)
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer
	if err := Protocol.EncodeEnveloped(responseEnvelope, &buff); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// HandleEnvelope receives an enveloped request for Plugin service
// and returns an enveloped response.
func (h Handler) HandleEnvelope(envelope wire.Envelope) (wire.Envelope, error) {
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
