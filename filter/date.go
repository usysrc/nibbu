package filter

import (
	"fmt"
	"time"
)

// human readable dates from the db saved date
func Date(input string) (string, error) {
	// Parse the input date in RFC 3339 format
	t, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return "", err
	}

	// Get the day and determine the ordinal suffix
	day := t.Day()
	var suffix string
	switch day % 10 {
	case 1:
		if day != 11 {
			suffix = "st"
		} else {
			suffix = "th"
		}
	case 2:
		if day != 12 {
			suffix = "nd"
		} else {
			suffix = "th"
		}
	case 3:
		if day != 13 {
			suffix = "rd"
		} else {
			suffix = "th"
		}
	default:
		suffix = "th"
	}

	// Format the date into the desired format
	formattedDate := fmt.Sprintf("%d%s of %s, %d", day, suffix, t.Month(), t.Year())
	return formattedDate, nil
}
