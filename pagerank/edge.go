package pagerank

type Edge struct {
	src int32
	dst int32
}

func (e Edge) GetSrc() int32 {
	return e.src
}

func (e Edge) GetDst() int32 {
	return e.dst
}
