package main

import "errors"

// ErrNoTarget is returned when the target section is missing
var (
	ErrNoTarget = errors.New("no 'target' section")
	ErrNoJob    = errors.New("no 'job' field specified")
	ErrNoURLs   = errors.New("no URLs specified")
	ErrKV       = errors.New("key/value missing")
)
