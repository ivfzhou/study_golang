package b

import "gitee.com/ivfzhou/study_golang/language_specification/mvs/e"

func Version() string {
	_ = e.Version()
	return "v1.3.0"
}
