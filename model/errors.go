package model

import "errors"

var (
	// ErrPostNotFound はpostが存在しないerrorを生成する
	ErrPostNotFound = errors.New("post not found")
	// ErrPostAlreadyExisted はpostが既に存在しているerrorを生成する
	ErrPostAlreadyExisted = errors.New("post has already existed")

	// ErrTagNotFound はタグが存在しないエラーを表します。
	ErrTagNotFound = errors.New("tag not found")
	// ErrTagAlreadyExisted はpostが既に存在しているerrorを生成する
	ErrTagAlreadyExisted = errors.New("tag has already existed")
)
