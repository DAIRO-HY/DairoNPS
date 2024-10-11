// package B
package B

import (
	"DairoNPS/AA"
	"fmt"
)

type B struct{}

func (b *B) Out() {
	fmt.Println("this is B class")
}

func OutToA(aa AA.AA) {
	aa.Out()
}
