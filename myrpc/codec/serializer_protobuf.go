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

import "google.golang.org/protobuf/proto"

const ProtobufSerializerType SerializerType = "protobuf"

type protobufSerializer struct{}

func init() {
	RegisterSerializer(ProtobufSerializerType, &protobufSerializer{})
}

func (*protobufSerializer) Serialize(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (*protobufSerializer) Deserialize(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
