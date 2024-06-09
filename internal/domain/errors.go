package domain

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log/slog"
)

type PublicI18NError struct {
	I18NMessageID string
	Translations  *map[language.Tag]string
	Language      *string
	TemplateData  *map[string]interface{}
}

func (e *PublicI18NError) Error() string {
	if e.Language != nil {
		return e.Localize(*e.Language)
	} else {
		return e.Localize("en")
	}
}

func (e *PublicI18NError) WithLanguage(lang string) *PublicI18NError {
	return &PublicI18NError{
		I18NMessageID: e.I18NMessageID,
		Translations:  e.Translations,
		Language:      &lang,
		TemplateData:  e.TemplateData,
	}
}

func (e *PublicI18NError) Localize(langs ...string) string {
	bundle := i18n.NewBundle(language.English)
	for lang, translation := range *e.Translations {
		err := bundle.AddMessages(lang, &i18n.Message{
			ID:    e.I18NMessageID,
			Other: translation,
		})
		if err != nil {
			slog.Error(err.Error())
			return ""
		}
	}

	return i18n.NewLocalizer(bundle, langs...).MustLocalize(&i18n.LocalizeConfig{
		MessageID:    e.I18NMessageID,
		TemplateData: *e.TemplateData,
	})
}

const errorPasswordTooLongID = "err_password_too_long"

var errorPasswordTooLongTranslations = map[language.Tag]string{
	language.English: "password must be at most {{.max_length}} characters long",
	language.Latvian: "parolei jābūt ne vairāk kā {{.max_length}} simbolus garai",
}

func NewErrorPasswordTooLong(maxLength int) *PublicI18NError {
	return &PublicI18NError{
		I18NMessageID: errorPasswordTooLongID,
		Translations:  &errorPasswordTooLongTranslations,
		TemplateData:  &map[string]interface{}{"max_length": maxLength},
	}
}

const errorPasswordTooShortID = "err_password_too_short"

var errorPasswordTooShortTranslations = map[language.Tag]string{
	language.English: "password must be at least {{.min_length}} characters long",
	language.Latvian: "parolei jābūt vismaz {{.min_length}} simbolus garai",
}

func NewErrorPasswordTooShort(minLength int) *PublicI18NError {
	return &PublicI18NError{
		I18NMessageID: errorPasswordTooShortID,
		Translations:  &errorPasswordTooShortTranslations,
		TemplateData:  &map[string]interface{}{"min_length": minLength},
	}
}

const errorNotLoggedInID = "err_not_logged_in"

var errorNotLoggedInTranslations = map[language.Tag]string{
	language.English: "user is not logged in",
	language.Latvian: "lietotājs nav pieslēdzies",
}

func NewErrorNotLoggedIn() *PublicI18NError {
	return &PublicI18NError{
		I18NMessageID: errorNotLoggedInID,
		Translations:  &errorNotLoggedInTranslations,
	}
}

const errorInternalServerID = "err_internal_server"

var errorInternalServerTranslations = map[language.Tag]string{
	language.English: "internal server error",
	language.Latvian: "iekšējā servera kļūda",
}

func NewErrorInternalServer() *PublicI18NError {
	return &PublicI18NError{
		I18NMessageID: errorInternalServerID,
		Translations:  &errorInternalServerTranslations,
	}
}

const errorUsernameOrPasswordEmptyID = "err_username_or_password_empty"

var errorUsernameOrPasswordEmptyTranslations = map[language.Tag]string{
	language.English: "username or password is empty",
	language.Latvian: "lietotājvārds vai parole ir tukša",
}

func NewErrorUsernameOrPasswordEmpty() *PublicI18NError {
	return &PublicI18NError{
		I18NMessageID: errorUsernameOrPasswordEmptyID,
		Translations:  &errorUsernameOrPasswordEmptyTranslations,
	}
}

const errorUsernameTooShortID = "err_username_too_short"

var errorUsernameTooShortTranslations = map[language.Tag]string{
	language.English: "username must be at least {{.min_length}} characters long",
	language.Latvian: "lietotājvārdam jābūt vismaz {{.min_length}} simbolus garam",
}

func NewErrorUsernameTooShort(minLength int) *PublicI18NError {
	return &PublicI18NError{
		I18NMessageID: errorUsernameTooShortID,
		Translations:  &errorUsernameTooShortTranslations,
		TemplateData:  &map[string]interface{}{"min_length": minLength},
	}
}

const errorUsernameTooLongID = "err_username_too_long"

var errorUsernameTooLongTranslations = map[language.Tag]string{
	language.English: "username must be at most {{.max_length}} characters long",
	language.Latvian: "lietotājvārdam jābūt ne vairāk kā {{.max_length}} simbolus garam",
}

func NewErrorUsernameTooLong(maxLength int) *PublicI18NError {
	return &PublicI18NError{
		I18NMessageID: errorUsernameTooLongID,
		Translations:  &errorUsernameTooLongTranslations,
		TemplateData:  &map[string]interface{}{"max_length": maxLength},
	}
}

const errorUserWithUsernameExistsID = "err_user_with_username_exists"

var errorUserWithUsernameExistsTranslations = map[language.Tag]string{
	language.English: "user with this username already exists",
	language.Latvian: "lietotājs ar šādu lietotājvārdu jau eksistē",
}

func NewErrorUserWithUsernameExists() *PublicI18NError {
	return &PublicI18NError{
		I18NMessageID: errorUserWithUsernameExistsID,
		Translations:  &errorUserWithUsernameExistsTranslations,
	}
}

const errorUserWithEmailExistsID = "err_user_with_email_exists"

var errorUserWithEmailExistsTranslations = map[language.Tag]string{
	language.English: "user with this email already exists",
	language.Latvian: "lietotājs ar šādu e-pastu jau eksistē",
}

func NewErrorUserWithEmailExists() *PublicI18NError {
	return &PublicI18NError{
		I18NMessageID: errorUserWithEmailExistsID,
		Translations:  &errorUserWithEmailExistsTranslations,
	}
}

const errorInvalidEmailID = "err_invalid_email"

var errorInvalidEmailTranslations = map[language.Tag]string{
	language.English: "invalid email",
	language.Latvian: "nederīgs e-pasts",
}

func NewErrorInvalidEmail() *PublicI18NError {
	return &PublicI18NError{
		I18NMessageID: errorInvalidEmailID,
		Translations:  &errorInvalidEmailTranslations,
	}
}

//const errorUsernameOrPasswordIncorrectID = "err_username_or_password_incorrect"
//
//var errorUsernameOrPasswordIncorrectTranslations = map[language.Tag]string{
//	language.English: "username or password is incorrect",
//	language.Latvian: "lietotājvārds vai parole ir nepareiza",
//}
//
//func NewErrorUsernameOrPasswordIncorrect() *PublicI18NError {
//	return &PublicI18NError{
//		I18NMessageID: errorUsernameOrPasswordIncorrectID,
//		Translations:  &errorUsernameOrPasswordIncorrectTranslations,
//	}
//}
