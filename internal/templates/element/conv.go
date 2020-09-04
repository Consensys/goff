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
	var b [Limbs*8]byte
	{{- range $i := reverse .NbWordsIndexesFull}}
		{{- $j := mul $i 8}}
		{{- $k := sub $.NbWords 1}}
		{{- $k := sub $k $i}}
		{{- $jj := add $j 8}}
		binary.BigEndian.PutUint64(b[{{$j}}:{{$jj}}], z[{{$k}}])
	{{- end}}

	return res.SetBytes(b[:])
}

// ToBigIntRegular returns z as a big.Int in regular form 
func (z *{{.ElementName}}) ToBigIntRegular(res *big.Int) *big.Int {
	zRegular := z.ToRegular()
	var b [Limbs*8]byte
	{{- range $i := reverse .NbWordsIndexesFull}}
		{{- $j := mul $i 8}}
		{{- $k := sub $.NbWords 1}}
		{{- $k := sub $k $i}}
		{{- $jj := add $j 8}}
		binary.BigEndian.PutUint64(b[{{$j}}:{{$jj}}], zRegular[{{$k}}])
	{{- end}}
	return res.SetBytes(b[:])
}

// SetBigInt sets z to v (regular form) and returns z in Montgomery form
func (z *{{.ElementName}}) SetBigInt(v *big.Int) *{{.ElementName}} {
	z.SetZero()

	var zero big.Int 
	q := Modulus()

	// fast path
	c := v.Cmp(q)
	if c == 0 {
		// v == 0
		return z
	} else if c != 1 && v.Cmp(&zero) != -1 {
		// 0 < v < q 
		return z.setBigInt(v)
	}
	
	// copy input + modular reduction
	vv := new(big.Int).Set(v)
	vv.Mod(v, q)
	
	return z.setBigInt(vv)
}

// setBigInt assumes 0 <= v < q 
func (z *{{.ElementName}}) setBigInt(v *big.Int) *{{.ElementName}} {
	vBits := v.Bits()

	if bits.UintSize == 64 {
		for i := 0; i < len(vBits); i++ {
			z[i] = uint64(vBits[i])
		}
	} else {
		for i := 0; i < len(vBits); i++ {
			if i%2 == 0 {
				z[i/2] = uint64(vBits[i])
			} else {
				z[i/2] |= uint64(vBits[i]) << 32
			}
		}
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
