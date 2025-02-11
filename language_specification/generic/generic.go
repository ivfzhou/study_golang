package generic

import "fmt"

type ArrayList[E any] struct {
	val []E
}

func (l *ArrayList[E]) Add(e E) {
	l.val = append(l.val, e)
}

func (l *ArrayList[E]) Get(index int) E {
	return l.val[index]
}

func (l *ArrayList[E]) Size() int {
	return len(l.val)
}

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 | complex64 | complex128
}

func Sum[K comparable, V Number](m map[K]V) V {
	var sum V
	for _, v := range m {
		sum += v
	}
	return sum
}

func ToString[E fmt.Stringer](v E) string {
	return v.String()
}

type Car[E any] interface {
	Name() E
}

type BaoMa[E any] struct {
	name E
}

func (c *BaoMa[E]) Name() E {
	return c.name
}

func (c *BaoMa[E]) SetName(name E) {
	c.name = name
}

type BYD struct {
	name string
}

func (c *BYD) Name() string {
	return c.name
}

func (c *BYD) SetName(name string) {
	c.name = name
}

// ~表示包含底层类型为 []int 的类型
func AFunc[T ~[]int](x T) {
	fmt.Println(x)
}
