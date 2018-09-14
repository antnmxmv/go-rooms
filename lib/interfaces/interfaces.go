package interfaces

import "encoding/json"

type Game interface {
	// must return game name
	GetName() string

	Initialize()
	// should return bool for different specific goals and error != nil if parameters are unacceptable
	Action(int, json.RawMessage) (string, error)
}
