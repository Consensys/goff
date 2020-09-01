package asm

func generateMulByNonResidueE2BN256() {
	// 	var a, b fp.Element
	// 	a.Double(&x.A0).Double(&a).Double(&a).Add(&a, &x.A0).Sub(&a, &x.A1)
	// 	b.Double(&x.A1).Double(&b).Double(&b).Add(&b, &x.A1).Add(&b, &x.A0)
	// 	z.A0.Set(&a)
	// 	z.A1.Set(&b)
	fnHeader("mulNonRes"+elementName, 0, 16)

	a := popRegisters(nbWords)
	b := popRegisters(nbWords)
	x := popRegister()

	movq("x+8(FP)", x)
	_mov(x, a) // a = a0

	_add(a, a)
	_reduce(a, a)

	_add(a, a)
	_reduce(a, a)

	_add(a, a)
	_reduce(a, a)

	_add(x, a)
	_reduce(a, a)

	_mov(x, b, nbWords) // b = a1
	_sub(b, a)
	_reduceAfterSub(a, true)

	_add(b, b)
	_reduce(b, b)

	_add(b, b)
	_reduce(b, b)

	_add(b, b)
	_reduce(b, b)

	_add(x, b, nbWords)
	_reduce(b, b)
	_add(x, b)
	_reduce(b, b)

	movq("res+0(FP)", x)
	_mov(a, x)
	_mov(b, x, 0, nbWords)

	ret()
}

func generateSquareE2BN256() {
	// var a, b fp.Element
	// a.Add(&x.A0, &x.A1)
	// b.Sub(&x.A0, &x.A1)
	// a.Mul(&a, &b)
	// b.Mul(&x.A0, &x.A1).Double(&b)
	// z.A0.Set(&a)
	// z.A1.Set(&b)
	fnHeader("squareAdx"+elementName, 16, 16, dx, ax)

	noAdx := newLabel()
	// check ADX instruction support
	cmpb("·supportAdx(SB)", 1)
	jne(noAdx)

	a := popRegisters(nbWords)
	b := popRegisters(nbWords)
	{
		x := popRegister()

		movq("x+8(FP)", x)
		_mov(x, a, nbWords) // a = a1
		_mov(x, b)          // b = a0

		// a = a0 + a1
		_add(b, a)
		_reduce(a, a)

		// b = a0 - a1
		_sub(x, b, nbWords)
		pushRegister(x)
		_reduceAfterSub(b, true)
	}

	// a = a * b
	{
		yat := func(i int) string {
			return string(b[i])
		}
		xat := func(i int) string {
			return string(a[i])
		}
		uglyHook := func(i int) {
			pushRegister(b[i])
		}
		t := mulAdx(yat, xat, uglyHook)
		_reduce(t, a)

		pushRegister(t...)
	}

	// b = a0 * a1 * 2
	{
		r := popRegister()
		movq("x+8(FP)", r)
		yat := func(i int) string {
			return r.at(i + nbWords)
		}
		xat := func(i int) string {
			return r.at(i)
		}
		b = mulAdx(yat, xat, nil)
		pushRegister(r)

		// reduce b
		_reduce(b, b)

		// double b (no reduction)
		_add(b, b)
	}

	// result.a1 = b
	r := popRegister()
	movq("res+0(FP)", r)
	_reduce(b, r, nbWords)

	// result.a0 = a
	_mov(a, r)

	ret()

	// No adx
	label(noAdx)
	movq("res+0(FP)", ax)
	movq(ax, "(SP)")
	movq("x+8(FP)", ax)
	movq(ax, "8(SP)")
	writeLn("CALL ·squareGenericE2(SB)")
	ret()
}

func generateMulE2BN256() {
	// var a, b, c fp.Element
	// a.Add(&x.A0, &x.A1)
	// b.Add(&y.A0, &y.A1)
	// a.Mul(&a, &b)
	// b.Mul(&x.A0, &y.A0)
	// c.Mul(&x.A1, &y.A1)
	// z.A1.Sub(&a, &b).Sub(&z.A1, &c)
	// z.A0.Sub(&b, &c)
	fnHeader("mulAdx"+elementName, 24, 24, dx, ax)

	noAdx := newLabel()
	// check ADX instruction support
	cmpb("·supportAdx(SB)", 1)
	jne(noAdx)

	a := popRegisters(nbWords)
	b := popRegisters(nbWords)
	{
		x := popRegister()

		movq("x+8(FP)", x)

		_mov(x, a, nbWords) // a = x.a1
		_add(x, a)          // a = x.a0 + x.a1
		_reduce(a, a)

		movq("y+16(FP)", x)
		_mov(x, b, nbWords) // b = y.a1
		_add(x, b)          // b = y.a0 + y.a1
		_reduce(b, b)

		pushRegister(x)
	}

	// a = a * b
	{
		yat := func(i int) string {
			return string(b[i])
		}
		xat := func(i int) string {
			return string(a[i])
		}
		uglyHook := func(i int) {
			pushRegister(b[i])
		}
		t := mulAdx(yat, xat, uglyHook)
		_reduce(t, a)

		pushRegister(t...)
	}

	// b = x.A0 * y.AO
	{
		r := popRegister()
		yat := func(i int) string {
			movq("y+16(FP)", r)
			return r.at(i)
		}
		xat := func(i int) string {
			movq("x+8(FP)", r)
			return r.at(i)
		}
		b = mulAdx(yat, xat, nil)
		pushRegister(r)
		_reduce(b, b)
	}
	// a - = b
	_sub(b, a)
	_reduceAfterSub(a, true)

	// push a to the stack for later use
	for i := 0; i < nbWords; i++ {
		pushq(a[i])
	}
	pushRegister(a...)

	var c []register
	// c = x.A1 * y.A1
	{
		r := popRegister()
		yat := func(i int) string {
			movq("y+16(FP)", r)
			return r.at(i + nbWords)
		}
		xat := func(i int) string {
			movq("x+8(FP)", r)
			return r.at(i + nbWords)
		}
		c = mulAdx(yat, xat, nil)
		pushRegister(r)
		_reduce(c, c)
	}

	// b = b - c
	_sub(c, b)
	_reduceAfterSub(b, true)

	// dereference result
	r := popRegister()
	movq("res+0(FP)", r)

	// z.A0 = b
	_mov(b, r)

	// restore a
	a = b
	for i := nbWords - 1; i >= 0; i-- {
		popq(a[i])
	}

	// a = a - c
	_sub(c, a)
	pushRegister(c...)

	// reduce a
	_reduceAfterSub(a, true)

	// z.A1 = a
	_mov(a, r, 0, nbWords)

	ret()

	// No adx
	label(noAdx)
	movq("res+0(FP)", ax)
	movq(ax, "(SP)")
	movq("x+8(FP)", ax)
	movq(ax, "8(SP)")
	movq("y+16(FP)", ax)
	movq(ax, "16(SP)")
	writeLn("CALL ·mulGenericE2(SB)")
	ret()
}
