package config

import "os"

// IsLocal はlocal環境か判定する
func IsLocal() bool {
	return os.Getenv("ENV") == "local"
}

// CORSAllowOrigin はCORSのAllow Originを取得する
func CORSAllowOrigin() string {
	return os.Getenv("CORS_ALLOW_ORIGIN")
}
