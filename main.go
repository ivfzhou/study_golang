package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/net/websocket"
)

func main() {
	to := 15
	switch to {
	case 1:
		x := byte(1)
		bitPrint(x)
		y := float32(7.3125)
		bitPrint(y)
		s := struct {
			x int32
			y float32
		}{x: int32(x), y: y}
		bitPrint(s)
	case 2:
		garbage()
	case 3:
		garbageNoFree()
	case 4:
		pointerSize()
	case 5:
		var x, y uint8 = 15, 17
		fmt.Printf("isMultiplicationOverflow(%d, %d) is %v, and x*y = %d\n",
			x, y, isMultiplicationOverflow(x, y), x*y)
	case 6:
		x := uint8(3)
		fmt.Printf("isPowerOfTwo(%d) is %v\n", x, isPowerOfTwo(x))
	case 7:
		x, y := 6, 9
		gcd, lcm := greatestCommonDivisor(x, y)
		fmt.Printf("greatestCommonDivisor(%d, %d) is %d, %d\n", x, y, gcd, lcm)
	case 8:
		primeSieve()
	case 9:
		fmt.Println(getMachineIPv4())
	case 10:
		createPng()
	case 11:
		regexpNonCapture()
	case 12:
		printFileBits("file_bit.txt", "file_bit_out.txt")
	case 13:
		shortestPath()
	case 14:
		websocketServer()
	case 15:
		arr := bytes.Split([]byte("1="), []byte("="))
		fmt.Println(len(arr))
		fmt.Println(string(arr[0]))
		fmt.Println(string(arr[1]))
	}
}

func websocketServer() {
	var req string
	http.Handle("/", websocket.Handler(func(conn *websocket.Conn) {
		for {
			if err := websocket.Message.Receive(conn, &req); err != nil {
				log.Fatal(err)
			} else {
				log.Println("req: ", req)
			}
			if err := websocket.Message.Send(conn, "reto: "+req); err != nil {
				log.Fatal(err)
			}
		}
	}))
	http.ListenAndServe(":12345", nil)
}

func shortestPath() {
	matrix := [3][5]int{}
	rand.Seed(time.Now().UnixMilli())
	for i := range matrix {
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(10)
		}
	}
	for i := range matrix {
		log.Println(matrix[i])
	}

	states := [3][5]int{}
	sum := 0
	for i := range matrix[0] {
		sum += matrix[0][i]
		states[0][i] = sum
	}
	sum = 0
	for i := range matrix {
		sum += matrix[i][0]
		states[i][0] = sum
	}
	for i := 1; i < len(matrix); i++ {
		for j := 1; j < len(matrix[i]); j++ {
			states[i][j] = matrix[i][j] + min(states[i-1][j], states[i][j-1])
		}
	}

	for i := range states {
		log.Println(states[i])
	}
}

func printFileBits(file, out string) {
	inFileObj, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer inFileObj.Close()

	inFileState, err := inFileObj.Stat()
	if err != nil {
		panic(err)
	}
	inFileSize := inFileState.Size()

	if err = os.MkdirAll(filepath.Dir(out), 0755); err != nil {
		panic(err)
	}
	outFileObj, err := os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer outFileObj.Close()

	const bytesPreLine = 16
	readTmp := make([]byte, 1)
	writeTmp := make([]byte, bytesPreLine)
	for i, j := int64(0), 0; i < inFileSize; i++ {
		n, err := inFileObj.Read(readTmp)
		if err != nil {
			panic(err)
		}
		if n != len(readTmp) {
			panic("read fail")
		}

		str := fmt.Sprintf("%02X", readTmp[0])
		n, err = outFileObj.WriteString(str)
		if err != nil {
			panic(err)
		}
		if n != len(str) {
			panic("write fail")
		}

		writeTmp[j] = readTmp[0]
		j++
		j %= bytesPreLine
		if j == 0 {
			str = ": " + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(string(writeTmp), "\n", ""), "\t", ""), "\r", "") + "\n"
			n, err = outFileObj.WriteString(str)
			if err != nil {
				panic(err)
			}
			if n != len(str) {
				panic("write fail")
			}
		} else if i+1 == inFileSize {
			str = ": " + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(string(writeTmp[:j]), "\n", ""), "\t", ""), "\r", "")
			n, err = outFileObj.WriteString(str)
			if err != nil {
				panic(err)
			}
			if n != len(str) {
				panic("write fail")
			}
		} else {
			n, err = outFileObj.WriteString(" ")
			if err != nil {
				panic(err)
			}
			if n != len(" ") {
				panic("write fail")
			}
		}
	}
}

// (?:...) 非捕获组
func regexpNonCapture() {
	r1 := regexp.MustCompile(`^(?:a|b)c$`)
	r2 := regexp.MustCompile(`^(a|b)c$`)
	fmt.Println(r1.FindStringSubmatch("ac")) // [ac]
	fmt.Println(r2.FindStringSubmatch("ac")) // [ac a]
}

func createPng() {
	width := 1024
	height := 768
	rect := image.Rect(0, 0, width, height)
	rgba := image.NewRGBA(rect)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			rgba.SetRGBA(x, y, color.RGBA{A: 255})
		}
	}
	buf := &bytes.Buffer{}
	err := png.Encode(buf, rgba)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(`black.png`, buf.Bytes(), 0666)
	if err != nil {
		panic(err)
	}
}

func getMachineIPv4() (res []string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, ife := range interfaces {
		addrs, err := ife.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, _ := addr.(*net.IPNet)
			if ipNet != nil {
				to4 := ipNet.IP.To4()
				if to4 != nil {
					if !to4.IsLoopback() {
						res = append(res, to4.String())
					}
				}
			}
		}
	}

	return
}

func primeSieve() {
	generate := func(ch chan<- int) {
		for i := 2; ; i++ {
			ch <- i
		}
	}
	filter := func(src <-chan int, dst chan<- int, prime int) {
		for i := range src {
			if i%prime != 0 {
				dst <- i
			}
		}
	}

	ch := make(chan int)
	go generate(ch)
	for {
		prime := <-ch
		fmt.Printf("%d\n", prime)
		time.Sleep(time.Millisecond * 100)
		ch1 := make(chan int)
		go filter(ch, ch1, prime)
		ch = ch1
	}
}

// 两数乘积 = 两数的最大公约数 * 两数最小公倍数。
// 最小公倍数：两数共用的因数和各自的因数相乘。
func greatestCommonDivisor[T Integer](x, y T) (gcd T, leastCommonMultiple T) {
	mul := x * y

	for y != 0 {
		x, y = y, x%y
	}
	return x, mul / x
}

type Integer interface {
	int | int8 | int32 | int64 | UInteger
}

func isPowerOfTwo[T Integer](v T) bool {
	return (v-1)&v == 0
}

type UInteger interface {
	uint | uint8 | uint32 | uint64 | uintptr
}

func isMultiplicationOverflow[T UInteger](x, y T) bool {
	typeSize := unsafe.Sizeof(x)
	if uint64(x|y) < uint64(1<<4*typeSize) { // 该类型的一半位数量的最大值的平方不会大于该类型最大数。
		return false
	}
	return x > (^T(0))/y // 一个数超过了类型最大数的一半。
}

func pointerSize() {
	const PtrSize = 4 << (^uintptr(0) >> 63)
	fmt.Println("pointer size is", PtrSize) // 4 或者 8。
}

func garbageNoFree() {
	type S struct {
		Field *S
	}

	func() {
		f := new(S)
		s := new(S)
		s.Field = f
		f.Field = s
		runtime.SetFinalizer(s, func(s *S) {
			fmt.Println("s is freed")
		})
		runtime.SetFinalizer(f, func(s *S) {
			fmt.Println("f is freed")
		})
	}()
	for i := 0; i < 10; i++ {
		runtime.GC()
		time.Sleep(time.Second)
	}
}

func garbage() {
	type S struct {
		Field *int
	}

	func() {
		f := new(int)
		s := new(S)
		s.Field = f
		runtime.SetFinalizer(s, func(s *S) {
			fmt.Println("s is freed")
		})
		runtime.SetFinalizer(f, func(s *int) {
			fmt.Println("f is freed")
		})
	}()
	for i := 0; i < 10; i++ {
		runtime.GC()
		time.Sleep(time.Second)
	}
}

func bitPrint[T any](t T) {
	size := int64(unsafe.Sizeof(t))
	ptr := unsafe.Pointer(&t)
	var b byte
	for i := size - 1; i > 0; i-- {
		b = *(*byte)(unsafe.Add(ptr, i))
		for j := 7; j >= 0; j-- {
			fmt.Printf("%d", b>>j&0b1)
		}
		fmt.Print("_")
	}
	b = *(*byte)(ptr)
	for j := 7; j >= 0; j-- {
		fmt.Printf("%d", b>>j&0b1)
	}
	fmt.Println()
}
