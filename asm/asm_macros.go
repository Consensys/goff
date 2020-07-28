package asm

func _mov(i1, i2 interface{}, offsets ...int) {
	var o1, o2 int
	if len(offsets) >= 1 {
		o1 = offsets[0]
		if len(offsets) >= 2 {
			o2 = offsets[1]
		}
	}
	switch c1 := i1.(type) {
	case []uint64:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []register:
			for i := 0; i < nbWords; i++ {
				movq(c1[i+o1], c2[i+o2])
			}
		}
	case register:
		switch c2 := i2.(type) {
		case register:
			for i := 0; i < nbWords; i++ {
				movq(c1.at(i+o1), c2.at(i+o2))
			}
		case []register:
			for i := 0; i < nbWords; i++ {
				movq(c1.at(i+o1), c2[i+o2])
			}
		default:
			panic("unsupported")
		}
	case []register:
		switch c2 := i2.(type) {
		case register:
			for i := 0; i < nbWords; i++ {
				movq(c1[i+o1], c2.at(i+o2))
			}
		case []register:
			for i := 0; i < nbWords; i++ {
				movq(c1[i+o1], c2[i+o2])
			}
		default:
			panic("unsupported")
		}
	default:
		panic("unsupported")
	}

}

func _add(i1, i2 interface{}, offsets ...int) {
	var o1, o2 int
	if len(offsets) >= 1 {
		o1 = offsets[0]
		if len(offsets) >= 2 {
			o2 = offsets[1]
		}
	}
	switch c1 := i1.(type) {

	case register:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []register:
			for i := 0; i < nbWords; i++ {
				if i == 0 {
					addq(c1.at(i+o1), c2[i+o2])
				} else {
					adcq(c1.at(i+o1), c2[i+o2])
				}
			}
		}
	case []register:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []register:
			for i := 0; i < nbWords; i++ {
				if i == 0 {
					addq(c1[i+o1], c2[i+o2])
				} else {
					adcq(c1[i+o1], c2[i+o2])
				}
			}
		}
	default:
		panic("unsupported")
	}
}

func _sub(i1, i2 interface{}, offsets ...int) {
	var o1, o2 int
	if len(offsets) >= 1 {
		o1 = offsets[0]
		if len(offsets) >= 2 {
			o2 = offsets[1]
		}
	}
	switch c1 := i1.(type) {

	case register:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []register:
			for i := 0; i < nbWords; i++ {
				if i == 0 {
					subq(c1.at(i+o1), c2[i+o2])
				} else {
					sbbq(c1.at(i+o1), c2[i+o2])
				}
			}
		}
	case []register:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []register:
			for i := 0; i < nbWords; i++ {
				if i == 0 {
					subq(c1[i+o1], c2[i+o2])
				} else {
					sbbq(c1[i+o1], c2[i+o2])
				}
			}
		}
	default:
		panic("unsupported")
	}
}
