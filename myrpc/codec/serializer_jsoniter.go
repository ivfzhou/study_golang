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

import jsoniter "github.com/json-iterator/go"

const JsoniterSerializerType SerializerType = "jsoniter"

type jsoniterSerializer struct {
	jsoniter.API
}

func init() {
	RegisterSerializer(JsoniterSerializerType, &jsoniterSerializer{jsoniter.ConfigCompatibleWithStandardLibrary})
}

func (s *jsoniterSerializer) Serialize(v interface{}) ([]byte, error) {
	return s.Marshal(v)
}

func (s *jsoniterSerializer) Deserialize(data []byte, v interface{}) error {
	return s.Unmarshal(data, v)
}
