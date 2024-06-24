package logs

import (
	"log/slog"
)

func LogCommandExecution(commandName string, cmd interface{}, err error) {
	log := slog.With("cmd", cmd)

	if err == nil {
		log.Info(commandName + " command succeeded")
	} else {
		log.Error(commandName+" command failed", "error", err)
	}
}
