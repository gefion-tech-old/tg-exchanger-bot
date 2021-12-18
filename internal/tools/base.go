package tools

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func Recovery(logger *logrus.Logger) {
	if err := recover(); err != nil {
		go fmt.Printf("Panic: %s\n", err.(string))
		go logger.Error(err)
	}
}
