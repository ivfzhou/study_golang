package types_test

type Struct struct{}

type Struct1 struct {
	_    Struct
	_    Struct
	x, y Struct
	Struct
}

func (Struct1) _() {}
