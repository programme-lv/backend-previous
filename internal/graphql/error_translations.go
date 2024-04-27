package graphql

import "fmt"

func ErrUserNotFound(lang string) error {
	switch lang {
	case "lv":
		return fmt.Errorf("lietotājs nav atrasts")
	default:
		return fmt.Errorf("user not found")
	}
}

func ErrUsernameOrPasswordIncorrect(lang string) error {
	if lang == "lv" {
		return fmt.Errorf("lietotājvārds vai parole ir nepareiza")
	}
	return fmt.Errorf("username or password is incorrect")
}
