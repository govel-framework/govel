package govel

import (
	"net/url"
)

/*
 * Private methods
 */

func (f *Form) skipValidation() bool {
	return f.skip
}

func (f *Form) setSkipValidation(value bool) {
	f.skip = value
}

func (f *Form) setVariableForValidation(v interface{}) {
	f.validationVar = v
}

func (f *Form) getVariableForValidation() interface{} {
	return f.validationVar
}

/*
 * Public methods
 */

// Gets all POST data of the form
func (f *Form) GetAll() url.Values {
	return f.request.Form
}

// Gets the first value associated with the given key
func (f *Form) Get(value string) string {
	return f.request.Form.Get(value)
}

// Validate validates a form with the provided rules.
func (f *Form) Validate(rules Map, onError OnError) (data Map, errors map[string]string) {
	// All the valid validation rules
	formRules := map[string]validateFunc{
		"required":  validatorRequiredField,
		"int":       validatorInt,
		"string":    validatorString,
		"min":       validatorMin,
		"max":       validatorMax,
		"email":     validatorEmail,
		"optional":  validatorOptionalField,
		"url":       validatorUrl,
		"date":      validatorDate,
		"true":      validatorTrue,
		"boolean":   validatorBoolean,
		"confirm":   validatorCofirm,
		"alpha_num": validatorAlphaNum,
	}

	// start validation
	data = make(map[string]interface{})
	errors = make(map[string]string)

	formValidator(formRules, data, errors, rules, onError, f)

	if len(errors) > 0 {
		return nil, errors
	}

	return data, nil
}
