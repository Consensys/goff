package element

const Conv = `

// ToMont converts z to Montgomery form
// sets and returns z = z * r^2
func (z *{{.ElementName}}) ToMont() *{{.ElementName}} {
	return z.Mul(z, &rSquare)
}

// ToRegular returns z in regular form (doesn't mutate z)
func (z {{.ElementName}}) ToRegular() {{.ElementName}} {
	return *z.FromMont()
}

// String returns the string form of an {{.ElementName}} in Montgomery form
func (z *{{.ElementName}}) String() string {
	var _z big.Int
	return z.ToBigIntRegular(&_z).String()
}

// ToBigInt returns z as a big.Int in Montgomery form 
func (z *{{.ElementName}}) ToBigInt(res *big.Int) *big.Int {
	bits := (*[{{.NbWords}}]big.Word)(unsafe.Pointer(z))
	return res.SetBits(bits[:])
}

// ToBigIntRegular returns z as a big.Int in regular form 
func (z {{.ElementName}}) ToBigIntRegular(res *big.Int) *big.Int {
	z.FromMont()
	bits := (*[{{.NbWords}}]big.Word)(unsafe.Pointer(&z))
	return res.SetBits(bits[:])
}

// SetBigInt sets z to v (regular form) and returns z in Montgomery form
func (z *{{.ElementName}}) SetBigInt(v *big.Int) *{{.ElementName}} {
	z.SetZero()

	zero := big.NewInt(0)
	q := Modulus()

	// fast path
	c := v.Cmp(q)
	if c == 0 {
		return z
	} else if c != 1 && v.Cmp(zero) != -1 {
		// v should
		vBits := v.Bits()
		for i := 0; i < len(vBits); i++ {
			z[i] = uint64(vBits[i])
		}
		return z.ToMont()
	}
	
	// copy input
	vv := new(big.Int).Set(v)
	vv.Mod(v, q)
	
	// v should
	vBits := vv.Bits()
	for i := 0; i < len(vBits); i++ {
		z[i] = uint64(vBits[i])
	}
	return z.ToMont()
}

// SetString creates a big.Int with s (in base 10) and calls SetBigInt on z
func (z *{{.ElementName}}) SetString( s string) *{{.ElementName}} {
	x, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("{{.ElementName}}.SetString failed -> can't parse number in base10 into a big.Int")
	}
	return z.SetBigInt(x)
}

`
