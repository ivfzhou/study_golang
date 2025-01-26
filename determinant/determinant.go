package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// [["1","4","9","16"],["4","9","16","25"],["9","16","25","36"],["16","25","36","49"]] = 0
// [["1","1","1","1"],["1","2","3","4"],["1","3","6","10"],["1","4","10","20"]] = 1
// [["5","0","4","7"],["1","-1","2","1"],["4","1","2","0"],["1","1","1","1"]] = -234/5
// [["0","1","2","4","1"],["2","0","1","1","3"],["-1","3","5","2","6"],["1","1","1","6","6"]] = 0
// [["1","1/2","0","1","-1"],["2","0","-1","1","1"],["3","2","1","1/2","-1/2"],["1","-1","0","1","2"]] = 0
func TestDeterminant() {
	var data [][]string
	bs := bufio.NewReader(os.Stdin)
	fmt.Println("输入一行JSON数据；小数请输入最简分式")
	line, err := bs.ReadString('\n')
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(line), &data)
	if err != nil {
		panic(err)
	}
	fmt.Println(Calculate(data))
}

// Calculate 计算行列式值
func Calculate(data [][]string) string {
	// 校验数据
	if data == nil || len(data) == 0 || len(data[0]) != len(data) {
		return "0"
	}
	size := len(data[0])
	for i := 1; i < len(data); i++ {
		if len(data[i]) < size {
			panic("每行数个数不一致")
		}
	}

	// 复制数据
	tmp := make([][][2]int64, 0, len(data))
	for i := range data {
		arr := make([][2]int64, len(data[0]))
		for j := range arr {
			fraction := strings.Split(data[i][j], "/")
			if len(fraction) == 2 {
				num1, err := strconv.ParseInt(fraction[0], 10, 64)
				if err != nil {
					panic(err)
				}
				num2, err := strconv.ParseInt(fraction[1], 10, 64)
				if err != nil {
					panic(err)
				}
				arr[j] = [2]int64{num1, num2}
			} else if len(fraction) == 1 {
				num1, err := strconv.ParseInt(fraction[0], 10, 64)
				if err != nil {
					panic(err)
				}
				arr[j] = [2]int64{num1, 1}
			} else {
				panic("数据非法 " + data[i][j])
			}
		}
		tmp = append(tmp, arr)
	}
	data = nil

	// 系数
	coefficient := [2]int64{1, 1}
	for len(tmp) != 1 {
		fmt.Println(coefficient)
		for i := range tmp {
			bytes, _ := json.Marshal(tmp[i])
			fmt.Println(string(bytes))
		}
		fmt.Println()

		// 处理第一个元素非0
		isPositive := 0
		for i := 1; tmp[0][0][0] == 0 && i < len(tmp); i++ {
			tmp[0], tmp[i] = tmp[i], tmp[0]
			isPositive++
		}
		if tmp[0][0][0] == 0 {
			return "0"
		}

		// 系数
		if isPositive%2 != 0 {
			coefficient[0] *= -1
		}

		// 提取每行第一个数为系数
		for i := range tmp {
			factor := tmp[i][0]
			if factor[0] == 0 {
				continue
			}
			for j := 1; j < len(tmp[i]); j++ {
				if tmp[i][j][0] == 0 || factor[0] == 0 {
					tmp[i][j] = [2]int64{0, 1}
				} else {
					tmp[i][j][0] *= factor[1]
					tmp[i][j][1] *= factor[0]
					tmp[i][j] = divide(tmp[i][j])
				}
			}
			coefficient[0] *= factor[0]
			coefficient[1] *= factor[1]
			coefficient = divide(coefficient)
		}

		// 每行依次减第一行
		for i := 1; i < len(tmp); i++ {
			for j := 1; j < len(tmp[i]); j++ {
				denominator := tmp[i][j][1] * tmp[0][j][1]
				tmp[i][j][0] = tmp[i][j][0]*tmp[0][j][1] - tmp[0][j][0]*tmp[i][j][1]
				tmp[i][j][1] = denominator
				tmp[i][j] = divide(tmp[i][j])
			}
		}

		// 除第一列每列提取第一个数为系数
		for i := 1; i < len(tmp); i++ {
			for j := 1; j < len(tmp[i]); j++ {
				if tmp[0][j][0] == 0 {
					continue
				}
				if i == 1 {
					coefficient[0] *= tmp[0][j][0]
					coefficient[1] *= tmp[0][j][1]
					coefficient = divide(coefficient)
				}
				if tmp[i][j][0] == 0 || tmp[0][j][0] == 0 {
					tmp[i][j] = [2]int64{0, 1}
				} else {
					tmp[i][j][0] *= tmp[0][j][1]
					tmp[i][j][1] *= tmp[0][j][0]
					tmp[i][j] = divide(tmp[i][j])
				}
			}
		}

		// tmp切除第一行第一列
		for i := 0; i < len(tmp)-1; i++ {
			tmp[i] = tmp[i+1][1:]
		}
		tmp = tmp[:len(tmp)-1]
	}

	member := tmp[0][0][0] * coefficient[0]
	denominator := tmp[0][0][1] * coefficient[1]
	num := divide([2]int64{member, denominator})

	if num[1] == 1 {
		return strconv.FormatInt(num[0], 10)
	} else {
		return fmt.Sprintf("%d/%d", num[0], num[1])
	}
}

// 最大公约数 最小公倍数。
func gcd(x, y int64) int64 {
	for y != 0 {
		x, y = y, x%y
	}

	return x
}

func divide(num [2]int64) [2]int64 {
	member := num[0]
	denominator := num[1]
	if member == 0 {
		return [2]int64{0, 1}
	}

	// 处理符号
	symbol := int64(1)
	if member < 0 && denominator < 0 {
		member = -member
		denominator = -denominator
	} else if member > 0 && denominator < 0 {
		symbol = -1
		denominator = -denominator
	} else if member < 0 && denominator > 0 {
		symbol = -1
		member = -member
	}

	// 约分
	d := gcd(member, denominator)
	member /= d
	denominator /= d
	return [2]int64{symbol * member, denominator}
}
