package logic

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	. "github.com/weiwenchong/calculator/common"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

type Calculator struct {
	Expression string `json:"expression"`
}

func Calculate(c *gin.Context) {
	cal := new(Calculator)
	if err := c.ShouldBind(&cal); err != nil {
		c.JSON(http.StatusOK, ApiReturn{
			Ret : -1,
			Msg: "request paramenters format error",
		})
		return
	}

	res, err := cal.Calculate(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, ApiReturn{
			Ret : -2,
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ApiReturn{
		Ret: 1,
		Data: &ApiData{
			Ext: res,
			Ent: nil,
		},
	})
}

func (c *Calculator) Calculate(ctx context.Context) (res string, err error) {
	fun := "Calculator.Calculate -->"
	log.Printf("%s incall", fun)

	err = c.valid(ctx)
	if err != nil {
		return
	}

	stacks := make([]*calStack, 0)
	stacks = append(stacks, &calStack{
		Operator: '+',
		Operands: make([]int64, 0),
	})
	var cur int64 = 0

	for index, s := range c.Expression {
		flag := true
		if unicode.IsDigit(s) {
			n, _ := strconv.ParseInt(string(s), 10, 64)
			cur = cur*10 + n
		} else if s == '(' {
			stacks = append(stacks, &calStack{
				Operator: '+',
				Operands: make([]int64, 0),
			})
		} else if s == ')' {
			stacks[len(stacks)-1].calculate(cur)
			stacks[len(stacks)-1].Operator = s
			cur = 0
			cur = stacks[len(stacks)-1].sum()
			stacks = stacks[:len(stacks)-1]
		} else {
			flag = false
		}
		if !flag || index == len(c.Expression)-1 {
			stacks[len(stacks)-1].calculate(cur)
			stacks[len(stacks)-1].Operator = s
			cur = 0
		}
	}

	res = strconv.FormatInt(stacks[0].sum(), 10)

	log.Printf("%s succeed req:%s, rev:%s", fun, c.Expression, res)
	return
}

func (c *Calculator) valid(ctx context.Context) error {
	// 剔除空格
	// 是否有非法字符
	// 检查是否有连续两个计算符
	// 检查是否括号是否匹配
	c.Expression = strings.ReplaceAll(c.Expression, " ", "")

	bracketStack := make([]bool, 0)
	// 0:无 1:数字 2:计算符 3:前括号 4:后括号
	pre := 0

	err := errors.New("invalid expression")

	for _, s := range c.Expression {
		if unicode.IsDigit(s) {
			pre = 1
			continue
		} else if s == '(' {
			if pre == 1 {
				return err
			}
			pre = 3
			bracketStack = append(bracketStack, true)
		} else if s == ')' {
			if len(bracketStack) == 0 || bracketStack[len(bracketStack)-1] != true || pre == 2 {
				return err
			}
			bracketStack = bracketStack[:len(bracketStack)-1]
			pre = 4
		} else if s == '+' || s == '-' {
			if pre == 2 {
				return err
			}
			pre = 2
		} else if s == '*' || s == '/' {
			if pre == 2 || pre == 3 || pre == 0 {
				return err
			}
			pre = 2
		} else {
			return err
		}
	}
	if len(bracketStack) != 0 {
		return err
	}
	return nil
}

type calStack struct {
	Operator int32
	Operands []int64
}

func (c *calStack) push(s int64) {
	c.Operands = append(c.Operands, s)
}

func (c *calStack) pop() int64 {
	if len(c.Operands) == 0 {
		return 0
	}
	res := c.Operands[len(c.Operands)-1]
	c.Operands = c.Operands[:len(c.Operands)-1]
	return res
}

func (c *calStack) len() int {
	return len(c.Operands)
}

func (c *calStack) sum() (res int64) {
	for _, s := range c.Operands {
		res += s
	}
	return
}

func (c *calStack) calculate(num int64) {
	switch c.Operator {
	case '+':
		c.push(num)
	case '-':
		c.push(-num)
	case '*':
		c.push(c.pop()*num)
	case '/':
		c.push(c.pop()/num)
	}
}

