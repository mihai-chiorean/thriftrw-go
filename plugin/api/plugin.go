package api

type Plugin interface {
	Generate(
		Request *GenerateRequest,
	) (*GenerateResponse, error)

	Goodbye() error

	Handshake(
		Request *HandshakeRequest,
	) (*HandshakeResponse, error)
}
