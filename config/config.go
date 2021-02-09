package config

import "os"

// IsLocal はlocal環境か判定する
func IsLocal() bool {
	return os.Getenv("ENV") == "local"
}
