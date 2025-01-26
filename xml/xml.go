package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
)

type Server struct {
	XMLName     xml.Name  `xml:"server"`    // 容纳元素名及命名空间。
	Type        string    `xml:"type,attr"` // 容纳元素属性。
	AnyAttr     string    `xml:",any,attr"`
	Name        Name      `xml:"name"`     // 容纳元素内部元素的名称。
	Comment     string    `xml:",comment"` // 容纳注释内容。
	Services    []Service `xml:"service"`
	Password    string    `xml:"secret>password"` // 容纳指定元素。
	CharData    string    `xml:",cdata"`          // 容纳不解析的 xml 文本。
	Content     string    `xml:",innerxml"`       // 所有元素内的内容。
	NonOut      string    `xml:"-"`               // 忽略该字段。
	IgnoreEmpty string    `xml:",omitempty"`      // 转 xml 文本时字段为零值忽略输出。
	Any         string    `xml:",any"`            // 接收任何内容。
}

func (s *Server) String() string {
	bytes, _ := json.MarshalIndent(s, "", "\t")
	return string(bytes)
}

type Name struct {
	XMLName xml.Name `xml:"name"`
	Info    string   `xml:"info,attr"`
	Content string   `xml:",innerxml"`
}

type Service struct {
	Name     string `xml:"name,attr"`
	IP       string `xml:"ip"`
	Port     uint16 `xml:"port"`
	Location string `xml:"location"`
}

func TestUnmarshal() {
	data, err := os.ReadFile("./server.xml")
	if err != nil {
		log.Fatal(err)
	}

	server := new(Server)
	err = xml.Unmarshal(data, server)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", server)
}

func TestMarshal() {
	server := &Server{
		XMLName: xml.Name{
			Space: "",
			Local: "server",
		},
		Type:    "web",
		Comment: " Provide shopping, email and other services. ",
		Name: Name{
			XMLName: xml.Name{
				Space: "",
				Local: "name",
			},
			Info:    "info",
			Content: "web",
		},
		Services: []Service{
			{
				Name:     "shopping",
				IP:       "127.0.0.1",
				Port:     80,
				Location: "beijing",
			},
			{
				Name:     "email",
				IP:       "::1",
				Port:     443,
				Location: "shanghai",
			},
		},
		Password:    "123456",
		CharData:    "this is character data.",
		NonOut:      "nothings",
		IgnoreEmpty: "",
		Any:         "any info",
	}

	bytes, err := xml.MarshalIndent(server, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(bytes))
}
