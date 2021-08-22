package util

import "runtime"

func GetDefaultXrayPath() string {
	defaultXrayPath := "./xray"
	if runtime.GOOS == "windows" {
		defaultXrayPath = defaultXrayPath + ".exe"
	}
	return defaultXrayPath
}
