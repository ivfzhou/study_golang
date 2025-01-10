package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	// 1.读取标准输入，接收 proto 解析的文件内容，并解析成结构体。
	input, _ := io.ReadAll(os.Stdin)
	var req pluginpb.CodeGeneratorRequest
	err := proto.Unmarshal(input, &req)
	if err != nil {
		panic(err)
	}

	// 2.生成插件。
	opts := protogen.Options{}
	plugin, err := opts.New(&req)
	if err != nil {
		panic(err)
	}

	// 3.在插件 plugin.Files 就是 *.proto 的内容了，是一个切片，每个切片元素代表一个文件内容，我们只需要遍历这个文件就能获取到文件的信息了。
	for _, v := range plugin.Files {
		// 创建一个 buf 写入生成的文件内容
		var buf bytes.Buffer
		// 写入 go 文件的 package 名
		pkg := fmt.Sprintf("package %s", v.GoPackageName)
		buf.Write([]byte(pkg))
		context := ""
		// 遍历消息，这个内容就是 protobuf 的每个消息
		for _, msg := range v.Messages {
			// 遍历消息的每个字段
			for _, field := range msg.Fields {
				op, ok := field.Desc.Options().(*descriptorpb.FieldOptions)
				if ok {
					value := GetJsonTag(op)
					context += fmt.Sprintf("%v\n", value)
				}
			}
			buf.Write([]byte(fmt.Sprintf(`
           func (x *%s) optionsTest() {
            %s
           }`, msg.GoIdent.GoName, context)))
		}
		// 指定输入文件名，输出文件名为 demo.foo.go
		filename := v.GeneratedFilenamePrefix + ".foo.go"
		file := plugin.NewGeneratedFile(filename, ".")
		// 将内容写入插件文件内容
		file.Write(buf.Bytes())
	}

	// 生成响应
	stdout := plugin.Response()
	out, err := proto.Marshal(stdout)
	if err != nil {
		panic(err)
	}

	// 将响应写回标准输入, protoc 会读取这个内容
	fmt.Fprintf(os.Stdout, string(out))
}

func GetJsonTag(field *descriptorpb.FieldOptions) interface{} {
	if field == nil {
		return ""
	}
	v := proto.GetExtension(field, E_Tag)
	return v.(string)
}
