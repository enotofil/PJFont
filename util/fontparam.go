package util

import "fmt"

type fontParam struct {
	name  string
	value int
	min   int
	max   int
	step  int
}

func (fp *fontParam) change(amount int) {
	fp.value += amount * fp.step
	if fp.value > fp.max {
		fp.value = fp.max
	} else if fp.value < fp.min {
		fp.value = fp.min
	}
}

func (fp *fontParam) toString() string {
	str := fmt.Sprintf("%22s : %d", fp.name, fp.value)
	if fp.max == 1 {
		str = fmt.Sprintf("%22s : %t", fp.name, fp.value == 1)
	}
	return str
}
