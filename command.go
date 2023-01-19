package main

import "time"

type Command string

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
	CMDDel Command = "DEL"
)

type Message struct {
	Command Command
	Key     []byte
	Value   []byte
	TTL     time.Duration
}

// func pa
