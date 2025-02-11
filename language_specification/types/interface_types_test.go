package types_test

type I0 interface {
	I1
	/* M() error // 相同名字的方法需要签名也相同。冲突 I1.M。*/

	/*I2 // 不能直接或间接内嵌自己。*/
}

type I1 interface {
	M()
}

type I2 interface {
	I0
}
