package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func CheckErr(err error) {
	if err == nil {
		return
	}
	_, _ = fmt.Fprintln(os.Stderr, "Error: ", err)
	logrus.Errorf("Error: %+v", err)
	os.Exit(1)
}

func Closeq(c io.Closer) {
	silently(c.Close())
}

func silently(_ ...interface{}) {}
