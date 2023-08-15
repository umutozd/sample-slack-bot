package storage

type SlackTeam struct {
	ID           string `json:"id,omitempty"`
	AppID        string `json:"app_id,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
