package api

type (
	FlagEvent struct {
		Name  string `json:"name"`
		Kind  string `json:"kind"`
		Value any    `json:"value"`
	}
)
