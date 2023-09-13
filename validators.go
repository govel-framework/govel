package govel

import (
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
)

/*
 * General validators.
 */

func validatorRequiredField(key string, value string, form formInterface) error {
	if value == "" {
		return &validationError{Format: "%s is required", Key: key}
	}

	form.setVariableForValidation(value)

	return nil
}

func validatorInt(key string, value string, form formInterface) error {
	number, err := strconv.Atoi(value)

	if err != nil {
		return &validationError{Format: "%s value must be an integer", Key: key}
	}

	form.setVariableForValidation(number)

	return nil
}

func validatorString(key string, value string, form formInterface) error {

	form.setVariableForValidation(value)

	return nil
}

func validatorMin(key string, value string, form formInterface) error {
	value_as_int, err := strconv.Atoi(value)

	if err != nil {
		panic(fmt.Sprintf("The key %s is not valid", key))
	}

	switch validationVarValue := form.getVariableForValidation().(type) {
	case string:
		if len(validationVarValue) < value_as_int {
			return &validationError{Format: "%s must have at least " + value + " characters", Key: key}
		}

	case int:
		if validationVarValue < value_as_int {
			return &validationError{Format: "%s cannot be less than " + value, Key: key}
		}

	}

	return nil
}

func validatorMax(key string, value string, form formInterface) error {
	value_as_int, err := strconv.Atoi(value)

	if err != nil {
		panic(fmt.Sprintf("The key %s is not valid", key))
	}

	switch validationVarValue := form.getVariableForValidation().(type) {
	case string:
		if len(validationVarValue) > value_as_int {
			return &validationError{Format: "%s cannot be longer than " + value + " characters", Key: key}
		}

	case int:
		if validationVarValue > value_as_int {
			return &validationError{Format: "%s cannot be greater than " + value, Key: key}
		}

	}

	return nil
}

func validatorEmail(key string, value string, form formInterface) error {
	_, err := mail.ParseAddress(value)

	if err != nil {
		return &validationError{Format: "%s is not a valid email address", Key: key}
	}

	return nil
}

func validatorOptionalField(key string, value string, form formInterface) error {
	if value == "" {
		form.setSkipValidation(true)
	}

	form.setVariableForValidation(value)

	return nil
}

func validatorUrl(key string, value string, form formInterface) error {
	url, err := url.ParseRequestURI(value)

	if err != nil || url.Scheme == "" || url.Host == "" || (strings.Index(url.Host, ".") <= 0 || strings.LastIndex(url.Host, ".") == len(url.Host)-1) {
		return &validationError{Format: "%s is not a valid url", Key: key}
	}

	form.setVariableForValidation(url)

	return nil
}

func validatorCofirm(key string, value string, f formInterface) error {
	var confirmValue string

	form, isForm := f.(*Form)

	if isForm {
		confirmValue = form.Get(fmt.Sprintf("%s_confirm", key))
	}

	multiPartForm, isMultipartForm := f.(*MultipartFormData)

	if isMultipartForm {
		confirmValue = multiPartForm.Get(fmt.Sprintf("%s_confirm", key))
	}

	if value != confirmValue {
		e := fmt.Sprintf("%s and %s do not match", value, confirmValue)

		return &validationError{Format: e}
	}

	return nil
}

func validatorAlphaNum(key string, value string, f formInterface) error {
	var isAlphaNum bool = true

	for _, char := range value {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != ' ' {
			isAlphaNum = false
			break
		}
	}

	if !isAlphaNum {
		return &validationError{Format: "%s is not alphanumeric", Key: key}
	}

	return nil
}

// Date validators

func validatorDate(key string, value string, form formInterface) error {
	format := make(map[string]string)
	format["year"] = "2006"
	format["month"] = "01"
	format["day"] = "02"
	format["hour"] = "15"
	format["minute"] = "04"
	format["second"] = "05"

	// format the date
	for key, keyValue := range format {
		value = strings.ReplaceAll(value, key, keyValue)
	}

	date, err := time.Parse(value, form.getVariableForValidation().(string))

	if err != nil {
		return &validationError{Format: "%s is not a valid date", Key: key}
	}

	form.setVariableForValidation(date)

	return nil
}

// Boolean validators

func validatorBoolean(key string, value string, form formInterface) error {
	value_to_boolean, err := strconv.ParseBool(value)

	if err != nil {
		return &validationError{Format: "%s is not valid", Key: key}
	}

	form.setVariableForValidation(value_to_boolean)

	return nil
}

func validatorTrue(key string, value string, form formInterface) error {
	if form.getVariableForValidation().(bool) == false {
		return &validationError{Format: "%s is not valid", Key: key}
	}

	return nil
}

/*
 * MultipartFormData validators
 */

func validatorOptionalFile(key string, value string, form formInterface) error {
	f, _ := form.(*MultipartFormData)

	file, err := f.GetFile(key)

	if err != nil {
		f.setSkipValidation(true)
	}

	form.setVariableForValidation(file)

	return nil
}

func validatorRequiredFile(key string, value string, form formInterface) error {
	f, _ := form.(*MultipartFormData)

	file, err := f.GetFile(key)

	if err != nil {
		return &validationError{Format: "%s file is required", Key: key}
	}

	form.setVariableForValidation(file)

	return nil
}

func validatorContentType(key string, value string, form formInterface) error {

	// validate the file
	content, _ := io.ReadAll(form.getVariableForValidation().(FormFile).File)

	contentType := http.DetectContentType(content)

	// validate the contentType
	contentTypes := strings.Split(value, ",")

	isValidContentType := false

	for _, ct := range contentTypes {
		if contentType == ct {
			isValidContentType = true
			break
		}
	}

	if !isValidContentType {
		return &validationError{Format: "File %s does not have a valid content type", Key: key}
	}

	return nil
}

func validatorMaxBytes(key string, value string, form formInterface) error {
	formFile := form.getVariableForValidation().(FormFile)

	value_as_int, err := strconv.Atoi(value)

	if err != nil {
		panic(fmt.Sprintf("\"%s\" is not valid for maxBytes", value))
	}

	if formFile.FileHeader.Size > int64(value_as_int) {
		return &validationError{Format: "%s is too large", Key: key}
	}

	return nil
}
