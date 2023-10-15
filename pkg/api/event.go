package api

import "time"

type Record struct {
	Id      string    `json:"id" yaml:"id"`
	Date    time.Time `json:"date" yaml:"date" `
	Message string    `json:"message" yaml:"message"`

	Target *Dispenser `json:"dispenser,omitempty" yaml:"dispenser,omitempty"`
}
