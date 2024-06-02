package domain

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
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
		bundle.AddMessages(lang, &i18n.Message{
			ID:    e.I18NMessageID,
			Other: translation,
		})
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
