package statements_test

import "testing"

func TestSwitch(t *testing.T) {
	switch "1" {
	default:
		fallthrough
	case "":
		println(1)
	}
}

func TestSwitch1(t *testing.T) {
	type (
		i1 interface {
			M()
		}
		i2 interface {
			M()
			M1()
		}
		i3 interface {
			M2()
		}
	)
	var i interface{}
	switch vi := 0; v := i.(type) {
	case nil:
		println("yes is nil.")
		// fallthrough
	case i1, i2:
		_, _ = v, vi
	case i3:
		v.M2()
	}
}
