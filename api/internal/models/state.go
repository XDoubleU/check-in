package models

type State struct {
	IsMaintenance    bool `json:"isMaintenance"`
	IsDatabaseActive bool `json:"isDatabaseActive"`
} //	@name	State

type StateKey string

const (
	IsMaintenanceKey StateKey = "IsMaintenance"
)
