// package A
package A

import (
	"DairoNPS/BB"
	"fmt"
)

type A struct {
	abc int
}

func (a *A) Out() {
	a.abc = 10
	fmt.Println("this is A class")
}

func OutToB(bb BB.BB) {
	bb.Out()
}
