package strings_utils

import (
	"github.com/iancoleman/strcase"
)

func ToSnake(s string) string {
	return strcase.ToSnake(s)
}
