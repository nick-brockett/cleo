package model

import "github.com/asaskevich/govalidator"

type ClinicalNote struct {
	Text string `json:"text" valid:"required,stringlength(1|500)"`
}

func (n *ClinicalNote) Valid() (bool, error) {
	return govalidator.ValidateStruct(n)
}
