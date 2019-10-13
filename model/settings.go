package model

import "time"

type Settings struct {
	Server struct {
		SessionExpires time.Duration
	}
}
