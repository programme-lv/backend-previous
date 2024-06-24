package eval

import (
	"github.com/programme-lv/backend/internal/common/i18nerror"
	"golang.org/x/text/language"
)

func newErrorSubmissionBodyTooLarge() i18nerror.I18NError {
	return i18nerror.New("err_submission_body_too_large", map[language.Tag]string{
		language.English: "submission body too large",
		language.Latvian: "iesūtījums teksts ir pārāk garš",
	}, nil)
}

func newErrorInvalidSubmissionID() i18nerror.I18NError {
	return i18nerror.New("err_invalid_submission_id", map[language.Tag]string{
		language.English: "invalid submission id",
		language.Latvian: "nederīgs iesūtījuma id",
	}, nil)
}

func newErrorInvalidSubmissionParams() i18nerror.I18NError {
	return i18nerror.New("err_invalid_submission_params", map[language.Tag]string{
		language.English: "invalid submission parameters",
		language.Latvian: "nederīgi iesūtījuma parametri",
	}, nil)
}
