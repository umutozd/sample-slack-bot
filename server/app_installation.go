package server

import (
	"fmt"
	"net/http"

	"github.com/slack-go/slack"
)

// AppInstallation receives redirects after successful OAuth flow whose result is the
// installation of the bot to a workspace. The request URL contains "code" parameter
// which is the authorization code in OAuth flow and can be exchanged with an access
// and a refresh token using the client's ID and secret.
//
// The resulting tokens belong to the entire team, not just a single user.
func (s *server) AppInstallation(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		writeJsonResponse(w, http.StatusBadRequest, newHttpError("missing 'code' url parameter for app installation", nil))
		return
	}

	// exchange authorization code with tokens
	resp, err := slack.GetOAuthV2Response(http.DefaultClient, s.cfg.SlackClientID, s.cfg.SlackClientSecret, code, "")
	if err != nil {
		writeJsonResponse(w, http.StatusInternalServerError, newHttpError("error exchanging authorization code with access token", detailedSlackError(err)))
		return
	}

	// save team information and tokens to database
	_, err = s.st.AddSlackTeam(resp.Team.ID, resp.AppID, resp.AccessToken, resp.RefreshToken)
	if err != nil {
		writeJsonResponse(w, http.StatusInternalServerError, newHttpError("error storing slack team information", err))
		return
	}

	// redirect to the "About" tab of the app
	http.Redirect(w, r, fmt.Sprintf("slack://app?team=%s&id=%s&tab=about", resp.Team.ID, resp.AppID), http.StatusFound)
}
