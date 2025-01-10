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
	"bytes"
	"compress/gzip"
	"io"
	"sync"
)

const GzipCompressorType CompressorType = "gzip"

type gzipCompressor struct {
	readers sync.Pool
	writers sync.Pool
}

func init() {
	RegisterCompressor(GzipCompressorType, &gzipCompressor{})
}

func (c *gzipCompressor) Compress(data []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer, ok := c.writers.Get().(*gzip.Writer)
	if ok {
		writer.Reset(buf)
	} else {
		writer = gzip.NewWriter(buf)
	}
	defer c.writers.Put(writer)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *gzipCompressor) Decompress(data []byte) ([]byte, error) {
	buf := bytes.NewReader(data)
	reader, ok := c.readers.Get().(*gzip.Reader)
	var err error
	if ok {
		err = reader.Reset(buf)
	} else {
		reader, err = gzip.NewReader(buf)
	}
	if err != nil {
		return nil, err
	}
	defer c.readers.Put(reader)
	return io.ReadAll(reader)
}
