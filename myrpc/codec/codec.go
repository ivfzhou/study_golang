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

import "sync"

var (
	serializerMap sync.Map
	compressorMap sync.Map
)

type Serializer interface {
	Serialize(interface{}) ([]byte, error)
	Deserialize([]byte, interface{}) error
}

type Compressor interface {
	Compress([]byte) ([]byte, error)
	Decompress([]byte) ([]byte, error)
}

type (
	SerializerType string
	CompressorType string
)

func GetSerializer(name SerializerType) Serializer {
	value, ok := serializerMap.Load(name)
	if ok && value != nil {
		return value.(Serializer)
	}
	return nil
}

func GetCompressor(name CompressorType) Compressor {
	value, ok := compressorMap.Load(name)
	if ok && value != nil {
		return value.(Compressor)
	}
	return nil
}

func RegisterSerializer(name SerializerType, serializer Serializer) {
	serializerMap.Store(name, serializer)
}

func RegisterCompressor(name CompressorType, compressor Compressor) {
	compressorMap.Store(name, compressor)
}

func DeregisterSerializer(name SerializerType) {
	serializerMap.Delete(name)
}

func DeregisterCompressor(name CompressorType) {
	compressorMap.Delete(name)
}
