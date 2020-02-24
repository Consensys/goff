package element

const Exp = `
// Exp z = x^e mod q
func (z *{{.ElementName}}) Exp(x {{.ElementName}}, e uint64) *{{.ElementName}} {
	if e == 0 {
		return z.SetOne()
	} 

	z.Set(&x)

	l := bits.Len64(e) - 2
	for i := l; i >= 0; i-- {
		z.Square(z)
		if e&(1<<uint(i)) != 0 {
			z.MulAssign(&x)
		}
	}
	return z
}
`
