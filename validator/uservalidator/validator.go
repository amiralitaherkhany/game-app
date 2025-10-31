package uservalidator

import (
	"fmt"
	"gameapp/dto"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
}
type Validator struct {
	repo Repository
}

func New(repo Repository) *Validator {
	return &Validator{repo: repo}
}

func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateRegisterRequest"

	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.Name,
			validation.Required,
			validation.Length(3, 50),
		),
		// Minimum eight characters, at least one letter, one number and one special character:
		validation.Field(
			&req.Password,
			validation.Required,
			validation.Match(
				regexp.MustCompile(`^[a-zA-Z0-9@$%^#!*]{8,}$`),
			),
		),
		validation.Field(
			&req.PhoneNumber,
			validation.Required,
			validation.Match(
				regexp.MustCompile("^09[0-9]{9}$"),
			),
			validation.By(v.checkPhoneNumberUniqueness),
		),
	)
	if err != nil {
		fieldErrors := make(map[string]string)

		errors, ok := err.(validation.Errors)
		if ok {
			for k, err := range errors {
				fieldErrors[k] = err.Error()
			}
		}

		return fieldErrors, richerror.
			New(op).
			WithMessage(errmsg.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).
			WithErr(err)
	}

	return nil, nil
}

func (v Validator) checkPhoneNumberUniqueness(value any) error {
	phoneNumber := value.(string)
	isUnique, err := v.repo.IsPhoneNumberUnique(phoneNumber)
	if err != nil {
		return err
	}
	if !isUnique {
		return fmt.Errorf(errmsg.ErrorMsgPhoneNumberIsNotUnique)
	}
	return nil
}
