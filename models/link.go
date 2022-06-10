package models

import (
	"time"

	"github.com/dhruv42/hraasaka/enums"
)

type Link struct {
	Url       string
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    enums.ShortLinkStatus
}
