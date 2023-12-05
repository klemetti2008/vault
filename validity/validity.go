package validity

import (
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

func ApplyTranslations() {
	govalidity.SetDefaultErrorMessages(
		&govaliditym.Validations{
			IsEmail:           "{field} must be valid email",
			IsUrl:             "{field} must be valid url",
			IsIp:              "{field} must be valid ip",
			IsRequired:        "{field} should have values",
			IsMinLength:       "{field} should have minimum {min} characters",
			IsMin:             "{field} should be minimum {min}",
			IsMaxLength:       "{field} should have maximum {max} characters",
			IsMax:             "{field} should be {max}",
			IsMinMaxLength:    "{field} should be at least {min} and at most {max} characters",
			IsAlpha:           "{field} should consist of alphabetic characters",
			IsAlphaNum:        "{field} should consist of alphabetic and numeric characters",
			IsNumber:          "{field} should be number",
			IsMaxDate:         "{field} should be at most {max} characters",
			IsDNSName:         "{field} should be valid domain name",
			IsFloat:           "{field} should be floating-point number",
			IsFilepath:        "{field} should be valid file path",
			IsHost:            "{field} should be valid host",
			IsIn:              "{field} should be one of the values {in}",
			IsInt:             "{field} should be an integer",
			IsInRange:         "{field} should be between {from} and {to}",
			IsIpV4:            "{field} should be an IPv4 address",
			IsIpV6:            "{field} should be an IPv6 address",
			IsJson:            "{field} should be a valid JSON string",
			IsLatitude:        "{field} should be a valid latitude",
			IsLogitude:        "{field} should be a valid longitude",
			IsLowerCase:       "{field} should be in lowercase",
			IsMaxDateTime:     "{field} should be at most {max}",
			IsFilterOperators: "{field} should be one of the filter operators {in}",
			IsHexColor:        "{field} should be a valid color code",
			IsMaxTime:         "{field} should be at most {max}",
			IsMinDate:         "{field} should be at least {min}",
			IsPort:            "{field} should be a valid port",
			IsMinDateTime:     "{field} should be at least {min}",
			IsMinTime:         "{field} should be at least {min}",
			IsSlice:           "{field} should be an array",
			IsUpperCase:       "{field} should be in uppercase",
		},
	)
}
