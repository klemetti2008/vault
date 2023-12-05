package endpoints

import (
	"fmt"
	"regexp"

	"github.com/hoitek-go/govalidity"
)

var PasswordValidity = govalidity.New("password").CustomValidator(
	func(field string, args ...interface{}) (bool, error) {
		password := fmt.Sprintf("%v", args[0])

		upperCasePattern, err := regexp.Compile(`.*[A-Z].*`)
		if err != nil {
			panic(err)
		}

		if !upperCasePattern.MatchString(password) {
			return false, fmt.Errorf("%s must have uppercase letter", field)
		}

		lowerCasePattern, err := regexp.Compile(`.*[a-z].*`)
		if err != nil {
			panic(err)
		}

		if !lowerCasePattern.MatchString(password) {
			return false, fmt.Errorf("%s must have lowercase letter", field)
		}

		numberPattern, err := regexp.Compile(`.*\d.*`)
		if err != nil {
			panic(err)
		}

		if !numberPattern.MatchString(password) {
			return false, fmt.Errorf("%s must contain a number", field)
		}

		return true, nil
	},
).Required()

var PhoneNumberValidity = govalidity.New("phone").CustomValidator(
	func(field string, args ...interface{}) (bool, error) {
		phone := fmt.Sprintf("%v", args[0])

		patt, err := regexp.Compile(`^(09\d{9})$`)
		if err != nil {
			panic(err)
		}

		if !patt.MatchString(phone) {
			return false, fmt.Errorf("%s not valid", field)
		}
		return true, nil
	},
).Required()
