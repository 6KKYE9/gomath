package main

import (
	"math/big"
	"testing"
)

func TestCalc(t *testing.T) {
	cases := []struct {
		expr string
		want float64
	}{
		{"(2+3)*4", 20},
		{"2+3*4", 14},
		{"10/4", 2.5},
		{"-5+3", -2},
		{"3 - -2", 5},
		{"((1+2)*3)", 9},
		{"2*-3", -6},
		{"1.5+2.5", 4},
	}
	for _, c := range cases {
		got, err := Calc(c.expr)
		if err != nil {
			t.Fatalf("Calc(%q) 报错: %v", c.expr, err)
		}
		if got != c.want {
			t.Fatalf("Calc(%q) = %v, want %v", c.expr, got, c.want)
		}
	}
}

func TestCalcErrors(t *testing.T) {
	bad := []string{
		"2+",          // 不完整
		"(2+3",        // 缺右括号
		"2 3",         // 多余字符
		"2/0",         // 除零
		"abc",         // 非法字符
	}
	for _, e := range bad {
		if _, err := Calc(e); err == nil {
			t.Fatalf("Calc(%q) 应报错", e)
		}
	}
}

func TestIsPrime(t *testing.T) {
	cases := map[int]bool{
		0: false, 1: false, 2: true, 3: true, 4: false,
		17: true, 18: false, 97: true, 100: false,
	}
	for n, want := range cases {
		if got := IsPrime(n); got != want {
			t.Fatalf("IsPrime(%d) = %v, want %v", n, got, want)
		}
	}
}

func TestFib(t *testing.T) {
	cases := map[int]int{0: 0, 1: 1, 2: 1, 10: 55, 20: 6765}
	for n, want := range cases {
		if got := Fib(n); got != want {
			t.Fatalf("Fib(%d) = %d, want %d", n, got, want)
		}
	}
}

func TestFact(t *testing.T) {
	if Fact(5).String() != "120" {
		t.Fatalf("Fact(5) = %s, want 120", Fact(5).String())
	}
	// 21! 超过 int64 上限，用 big 仍能正确表示
	big21 := "51090942171709440000"
	if Fact(21).String() != big21 {
		t.Fatalf("Fact(21) = %s, want %s", Fact(21).String(), big21)
	}
}

func TestGcd(t *testing.T) {
	cases := []struct{ a, b, want int }{
		{12, 18, 6}, {18, 12, 6}, {7, 13, 1}, {0, 5, 5},
	}
	for _, c := range cases {
		if got := Gcd(c.a, c.b); got != c.want {
			t.Fatalf("Gcd(%d,%d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestLcm(t *testing.T) {
	cases := []struct{ a, b, want int }{
		{4, 6, 12}, {12, 18, 36}, {0, 5, 0},
	}
	for _, c := range cases {
		if got := Lcm(c.a, c.b); got != c.want {
			t.Fatalf("Lcm(%d,%d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestComb(t *testing.T) {
	c, err := Comb(10, 3)
	if err != nil {
		t.Fatalf("Comb(10,3) 报错: %v", err)
	}
	if c.String() != "120" {
		t.Fatalf("C(10,3) = %s, want 120", c.String())
	}
	// 大组合数不溢出（100891344545564193334812497256 为正确值）
	bigC, err := Comb(100, 50)
	if err != nil {
		t.Fatalf("Comb(100,50) 报错: %v", err)
	}
	want := "100891344545564193334812497256"
	if bigC.String() != want {
		t.Fatalf("C(100,50) = %s, want %s", bigC.String(), want)
	}
	if _, err := Comb(5, 9); err == nil {
		t.Fatal("Comb(5,9) 参数非法应报错")
	}
}

func TestFactEqualsBigInt(t *testing.T) {
	// 确认 Fact 返回 *big.Int 且与 big 直接计算一致
	if Fact(0).Cmp(big.NewInt(1)) != 0 {
		t.Fatal("Fact(0) 应为 1")
	}
}
