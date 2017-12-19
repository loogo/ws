package server

type Message struct {
	SendTo []string `json:"sendto"`
	Type   int `type` //0: personal, 1:multiple, 2: all
	Data   string   `json:"data"`
}
