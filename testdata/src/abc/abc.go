package abc

type Ta struct {
	Tb Tb
}

func (Ta) M1() {}
func (Ta) M2() {}
func Fa() Ta {
	var ta Ta
	return ta
}

type Tb int

func (Tb) M1() {}
func FCommon() {
	var ta Ta
	var tb Tb

	if ta.Tb == 0 || tb == 0 {
		return
	}
	return
}

func FCommonParams(ta Ta, tb Tb) {
	return
}

type tNonExported int

func (tNonExported) M1()    {}
func (tNonExported) M2()    {}
func (t *tNonExported) M3() {}
func (t *tNonExported) M4() {}
func F() *tNonExported {
	t := tNonExported(1)
	return &t
}
