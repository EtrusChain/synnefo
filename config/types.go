package config

import "time"

type Duration struct {
	time.Duration
}

type OptionalDuration struct {
	value *time.Duration
}

type OptionalString struct {
	value *string
}
