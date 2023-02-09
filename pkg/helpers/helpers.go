package helpers

import (
	"os"

	"github.com/piero0920/archiv-try/pkg/external_apis"
	"github.com/piero0920/archiv-try/pkg/models"
	"github.com/piero0920/archiv-try/pkg/queries"
)

func ImportEnvToDb() error {
	updatedSettings := false

	// get settings and create if none found
	var settings models.Settings
	if err := queries.GetSettings(&settings); err != nil {
		settings.TtvClientId = os.Getenv("TTV_CLIENT_ID")
		settings.TtvClientSecret = os.Getenv("TTV_CLIENT_SECRET")
		if errInit := queries.InitSettings(&settings); errInit != nil {
			return errInit
		}
		updatedSettings = true
	}

	// update ttv keys if not already created
	if !updatedSettings {
		settings.TtvClientId = os.Getenv("TTV_CLIENT_ID")
		settings.TtvClientSecret = os.Getenv("TTV_CLIENT_SECRET")
		if err := queries.PartiallyUpdateSettings(&settings); err != nil {
			return err
		}
	}

	// get broadcaster id from name
	if err := external_apis.TwitchUpdateBroadcasterID(); err != nil {
		return err
	}

	return nil
}
