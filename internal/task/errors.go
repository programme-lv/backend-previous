package task

import (
	"github.com/programme-lv/backend/internal/comm/i18nerror"
	"golang.org/x/text/language"
)

func newErrorUserDoesNotHaveEditAccessToTask() i18nerror.I18NError {
	return i18nerror.New("err_user_does_not_have_edit_access_to_task", map[language.Tag]string{
		language.English: "user does not have edit access to task",
		language.Latvian: "lietotājam nav rediģēšanas piekļuves uzdevumam",
	}, nil)
}

func newErrorTaskNotFound() i18nerror.I18NError {
	return i18nerror.New("err_task_not_found", map[language.Tag]string{
		language.English: "task not found",
		language.Latvian: "uzdevums nav atrasts",
	}, nil)
}

func newErrorNoStableVersion() i18nerror.I18NError {
	return i18nerror.New("err_no_stable_version", map[language.Tag]string{
		language.English: "task does not have a stable version",
		language.Latvian: "uzdevumam nav stabilas versijas",
	}, nil)
}
