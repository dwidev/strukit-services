package helper

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func ParseToIDR(value float64) string {
	f := message.NewPrinter(language.Indonesian)
	return f.Sprintf("Rp. %d", int(value))
}
