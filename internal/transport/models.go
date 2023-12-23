package transport

import "time"

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

type VerifySignatureRequest struct {
	UserJWT   string `json:"user_jwt"`
	Signature string `json:"signature"`
}

type VerifySignatureResponse struct {
	Data VerifiedSignature `json:"data"`
}
type VerifiedSignature struct {
	Answers   []Question `json:"answers"`
	Timestamp time.Time  `json:"timestamp"`
}

func NewVerifySignatureResponse(questions []Question, time time.Time) *VerifySignatureResponse {
	return &VerifySignatureResponse{
		Data: VerifiedSignature{
			Answers:   questions,
			Timestamp: time,
		},
	}
}
