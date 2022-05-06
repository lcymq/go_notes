package main

import "fmt"

func main() {
	a := [][]int{}
	n := 10
	for i := 0; i < n; i++ {
		tmp := make([]int, n)
		a = append(a, tmp)
	}
	tournament(a, n)
	for _, v := range a {
		fmt.Println(v)
	}
	fmt.Println()
	res := getRes(a, n)
	for _, v := range res {
		fmt.Println(v)
	}
}

func tournament(a [][]int, n int) {
	if n == 1 {
		a[0][0] = 1
		return
	}
	if n == 2 {
		a[0][0] = 1
		a[0][1] = 2
		a[1][0] = 2
		a[1][1] = 1
		return
	}
	if odd(n) {
		tournament(a, n+1)
		return
	} else {
		tournament(a, n/2)
		makecopy(a, n)
	}
}

func makecopy(a [][]int, n int) {
	m := n / 2
	if odd(m) {
		copyodd(a, n)
	} else {
		copy(a, n)
	}
}

func copy(a [][]int, n int) {
	m := n / 2
	for i := 0; i < m; i++ {
		for j := 0; j < m; j++ {
			a[i][j+m] = a[i][j] + m
			a[i+m][j] = a[i][j] + m
			a[i+m][j+m] = a[i][j]
		}
	}
}

func copyodd(a [][]int, n int) {
	m := n / 2
	b := make([]int, n)
	for i := 0; i < m; i++ {
		b[i] = m + i + 1
		b[m+i] = b[i]
	}
	for i := 0; i < m; i++ {
		for j := 0; j <= m; j++ {
			if a[i][j] > m {
				a[i][j] = b[i]
				a[i+m][j] = (b[i] + m) % n
			} else {
				a[i+m][j] = a[i][j] + m
			}
		}
		for j := 1; j < m; j++ {
			a[i][j+m] = b[i+j]
			a[b[i+j]-1][j+m] = i + 1
		}
	}
}

func odd(n int) bool {
	return n%2 == 1
}

func getRes(a [][]int, n int) [][][]int {
	res := make([][][]int, 0)
	for j := 1; j < n; j++ {
		tmpRes := make([][]int, 0)
		m := make(map[int]bool)
		for i := 0; i < n; i++ {
			if m[a[i][0]] || m[a[i][j]] {
				continue
			}
			tmpRes = append(tmpRes, []int{a[i][0], a[i][j]})
			m[a[i][0]] = true
			m[a[i][j]] = true
		}
		res = append(res, tmpRes)
	}
	return res
}
