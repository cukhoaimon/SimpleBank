package val

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\-\\s]+$`).MatchString
)

func ValidateString(value string, min, max int) error {
	n := len(value)
	if n <= min || n >= max {
		return fmt.Errorf("must have [from %d to %d] characters", min, max)
	}

	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 8, 50); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return errors.New("username must contains only letters (uppercase or lowercase), digits and underscore")
	}

	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 8, 100); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return errors.New("full name must contains only letters, spaces and hyphen")
	}

	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 8, 50); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return errors.New("invalid email address")
	}

	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 8, 50)
}
