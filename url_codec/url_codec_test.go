package url_codec_test

import (
	"testing"

	urlcodec "gitee.com/ivfzhou/study_golang/url_codec"
)

func TestEncode(t *testing.T) {
	// 11100100-10111101-10100000 11100101-10100101-10111101
	s := "你好世界 hello world"
	shouldEncodeS := "%E4%BD%A0%E5%A5%BD%E4%B8%96%E7%95%8C%20hello%20world"
	encodedS := urlcodec.Encode(s)
	t.Log(encodedS)
	if encodedS != shouldEncodeS {
		t.Errorf("%s != %s", encodedS, shouldEncodeS)
	}
}

func TestDecode(t *testing.T) {
	s := "%E4%BD%A0%E5%A5%BD%E4%B8%96%E7%95%8C%20hello%20world"
	shouldDecodedS := "你好世界 hello world"
	decodedS := urlcodec.Decode(s)
	t.Log(decodedS)
	if decodedS != shouldDecodedS {
		t.Errorf("%s != %s", decodedS, shouldDecodedS)
	}
}

func TestEscapeNonASCII(t *testing.T) {
	s := "你好世界 hello world"
	shouldEncoedS := "%e4%bd%a0%e5%a5%bd%e4%b8%96%e7%95%8c hello world"
	encodedS := urlcodec.EscapeNonASCII(s)
	t.Log(encodedS)
	if encodedS != shouldEncoedS {
		t.Errorf("%s != %s", encodedS, shouldEncoedS)
	}
}

func TestUnescapeNonASCII(t *testing.T) {
	s := "%E4%BD%A0%E5%A5%BD%E4%B8%96%E7%95%8C%20hello%20world"
	shouldDecodedS := "你好世界 hello world"
	decodedS := urlcodec.UnescapeNonASCII(s)
	t.Log(decodedS)
	if decodedS != shouldDecodedS {
		t.Errorf("%s != %s", decodedS, shouldDecodedS)
	}
}
