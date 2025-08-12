package openshock

type ShockType string

const (
	Shock   ShockType = "Shock"
	Vibrate           = "Vibrate"
)

type ShockControl struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Intensity int    `json:"intensity"`
	Duration  int    `json:"duration"`
	Exclusive bool   `json:"exclusive"`
}

type ShockerControlMessage struct {
	Shocks     []ShockControl `json:"shocks"`
	CustomName string         `json:"customName"`
}

type ShockerEntry struct {
	Name      string `json:"name"`
	IsPaused  bool   `json:"isPaused"`
	CreatedOn string `json:"createdOn"`
	Id        string `json:"id"`
	RfId      int    `json:"rfId"`
	Model     string `json:"model"`
}

type LoadShockersMessage struct {
	Message string `json:"message"`
	Data    []struct {
		Shockers  []ShockerEntry `json:"shockers"`
		Id        string         `json:"id"`
		Name      string         `json:"name"`
		CreatedOn string         `json:"createdOn"`
	} `json:"data"`
}
