package element

// note: not thourougly tested on moduli != .NoCarry
const FromMont = `
// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func (z *{{.ElementName}}) FromMont() *{{.ElementName}} {

	// the following lines implement z = z * 1
	// with a modified CIOS montgomery multiplication
	{{- range $j := .NbWordsIndexesFull}}
	{
		// m = z[0]n'[0] mod W
		m := z[0] * {{index $.QInverse 0}}
		C := madd0(m, {{index $.Q 0}}, z[0])
		{{- range $i := $.NbWordsIndexesNoZero}}
			C, z[{{sub $i 1}}] = madd2(m, {{index $.Q $i}}, z[{{$i}}], C)
		{{- end}}
		z[{{sub $.NbWords 1}}] = C
	}
	{{- end}}

	{{ template "reduce" .}}
	return z 
}
`

const Conv = `
// ToMont converts z to Montgomery form
// sets and returns z = z * r^2
func (z *{{.ElementName}}) ToMont() *{{.ElementName}} {
	var rSquare = {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	return z.MulAssign(&rSquare)
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
	q := {{toLower .ElementName}}ModulusBigInt()

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

	// while v < 0, v+=q 
	for vv.Cmp(zero) == -1 {
		vv.Add(vv, q)
	}
	// while v > q, v-=q
	for vv.Cmp(q) == 1 {
		vv.Sub(vv, q)
	}
	// if v == q, return 0
	if vv.Cmp(q) == 0 {
		return z
	}
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
