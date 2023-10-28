package api

import "time"

type Metadata struct {
	CreatedAt time.Time  `json:"createdAt,omitempty" yaml:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" yaml:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" yaml:"deletedAt,omitempty"`
	StolenAt  *time.Time `json:"stolenAt,omitempty" yaml:"stolenAt,omitempty"`
	FatalAt   *time.Time `json:"fatalAt,omitempty" yaml:"fatalAt,omitempty"`
}

type DispenserEditable struct {
	Ref  string `json:"ref"`
	Beer string `json:"beer" binding:"required"`
	Size string `json:"size" binding:"required,oneof=S M L"`
}

type Dispenser struct {
	DispenserEditable `yaml:",inline" bson:",inline"`

	Metadata *Metadata `json:"metadata" yaml:"metadata" `

	State DispenserState `json:"state" yaml:"state"`

	// State for
	Status *DispenserStatus `json:"status" yaml:"status"`
}

type DispenserStatus struct {
	Status InternalStatus `json:"status,omitempty" yaml:"status,omitempty"`
	Reason string         `json:"reason,omitempty" yaml:"reason,omitempty"`
}

type DispenserSize string

const (
	DispenserSizeS DispenserSize = "S"
	DispenserSizeM DispenserSize = "M"
	DispenserSizeL DispenserSize = "L"
)

type DispenserState string

const (
	DispenserNone DispenserState = "NONE"

	DispenserReady      DispenserState = "READY"
	DispenserRefreshing DispenserState = "REFRESHING"
	DispenserEmpty      DispenserState = "EMPTY"
)

type InternalStatus string

const (
	InternalStatusPending  InternalStatus = "PENDING"
	InternalStatusArchived InternalStatus = "ARCHIVED"
	InternalStatusActive   InternalStatus = "ACTIVE"
)

type DispenserQuery struct {
	Statuses []InternalStatus
}
