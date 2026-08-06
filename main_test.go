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
		{"10/4", 2.5},
		{"2+3*4", 14},
		{"-5+2", -3},
		{"2*(3+(4-1))", 12},
	}
	for _, c := range cases {
		got, err := Calc(c.expr)
		if err != nil {
			t.Errorf("Calc(%q) error: %v", c.expr, err)
			continue
		}
		if got != c.want {
			t.Errorf("Calc(%q)=%v want %v", c.expr, got, c.want)
		}
	}
	if _, err := Calc("2+"); err == nil {
		t.Error("不完整表达式应报错")
	}
	if _, err := Calc("2 @ 3"); err == nil {
		t.Error("非法字符应报错")
	}
}

func TestIsPrime(t *testing.T) {
	if IsPrime(17) != true || IsPrime(1) != false || IsPrime(0) != false || IsPrime(9) != false {
		t.Error("IsPrime 判定错误")
	}
}

func TestFib(t *testing.T) {
	want := []int{0, 1, 1, 2, 3, 5, 8, 13}
	for i, w := range want {
		if Fib(i) != w {
			t.Errorf("Fib(%d)=%d want %d", i, Fib(i), w)
		}
	}
}

func TestGcdLcm(t *testing.T) {
	if Gcd(12, 18) != 6 {
		t.Errorf("Gcd(12,18)=%d want 6", Gcd(12, 18))
	}
	if Lcm(4, 6) != 12 {
		t.Errorf("Lcm(4,6)=%d want 12", Lcm(4, 6))
	}
	if Lcm(0, 5) != 0 {
		t.Errorf("Lcm(0,5) 应为 0")
	}
}

func TestFactComb(t *testing.T) {
	if Fact(5).String() != "120" {
		t.Errorf("Fact(5)=%s want 120", Fact(5).String())
	}
	c, err := Comb(10, 3)
	if err != nil {
		t.Fatal(err)
	}
	if c.String() != "120" {
		t.Errorf("Comb(10,3)=%s want 120", c.String())
	}
	if _, err := Comb(3, 5); err == nil {
		t.Error("Comb 越界应报错")
	}
}

func TestIntPow(t *testing.T) {
	if intPow(2, 10) != 1024 {
		t.Errorf("2^10=%d want 1024", intPow(2, 10))
	}
	if intPow(3, 0) != 1 || intPow(5, 1) != 5 {
		t.Error("intPow 边界错误")
	}
}

// 确保 big.Int 用法仍能编译
var _ = big.NewInt
