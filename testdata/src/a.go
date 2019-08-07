package a

type Ta struct {
	A int
}

func (Ta) M1() {}

func (Ta) M2() {}

func Fa() {
	var ta Ta

	if ta.A == 0 {
	}
}

type Tb struct {
	A int
}

func (Tb) M1() {}

// should go to a common file
func Fcommon() {
	var ta Ta
	var tb Tb

	if ta.A == 0 || tb.A == 0 {
	}
}

type t struct{}

func (t) M1() {}

func (t) M2() {}

func (t *t) M3() {}

func (t *t) M4() {}

func F() {
	tp := &t{}

	if tp == nil {
	}
}
