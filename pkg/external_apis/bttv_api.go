package external_apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/piero0920/archiv-try/pkg/models"
	"github.com/piero0920/archiv-try/pkg/queries"
)

type bttvEmote struct {
	ID   string
	Code string
}

type bttvEmoteResponse struct {
	SharedEmotes []bttvEmote
}

func BttvUpdateEmotes() error {
	var settings models.Settings
	if err := queries.GetSettings(&settings); err != nil {
		return err
	}

	url := "https://api.betterttv.net/3/cached/users/twitch/" + settings.BroadcasterId
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bttv emotes: status code was %d", resp.StatusCode)
	}

	var responseJson bttvEmoteResponse
	json.NewDecoder(resp.Body).Decode(&responseJson)

	// iterate emotes and save in db
	for _, emote := range responseJson.SharedEmotes {
		newEmote := models.Emote{
			ID:       emote.ID,
			Name:     emote.Code,
			URL:      fmt.Sprintf("https://cdn.betterttv.net/emote/%s/3x", emote.ID),
			Provider: "bttv",
			Outdated: false,
		}
		queries.UpdateOrCreateEmote(&newEmote, emote.ID)
	}

	return nil
}
