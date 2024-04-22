package misc

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
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

func CurrentTimeFull() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}
func CurrentTime() string {
	return time.Now().Format("2006-01-02")
}

func DateParseToFull(s string) string {
	t, err := time.Parse("2006-01-02", s[:10])
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05.000")
}
