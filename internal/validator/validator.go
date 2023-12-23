package validator

import (
	"fmt"
	"test-signer.stekels.lv/internal/transport"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) Validate(input transport.CreateSignatureRequest) {
	if input.UserJWT == "" {
		v.AddError("user_jwt", "must be provided")
	}
	if input.Questions == nil {
		v.AddError("questions", "must be provided")
	}
	if len(input.Questions) == 0 {
		v.AddError("questions", "can not be empty")
	}
	for index, question := range input.Questions {
		if question.Body == "" {
			v.AddError(fmt.Sprintf("question[%d]", index), "question body must be provided")
		}
		if question.Answer == "" {
			v.AddError(fmt.Sprintf("question[%d]", index), "question answer must be provided")
		}
	}
}
