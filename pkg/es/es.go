package es

import (
	"net/http"

	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

func New(l *logger.Logger, c *config.Config) *ES {
	client := &http.Client{}

	e := &ES{
		log:    l,
		cfg:    c,
		client: client,
	}

	return e
}
