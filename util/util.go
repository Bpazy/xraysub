package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func CheckErr(err error) {
	if err == nil {
		return
	}
	_, _ = fmt.Fprintln(os.Stderr, "Error: ", err)
	log.Errorf("Error: %+v", err)
	os.Exit(1)
}

func Closeq(c io.Closer) {
	silently(c.Close())
}

func silently(_ ...interface{}) {
	// empty
}
