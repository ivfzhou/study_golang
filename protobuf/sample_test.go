package protobuf_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"gitee.com/ivfzhou/study_golang/protobuf"
	importing "gitee.com/ivfzhou/study_golang/protobuf/import"
)

func TestSample(t *testing.T) {
	println("1:int32")
	d := protobuf.TestSample{Num: 10}
	bytes, err := proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("1:int32")
	d = protobuf.TestSample{Num: 300}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("2:sint32")
	d = protobuf.TestSample{Num2: -5}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("3:string")
	d = protobuf.TestSample{Str: "abccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc0000000000000000000000000000000"}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("4:repeated int32")
	d = protobuf.TestSample{Num3: []int32{1, 2}}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("4:repeated int32")
	d = protobuf.TestSample{Num3: []int32{}}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("5:enum")
	d = protobuf.TestSample{E: protobuf.Test_TEST_FIRST}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("6:map")
	d = protobuf.TestSample{M: map[int32]string{1: "a", 2: "b"}}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("7:nested_msg")
	d = protobuf.TestSample{Tm: &protobuf.TestSample_TestM{I: 1}}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("8:float")
	d = protobuf.TestSample{F: 7.3125}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("9:sfixed32")
	d = protobuf.TestSample{Fx: -5}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("10:bytes")
	d = protobuf.TestSample{Bs: []byte{255, 1}}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("11:repeated msg")
	d = protobuf.TestSample{Kvs: []*protobuf.KV{{Key: 1, Value: "a"}, {Key: 1, Value: "a"}}}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}

	println("1025: int32")
	d = protobuf.TestSample{Num4: 1}
	bytes, err = proto.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes {
		fmt.Printf("%8b\n", v)
	}
}

func TestZigZag(t *testing.T) {
	x := int32(-5)
	y := protobuf.Encode(x)
	println(y)
	r := protobuf.Decode(y)
	println(r)
}

func TestRepeatedMessage(t *testing.T) {
	bs := []byte{0b1011010, 0b1000, 0b1000, 0b1, 0b10010, 0b1, 0b1100010, 0b10010, 0b1, 0b1100001}
	kv := &protobuf.TestSample{}
	err := proto.Unmarshal(bs, kv)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(kv.Kvs)
}

func TestAny(t *testing.T) {
	a := &protobuf.TestAny{Any: &anypb.Any{
		TypeUrl: "a",
		Value:   []byte{1, 2},
	}}
	fmt.Println(a)
	bytes, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bytes))
}

func TestOneof(t *testing.T) {
	o := &protobuf.TestOneof{One: &protobuf.TestOneof_A{A: "a"}}
	bytes, err := proto.Marshal(o)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%b\n", bytes) // [11010 1 1100001]

	bs := []byte{0b1010, 0b1, 0b1100010, 0b11010, 0b1, 0b1100001}
	o = &protobuf.TestOneof{}
	err = proto.Unmarshal(bs, o)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(o.One)
}

func TestDesc(t *testing.T) {
	t.Log(protobuf.File_protobuf_sample_proto.Enums())
	t.Log(protobuf.File_protobuf_sample_proto.SourceLocations())
}

func TestText(t *testing.T) {
	kv := protobuf.KV{
		Key:   1,
		Value: "a",
	}
	bytes, err := prototext.Marshal(&kv)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bytes))

	kv = protobuf.KV{}
	err = prototext.Unmarshal(bytes, &kv)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(&kv)
}

func TestReflect(t *testing.T) {
	descriptor := protobuf.AllowAlias.Descriptor(protobuf.AllowAlias_ALLOW_ALIAS_UNSPECIFIED)
	t.Log(descriptor.Values().Get(0).Options().ProtoReflect().Get(protobuf.E_MyEnumValueOption.TypeDescriptor()).Uint())
	t.Log(descriptor.Values().Get(0).Options().ProtoReflect().Get(importing.E_MyEnumValueOption.TypeDescriptor()).Uint())
	t.Log((*protobuf.Req).ProtoReflect(&protobuf.Req{}).Descriptor().Options())
}

func TestOneof1(t *testing.T) {
	/*o := &protobuf.TestOneof{One: s{}}
	t.Log(o)*/
}
