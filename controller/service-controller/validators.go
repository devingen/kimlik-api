package service_controller

import "regexp"

var emailRegex *regexp.Regexp

func IsValidEmail(email string) bool {
	if emailRegex == nil {
		emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	}
	return emailRegex.MatchString(email)
}
