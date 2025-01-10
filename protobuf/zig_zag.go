package protobuf

func Encode(x int32) int32 {
	// -1 11111111 => 00000001 => 1
	// 1 00000001 => 00000010 => 2
	return x<<1 ^ x>>31
}

func Decode(x int32) int32 {
	if x%2 == 0 {
		// 2 00000010 => 00000001 => 1
		return (x ^ 0) >> 1
	} else {
		// 1 => 00000001 => 11111111 => -1
		return (x ^ -1) >> 1
	}
}
