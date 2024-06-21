package user

import (
	"github.com/programme-lv/backend/internal/common/i18nerror"
	"golang.org/x/text/language"
)

func newErrorPasswordTooLong(maxLength int) i18nerror.I18NError {
	return i18nerror.New("err_password_too_long", map[language.Tag]string{
		language.English: "password must be at most {{.max_length}} characters long",
		language.Latvian: "parolei jābūt ne vairāk kā {{.max_length}} simbolus garai",
	}, map[string]interface{}{
		"max_length": maxLength,
	})
}

func newErrorPasswordTooShort(minLength int) i18nerror.I18NError {
	return i18nerror.New("err_password_too_short", map[language.Tag]string{
		language.English: "password must be at least {{.min_length}} characters long",
		language.Latvian: "parolei jābūt vismaz {{.min_length}} simbolus garai",
	}, map[string]interface{}{
		"min_length": minLength,
	})
}

func newErrorUsernameOrPasswordIncorrect() i18nerror.I18NError {
	return i18nerror.New("err_username_or_password_incorrect", map[language.Tag]string{
		language.English: "username or password is incorrect",
		language.Latvian: "lietotājvārds vai parole ir nepareiza",
	}, nil)
}

func newErrorUsernameOrPasswordEmpty() i18nerror.I18NError {
	return i18nerror.New("err_username_or_password_empty", map[language.Tag]string{
		language.English: "username or password is empty",
		language.Latvian: "lietotājvārds vai parole ir tukša",
	}, nil)
}

func newErrorUsernameTooShort(minLength int) i18nerror.I18NError {
	return i18nerror.New("err_username_too_short", map[language.Tag]string{
		language.English: "username must be at least {{.min_length}} characters long",
		language.Latvian: "lietotājvārdam jābūt vismaz {{.min_length}} simbolus garam",
	}, map[string]interface{}{
		"min_length": minLength,
	})
}

func newErrorUsernameTooLong(maxLength int) i18nerror.I18NError {
	return i18nerror.New("err_username_too_long", map[language.Tag]string{
		language.English: "username must be at most {{.max_length}} characters long",
		language.Latvian: "lietotājvārdam jābūt ne vairāk kā {{.max_length}} simbolus garam",
	}, map[string]interface{}{
		"max_length": maxLength,
	})
}

func newErrorInvalidEmail() i18nerror.I18NError {
	return i18nerror.New("err_invalid_email", map[language.Tag]string{
		language.English: "invalid email",
		language.Latvian: "nederīgs e-pasts",
	}, nil)
}

func newErrorUsernameAlreadyExists() i18nerror.I18NError {
	return i18nerror.New("err_username_exists", map[language.Tag]string{
		language.English: "username with such username already exists",
		language.Latvian: "lietotājvārds ar šādu lietotājvārdu jau eksistē",
	}, nil)
}

func newErrorEmailAlreadyExists() i18nerror.I18NError {
	return i18nerror.New("err_email_exists", map[language.Tag]string{
		language.English: "user with such email already exists",
		language.Latvian: "lietotājs ar šādu e-pastu jau eksistē",
	}, nil)
}

func newErrorUserNotFound() i18nerror.I18NError {
	return i18nerror.New("err_user_not_found", map[language.Tag]string{
		language.English: "user not found",
		language.Latvian: "lietotājs nav atrasts",
	}, nil)
}
