package token

import "time"

var (
	accessExpTime  = time.Hour * 10 // its mean 10 Hour
	refreshExpTime = time.Hour * 30 // its mean 30 Hour
)
