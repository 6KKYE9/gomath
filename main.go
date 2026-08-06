// gomath：一个命令行小数学工具
//
// 功能：
//   1. 计算四则运算表达式（支持 + - * / 和括号），比如  (2+3)*4
//   2. 质数判断、斐波那契、阶乘、最大公约数（GCD）
//
// 全部用标准库实现。重点看 evalExpr —— 它用"递归下降"的方法把字符串表达式算成数字，
// 这是编译原理里最朴素也最清楚的写法，比一堆正则替换靠谱。
//
// 用法：
//   gomath calc "(2+3)*4"
//   gomath prime 17
//   gomath fib 10
//   gomath fact 5
//   gomath gcd 12 18
package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// ===== 1) 四则表达式求值（递归下降）=====

// parser 持有一个表达式和当前读取位置。
type parser struct {
	s   string
	pos int
}

// 跳过空白。
func (p *parser) skipSpace() {
	for p.pos < len(p.s) && unicode.IsSpace(rune(p.s[p.pos])) {
		p.pos++
	}
}

// 读一个数字（支持小数）。解析失败返回 error 而非直接退出。
func (p *parser) parseNumber() (float64, error) {
	start := p.pos
	for p.pos < len(p.s) && (unicode.IsDigit(rune(p.s[p.pos])) || p.s[p.pos] == '.') {
		p.pos++
	}
	v, err := strconv.ParseFloat(p.s[start:p.pos], 64)
	if err != nil {
		return 0, fmt.Errorf("数字解析失败: %q", p.s[start:p.pos])
	}
	return v, nil
}

// 因子：数字，或带括号的子表达式，或一元负号。
func (p *parser) parseFactor() (float64, error) {
	p.skipSpace()
	if p.pos < len(p.s) && p.s[p.pos] == '(' {
		p.pos++ // 吃掉 '('
		v, err := p.parseExpr()
		if err != nil {
			return 0, err
		}
		p.skipSpace()
		if p.pos < len(p.s) && p.s[p.pos] == ')' {
			p.pos++ // 吃掉 ')'
		} else {
			return 0, fmt.Errorf("括号不匹配，缺少 )")
		}
		return v, nil
	}
	if p.pos < len(p.s) && p.s[p.pos] == '-' {
		p.pos++
		v, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		return -v, nil
	}
	if p.pos < len(p.s) && p.s[p.pos] == '+' {
		p.pos++
		return p.parseFactor()
	}
	return p.parseNumber()
}

// 项：因子之间做 * 和 /。
func (p *parser) parseTerm() (float64, error) {
	v, err := p.parseFactor()
	if err != nil {
		return 0, err
	}
	for {
		p.skipSpace()
		if p.pos < len(p.s) && (p.s[p.pos] == '*' || p.s[p.pos] == '/') {
			op := p.s[p.pos]
			p.pos++
			rhs, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			if op == '*' {
				v *= rhs
			} else {
				if rhs == 0 {
					return 0, fmt.Errorf("除以零")
				}
				v /= rhs
			}
		} else {
			return v, nil
		}
	}
}

// 表达式：项之间做 + 和 -。这是最外层，处理完应读到字符串末尾。
func (p *parser) parseExpr() (float64, error) {
	v, err := p.parseTerm()
	if err != nil {
		return 0, err
	}
	for {
		p.skipSpace()
		if p.pos < len(p.s) && (p.s[p.pos] == '+' || p.s[p.pos] == '-') {
			op := p.s[p.pos]
			p.pos++
			rhs, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			if op == '+' {
				v += rhs
			} else {
				v -= rhs
			}
		} else {
			return v, nil
		}
	}
}

// Calc 计算一个表达式字符串，遇到多余字符或语法错误会返回 error。
func Calc(s string) (float64, error) {
	p := &parser{s: strings.TrimSpace(s)}
	v, err := p.parseExpr()
	if err != nil {
		return 0, err
	}
	p.skipSpace()
	if p.pos != len(p.s) {
		return 0, fmt.Errorf("无法解析的多余字符: %q", s[p.pos:])
	}
	return v, nil
}

// ===== 2) 几个常见数学小工具 =====

// IsPrime 判断 n 是否为质数（只试到 sqrt(n) 就够）。
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Fib 返回第 n 个斐波那契数（从 0 开始：fib(0)=0, fib(1)=1）。
func Fib(n int) int {
	a, b := 0, 1
	for i := 0; i < n; i++ {
		a, b = b, a+b
	}
	return a
}

// Fact 计算 n 的阶乘，用 math/big 避免大数溢出（n 较大时依然正确）。
// 返回 *big.Int（调用方可转字符串查看超大结果）。
func Fact(n int) *big.Int {
	r := big.NewInt(1)
	for i := 2; i <= n; i++ {
		r.Mul(r, big.NewInt(int64(i)))
	}
	return r
}

// Gcd 用辗转相除法求最大公约数。
func Gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Lcm 用公式 lcm(a,b) = |a*b| / gcd(a,b) 求最小公倍数。
func Lcm(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return abs(a*b) / Gcd(a, b)
}

// Comb 计算组合数 C(n, k) = n! / (k! * (n-k)!)，用 math/big 避免溢出。
func Comb(n, k int) (*big.Int, error) {
	if k < 0 || k > n {
		return nil, fmt.Errorf("组合数要求 0 <= k <= n")
	}
	num := Fact(n)
	den := new(big.Int).Mul(Fact(k), Fact(n-k))
	return new(big.Int).Div(num, den), nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	// 用一个"-help"之外的子命令来区分功能。
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "用法: gomath <calc|prime|fib|fact|gcd> 参数...")
		os.Exit(1)
	}
	cmd := os.Args[1]
	rest := os.Args[2:]
	// 让 flag 包在子命令后能继续解析（虽然这里大多直接读位置参数）。
	_ = flag.CommandLine.Parse([]string{})

	switch cmd {
	case "calc":
		if len(rest) == 0 {
			fmt.Fprintln(os.Stderr, "calc 需要一个表达式，如 gomath calc \"(2+3)*4\"")
			os.Exit(1)
		}
		// 允许把空格分开的参数拼回去，比如 gomath calc 2 + 3。
		expr := strings.Join(rest, " ")
		v, err := Calc(expr)
		if err != nil {
			fmt.Fprintln(os.Stderr, "计算失败:", err)
			os.Exit(1)
		}
		// 整数结果去掉尾部 .0，更美观。
		if v == float64(int64(v)) {
			fmt.Printf("%s = %d\n", expr, int64(v))
		} else {
			fmt.Printf("%s = %v\n", expr, v)
		}
	case "prime":
		n := atoi(rest)
		fmt.Printf("%d 是质数? %v\n", n, IsPrime(n))
	case "fib":
		n := atoi(rest)
		fmt.Printf("fib(%d) = %d\n", n, Fib(n))
	case "fact":
		n := atoi(rest)
		fmt.Printf("%d! = %s\n", n, Fact(n).String())
	case "gcd":
		if len(rest) < 2 {
			fmt.Fprintln(os.Stderr, "gcd 需要两个数")
			os.Exit(1)
		}
		a, b := atoi(rest[:1]), atoi(rest[1:2])
		fmt.Printf("gcd(%d, %d) = %d\n", a, b, Gcd(a, b))
	case "lcm":
		if len(rest) < 2 {
			fmt.Fprintln(os.Stderr, "lcm 需要两个数")
			os.Exit(1)
		}
		a, b := atoi(rest[:1]), atoi(rest[1:2])
		fmt.Printf("lcm(%d, %d) = %d\n", a, b, Lcm(a, b))
	case "comb":
		if len(rest) < 2 {
			fmt.Fprintln(os.Stderr, "comb 需要两个数 n k")
			os.Exit(1)
		}
		n, k := atoi(rest[:1]), atoi(rest[1:2])
		c, err := Comb(n, k)
		if err != nil {
			fmt.Fprintln(os.Stderr, "comb 出错:", err)
			os.Exit(1)
		}
		fmt.Printf("C(%d, %d) = %s\n", n, k, c.String())
	default:
		fmt.Fprintln(os.Stderr, "未知子命令:", cmd)
		os.Exit(1)
	}
}

// atoi 从参数里取第一个数（简单封装，出错即退出）。
func atoi(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "缺少数字参数")
		os.Exit(1)
	}
	v, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "参数不是合法整数:", args[0])
		os.Exit(1)
	}
	return v
}
