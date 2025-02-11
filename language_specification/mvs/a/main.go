package a

import "gitee.com/ivfzhou/study_golang/language_specification/mvs/c"

func Version() string {
	_ = c.Version()
	return "v1.1.0"
}
