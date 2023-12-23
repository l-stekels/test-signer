package transport

type CreateSignatureRequest struct {
	UserJWT   string     `json:"user_jwt"`
	Questions []Question `json:"questions"`
}
type Question struct {
	Body   string `json:"question"`
	Answer string `json:"answer"`
}

type CreateSignatureResponse struct {
	Data Signature `json:"data"`
}
type Signature struct {
	Signature string `json:"signature"`
}

func NewCreateSignatureResponse(signature string) *CreateSignatureResponse {
	return &CreateSignatureResponse{
		Data: Signature{
			Signature: signature,
		},
	}
}
