/*
 * Copyright (c) 2023 ivfzhou
 * myrpc is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *          http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

package codec

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const JsonpbSerializerType SerializerType = "jsonpb"

type jsonpbSerializer struct {
	protojson.MarshalOptions
	protojson.UnmarshalOptions
}

func init() {
	RegisterSerializer(JsonpbSerializerType, &jsonpbSerializer{
		protojson.MarshalOptions{AllowPartial: true, UseEnumNumbers: true, UseProtoNames: true},
		protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true},
	})
}

func (s *jsonpbSerializer) Serialize(v interface{}) ([]byte, error) {
	return s.Marshal(v.(proto.Message))
}

func (s *jsonpbSerializer) Deserialize(data []byte, v interface{}) error {
	return s.Unmarshal(data, v.(proto.Message))
}
