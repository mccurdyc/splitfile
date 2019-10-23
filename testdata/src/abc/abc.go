package abc

type (
	A int
	B int
	C int
)

func fabc(a A, b B, c C) {}

func (a A) ma(b B) {}
func (b B) mb()    {}
func (c C) mc(b B) {}
