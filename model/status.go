package model

type LimitStatus int32 // 1-正常 2-限制

const (
	LimitNormal LimitStatus = iota + 1
	LimitLocked
)
