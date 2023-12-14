package tg_bot

import (
	"fmt"
	"regexp"
)

func ReadMsg(msg string) (string, error) {
	re := regexp.MustCompile(`front:\s*(\S+)\s*back:\s*(\S+)\.`)

	matches := re.FindStringSubmatch(msg)

	if len(matches) != 3 {
		return msg, fmt.Errorf("Incorrect format")
	}

	return msg, nil
}
