package models

import "time"

type Scope int

const (
	AccessScope  Scope = 0
	RefreshScope Scope = 1
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    string
	Expiry    time.Time
	Scope     Scope
	Used      bool
}
