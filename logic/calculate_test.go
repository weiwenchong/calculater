package logic

import (
	"context"
	"fmt"
	"testing"
)

func TestCalculate(t *testing.T) {
	c := &Calculator{Expression: "-3*(-1+2)-(3*(2+5)+10/(2+3))"}
	res, err := c.Calculate(context.TODO())
	fmt.Println(res, err)
}