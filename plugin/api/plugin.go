package api

type Plugin interface {
	Goodbye() error

	Handshake(
		Request *HandshakeRequest,
	) (*HandshakeResponse, error)
}
