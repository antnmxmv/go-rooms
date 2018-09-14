package tic_tac_toe

import "encoding/json"

type response struct {
	Winner int
	Turn   int
	Grid   [][]int
}

func (r response) Marshal() string {
	msg, _ := json.Marshal(r)
	return string(msg)
}
