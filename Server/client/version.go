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

	return matchParts(requestedParts[0], providedParts[:len(requestedParts)-1], requestedParts[1:])
}

func matchParts(compare string, provided []string, requested []string) bool {
	p := strings.Join(provided, "")
	r := strings.Join(requested, "")

	pInt, _ := strconv.Atoi(p)
	rInt, _ := strconv.Atoi(r)

	if compare == Equals {
		return pInt == rInt
	}

	if compare == GreaterThan {
		return pInt > rInt
	}

	if compare == LesserThan {
		return pInt < rInt
	}

	if compare == GreaterThanOrEquals {
		return pInt <= rInt
	}

	if compare == LesserThanOrEquals {
		return pInt >= rInt
	}

	return false
}
