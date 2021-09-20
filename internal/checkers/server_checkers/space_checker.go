package server_checkers

import (
	"encoding/json"
	"fmt"
)

type (
	spaceAddress struct {
		Address string `json:"address"`
		Total   string `json:"total"`
		Free    string `json:"free"`
		Status  bool   `json:"status"`
	}

	SpaceChecker struct {
		spaces []*spaceAddress
	}
)

func NewSpaceChecker(spaceRaw json.RawMessage) *SpaceChecker {
	fmt.Printf("%v\n", string(spaceRaw))

	var spaces []*spaceAddress
	err := json.Unmarshal(spaceRaw, &spaces)
	if err != nil {
		println(err.Error())
	}

	return &SpaceChecker{
		spaces: spaces,
	}
}

func (c *SpaceChecker) Check() bool {
	for _, address := range c.spaces {
		println(address.Free, address.Status)
	}

	return true
}
