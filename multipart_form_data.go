package govel

/*
 * Private methods
 */

func (f *MultipartFormData) skipValidation() bool {
	return f.skip
}

func (f *MultipartFormData) setSkipValidation(skip bool) {
	f.skip = skip
}

func (f *MultipartFormData) setVariableForValidation(v interface{}) {
	f.validationVar = v
}

func (f *MultipartFormData) getVariableForValidation() interface{} {
	return f.validationVar
}

/*
 * Public methods
 */

// Gets all POST data of the form except files
func (f *MultipartFormData) GetAll() interface{} {
	return f.request.MultipartForm.Value
}

// Gets the first value associated with the given key except files.
func (f *MultipartFormData) Get(value string) string {
	return f.request.FormValue(value)
}

// Gets the first file associated with the given key.
func (f *MultipartFormData) GetFile(value string) (FormFile, error) {
	file, header, err := f.request.FormFile(value)

	return FormFile{File: file, FileHeader: header}, err
}

func (f *MultipartFormData) Validate(rules Map, onError OnError) (data map[string]interface{}, errors map[string]string) {
	// All the valid validation rules
	formRules := map[string]validateFunc{
		"required":     validatorRequiredField,
		"int":          validatorInt,
		"string":       validatorString,
		"min":          validatorMin,
		"max":          validatorMax,
		"email":        validatorEmail,
		"optional":     validatorOptionalField,
		"optionalFile": validatorOptionalFile,
		"requiredFile": validatorRequiredFile,
		"contentType":  validatorContentType,
		"url":          validatorUrl,
		"maxBytes":     validatorMaxBytes,
		"date":         validatorDate,
		"true":         validatorTrue,
		"boolean":      validatorBoolean,
		"confirm":      validatorCofirm,
		"alpha_num":    validatorAlphaNum,
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
