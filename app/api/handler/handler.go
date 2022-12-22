package handler

import (
	"github.com/go-pg/pg/v10"
)

type (
	Handler struct {
		DB *pg.DB
	}
)

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)
