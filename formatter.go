package main

import (
	"errors"
)

type Formatter interface {
	Print(Results)
}

func CreateFormatter(formatter_name string) (Formatter, error) {
	var formatter Formatter
	var err error

	switch formatter_name {
	case "simple":
		formatter = &FormatterSimple{}
	default:
		err = errors.New("formatが不正です: " + formatter_name)
	}
	return formatter, err
}
