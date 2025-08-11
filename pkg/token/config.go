package token

import "time"

var (
	accessExpTime  = time.Minute * 10 // its mean 10 minute
	refreshExpTime = time.Minute * 30 // its mean 30 minute
)
