package models

type School struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ReadOnly bool   `json:"readOnly"`
} //	@name	School
