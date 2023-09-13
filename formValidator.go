package govel

import (
	"fmt"
	"regexp"
	"strings"
)

// Main function for validations.
func formValidator(formRules map[string]validateFunc, data Map, errors SMap, rules Map, onError map[string]SMap, f formInterface) {
	for key, value := range rules {
		var rules []string

		// split the rules or get the rules
		rulesString, isString := value.(string)

		if isString {
			rules = strings.Split(rulesString, "|")
		}

		rulesSlice, isSlice := value.([]string)

		if isSlice {
			rules = rulesSlice
		}

		if !isSlice && !isString /* Is not a valid type */ {
			panic(fmt.Sprintf("The data type of the rules is not valid, only []string and string are allowed, but the type is: %T", value))
		}

		// now we iterate over the rules
		for _, rule := range rules {
			if f.skipValidation() {
				f.setSkipValidation(false)
				break
			}

			// here we check if the rule is of type key:value
			regex := regexp.MustCompile(`^[a-zA-Z]+:[^:\n]*$`)

			if regex.Match([]byte(rule)) {
				new_rule := strings.Split(rule, ":")

				if len(new_rule) != 2 {
					panic(fmt.Sprintf("Rule %s not valid", rule))
				}

				rule = new_rule[0]

				validationCallable, exists := formRules[rule]

				if !exists {
					panic(fmt.Sprintf("Rule %s not found.", rule))
				}

				err := validationCallable(key, new_rule[1], f)

				if err != nil {
					// check if there is an error message for this case
					error_msg, exists := onError[key][rule]

					if !exists {
						general_error, exists := onError[key]["*"]

						if !exists {
							errors[key] = err.Error()
						} else {
							errors[key] = general_error
						}

					} else {
						errors[key] = error_msg
					}

					break
				}

				if f.skipValidation() == false {
					data[key] = f.getVariableForValidation()
				}

				continue
			}

			validationCallable, exists := formRules[rule]

			if !exists {
				panic(fmt.Sprintf("Rule %s not found.", rule))
			}

			// do the validation
			formValue := strings.TrimSpace(f.Get(key))

			err := validationCallable(key, formValue, f)

			if err != nil {
				// check if there is an error message for this case
				error_msg, exists := onError[key][rule]

				if !exists {
					general_error, exists := onError[key]["*"]

					if !exists {
						errors[key] = err.Error()
					} else {
						errors[key] = general_error
					}

				} else {
					errors[key] = error_msg
				}

				break
			}

			if f.getVariableForValidation() == nil {
				f.setVariableForValidation(formValue)
			}

			if f.skipValidation() == false {
				data[key] = f.getVariableForValidation()
			}
		}

	}
}
