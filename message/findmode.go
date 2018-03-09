package message

import "strings"

type FindMode int32

const (
	FindIndividual FindMode = iota // 0
	FindAll
	FindAllDot
)

func determineFindMode(target string) FindMode {
	if target == "all" {
		return FindAll
	}
	if strings.HasPrefix(target, "all.") {
		return FindAllDot
	}
	return FindIndividual
}

func parseFindToken(token string) (findMode FindMode, dotPart string, target string) {
	findMode = determineFindMode(token)
	switch findMode {
	case FindAll:
		dotPart = ""
		target = token
		break

	case FindAllDot:
		dotPart = "all"
		target = strings.TrimPrefix(token, "all.")
		break

	case FindIndividual:
		firstDot := strings.Index(token, ".")
		if firstDot < 0 {
			// no dot (or first char is dot, same thing
			dotPart = ""
			target = token
		} else {
			dotPart = token[0:firstDot]
			target = token[firstDot+1:]
		}
		break
	}
	return
}
