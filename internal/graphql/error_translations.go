package graphql

import "fmt"

func ErrUsernameOrPasswordIncorrect(lang string) error {
	if lang == "lv" {
		return fmt.Errorf("lietotājvārds vai parole ir nepareiza")
	}
	return fmt.Errorf("username or password is incorrect")
}

func ErrInternalServer(lang string) error {
	if lang == "lv" {
		return fmt.Errorf("iekšējā servera kļūda")
	}
	return fmt.Errorf("internal server error")
}

func ErrUsernameOrPasswordEmpty(lang string) error {
	if lang == "lv" {
		return fmt.Errorf("lietotājvārds un parole ir obligāti")
	}
	return fmt.Errorf("username and password are required")
}

func ErrPasswordTooShort(lang string, min int) error {
	if lang == "lv" {
		return fmt.Errorf("parolei jābūt vismaz %d simbolus garai", min)
	}
	return fmt.Errorf("password must be at least %d characters", min)
}

func ErrPasswordTooLong(lang string, max int) error {
	if lang == "lv" {
		return fmt.Errorf("parolei jābūt ne vairāk kā %d simbolus garai", max)
	}
	return fmt.Errorf("password must be at most %d characters", max)
}

func ErrUsernameTooShort(lang string, min int) error {
	if lang == "lv" {
		return fmt.Errorf("lietotājvārdam jābūt vismaz %d simbolus garai", min)
	}
	return fmt.Errorf("username must be at least %d characters", min)
}

func ErrUsernameTooLong(lang string, max int) error {
	if lang == "lv" {
		return fmt.Errorf("lietotājvārdam jābūt ne vairāk kā %d simbolus garai", max)
	}
	return fmt.Errorf("username must be at most %d characters", max)
}

func ErrUserWithThatUsernameExists(lang string) error {
	if lang == "lv" {
		return fmt.Errorf("lietotājs ar šādu lietotājvārdu jau pastāv")
	}
	return fmt.Errorf("user with that username already exists")
}

func ErrUserWithThatEmailExists(lang string) error {
	if lang == "lv" {
		return fmt.Errorf("lietotājs ar šādu e-pastu jau pastāv")
	}
	return fmt.Errorf("user with that email already exists")
}

func ErrInvalidEmail(lang string) error {
	if lang == "lv" {
		return fmt.Errorf("nederīgs e-pasts")
	}
	return fmt.Errorf("invalid email")
}
