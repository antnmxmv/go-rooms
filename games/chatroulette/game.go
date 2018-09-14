package chatroulette

import (
	"encoding/json"
	"errors"
)

var name = "chatroulette"

type Chat struct{}

func (c Chat) GetName() string {
	return name
}

func (c Chat) Initialize() {

}

func (c Chat) Action(role int, raw json.RawMessage) (string, error) {
	var message = string(raw)
	if message == "\"\"" || message == "\"\\n\"" {
		return "", errors.New("empty message")
	}
	var response = struct {
		Role    int
		Message json.RawMessage
	}{role, raw}
	res, _ := json.Marshal(response)
	return string(res), nil
}
