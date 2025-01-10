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

const NoopSerializerType SerializerType = "noop"

type noopSerializer struct{}

func init() {
	RegisterSerializer(NoopSerializerType, &noopSerializer{})
}

func (s *noopSerializer) Serialize(v interface{}) ([]byte, error) {
	return v.([]byte), nil
}
func (s *noopSerializer) Deserialize(data []byte, v interface{}) error {
	copy(v.([]byte), data)
	return nil
}
