package util

import "strings"

func SplitName(name string) (string, string) {
	// Find the first space in the name.
	spaceIndex := strings.Index(name, " ")

	// If there is no space, the full name is the first name.
	if spaceIndex == -1 {
		return name, ""
	}

	// The last name is the part of the name after the first space.
	lastName := name[spaceIndex:]

	// The first name is the part of the name before the first space.
	firstName := name[:spaceIndex]

	// Return the first and last names.
	return strings.TrimSpace(firstName), strings.TrimSpace(lastName)
}
