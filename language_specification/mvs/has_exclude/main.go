package has_exclude

import _ "google.golang.org/protobuf/proto"

func Version() string {
	println("有 exclude replace 的项目无法被 go install")
	return "v1.0.0"
}
