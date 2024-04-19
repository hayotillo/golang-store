package misc

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func NewUUID() string {
	return strings.ToUpper(fmt.Sprintf("%s", uuid.New()))
}

func SQLPlaceHolder(count int, start int) string {
	res := ""
	for i := start; i < (start + count); i++ {
		res = fmt.Sprintf("%s$%d, ", res, i)
	}
	return strings.TrimSuffix(res, ", ")
}

func SqlConstraintErrorStr(constraint string) string {
	return fmt.Sprintf("pq: duplicate key value violates unique constraint \"%s\"", constraint)
}
