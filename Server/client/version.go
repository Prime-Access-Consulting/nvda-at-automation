package client

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	Equals              = "=="
	GreaterThan         = ">"
	LesserThan          = "<"
	GreaterThanOrEquals = ">="
	LesserThanOrEquals  = "<="
)

func VersionRequestMatches(provided string, requested string) bool {
	r, _ := regexp.Compile("^(==|<=|>=|>|<) (\\d{4})(?:\\.([1-9]{1,2})(?:\\.([1-9]{1,2}))?)?$")

	if !r.MatchString(requested) {
		return false
	}

	providedParts := strings.Split(provided, ".")
	requestedParts := r.FindStringSubmatch(requested)[1:]

	return matchParts(requestedParts[0], asInts(providedParts), asInts(requestedParts[1:]))
}

func asInts(parts []string) []int {
	var i []int

	for _, s := range parts {
		if len(s) == 0 {
			continue
		}
		str, _ := strconv.Atoi(s)
		i = append(i, str)
	}

	return i
}

func matchParts(compare string, providedParts []int, requestedParts []int) bool {
	if compare == Equals {
		return matchEqual(providedParts, requestedParts)
	}

	if compare == GreaterThan {
		return matchGreaterThan(providedParts, requestedParts)
	}

	if compare == LesserThan {
		return matchLesserThan(providedParts, requestedParts)
	}

	if compare == GreaterThanOrEquals {
		return matchGreaterThan(providedParts, requestedParts) || matchEqual(providedParts, requestedParts)
	}

	if compare == LesserThanOrEquals {
		return matchLesserThan(providedParts, requestedParts) || matchEqual(providedParts, requestedParts)
	}

	return false
}

func matchEqual(provided []int, requested []int) bool {
	for index, part := range requested {
		if provided[index] != part {
			return false
		}
	}

	return true
}

func matchGreaterThan(provided []int, requested []int) bool {
	for index, part := range requested {
		if provided[index] <= part {
			return false
		}
	}

	return true
}

func matchLesserThan(provided []int, requested []int) bool {
	for index, part := range requested {
		if provided[index] >= part {
			return false
		}
	}

	return true
}
