package model

import "errors"

var (
	// ErrPostAlreadyExisted はpostが既に存在しているerrorを生成する
	ErrPostAlreadyExisted = errors.New("post has already existed")
)
