package cronjobs

import (
	"time"

	"github.com/piero0920/archiv-try/pkg/external_apis"
	"github.com/piero0920/archiv-try/pkg/logger"
	"github.com/piero0920/archiv-try/pkg/models"
	"github.com/piero0920/archiv-try/pkg/queries"
)

func UpdateEmotes() {
	logger.Debug.Println("[cronjob] Updating all emotes")

	// mark all emotes outdated
	if err := queries.MarkAllEmotesOutdated(true); err != nil {
		logger.Error.Println("[cronjob] mark all emotes outdated failed:")
		logger.Error.Println(err)
	}

	// update twitch emotes
	if err := external_apis.TwitchUpdateEmotes(); err != nil {
		logger.Error.Println("[cronjob] update twitch emotes failed:")
		logger.Error.Println(err)
	}

	// update bttv emotes
	if err := external_apis.BttvUpdateEmotes(); err != nil {
		logger.Error.Println("[cronjob] update bttv emotes failed:")
		logger.Error.Println(err)
	}

	// update ffz emotes
	if err := external_apis.FfzUpdateEmotes(); err != nil {
		logger.Error.Println("[cronjob] update ffz emotes failed:")
		logger.Error.Println(err)
	}

	// delete all outdated emotes
	if err := queries.DeleteOutdatedEmotes(); err != nil {
		logger.Error.Println("[cronjob] delete all outdated emotes failed:")
		logger.Error.Println(err)
	}

	// write update time to settings
	var settings models.Settings
	settings.DateEmotesUpdate = time.Now()
	if err := queries.PartiallyUpdateSettings(&settings); err != nil {
		logger.Error.Println("[cronjob] write update time to settings failed:")
		logger.Error.Println(err)
	}
}
