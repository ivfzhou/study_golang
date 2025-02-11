package c

import "gitee.com/ivfzhou/study_golang/language_specification/mvs/d"

func Version() string {
	_ = d.Version()
	return "v1.3.0"
}
