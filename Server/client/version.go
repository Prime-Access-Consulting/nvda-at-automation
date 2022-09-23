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

func matchParts(compare string, provided []int, requested []int) bool {
	if compare == Equals {
		return match(provided, requested, matchEqual)
	}

	if compare == GreaterThan {
		return match(provided, requested, matchGreaterThan)
	}

	if compare == LesserThan {
		return match(provided, requested, matchLesserThan)
	}

	if compare == GreaterThanOrEquals {
		return match(provided, requested, matchGreaterThan) || match(provided, requested, matchEqual)
	}

	if compare == LesserThanOrEquals {
		return match(provided, requested, matchLesserThan) || match(provided, requested, matchEqual)
	}

	return false
}

type matcher func(int, int) bool

func match(provided []int, requested []int, f matcher) bool {
	for index, part := range requested {
		if !f(provided[index], part) {
			return false
		}
	}

	return true
}

func matchEqual(a int, b int) bool {
	return a == b
}

func matchGreaterThan(a int, b int) bool {
	return a > b
}

func matchLesserThan(a int, b int) bool {
	return a < b
}
