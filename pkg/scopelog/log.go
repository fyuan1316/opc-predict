package scopelog

import (
	"fmt"
	"log"
)

func Printf(scope, format string, p ...interface{}) {
	newFormat := fmt.Sprintf("[%s]: %s", scope, format)
	log.Printf(newFormat, p...)
}
