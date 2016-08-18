package api

type ServiceGenerator interface {
	Generate(
		Request *GenerateServiceRequest,
	) (*GenerateServiceResponse, error)
}
