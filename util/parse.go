package util

import (
	"net/mail"
	"regexp"
	"strconv"
	"time"
)

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidDate(d string) bool {
	if length := len(d); length != 10 {
		return false
	}
	if matched, err := regexp.MatchString("^[a-zA-Z]+$", d); err != nil || matched {
		return false
	}
	dd, err := strconv.Atoi(d[0:2])
	if err != nil {
		return false
	}
	if !InBetween(dd, 1, 31) {
		return false
	}
	mm, err := strconv.Atoi(d[3:5])
	if err != nil {
		return false
	}
	if !InBetween(mm, 1, 12) {
		return false
	}
	yy, err := strconv.Atoi(d[6:])
	if err != nil {
		return false
	}
	current_year := time.Now().Year()
	if !InBetween(yy, 1920, current_year+5) {
		return false
	}
	if yy > current_year-18 {
		return false
	}
	return true
}

func InBetween(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	} else {
		return false
	}
}
