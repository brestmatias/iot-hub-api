package model

type BeaconResponse struct {
	ID      string   `json:id`
	Outputs []string `json:outputs`
}
