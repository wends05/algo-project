package structs

type Edge struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Weight int `json:"weight"`
}
