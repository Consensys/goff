package asm

import (
	"fmt"
	"os"

	"github.com/consensys/bavard"
)

const SmallModulus = 6

func qAt(index int, elementName string) string {
	return fmt.Sprintf("·q%s+%d(SB)", elementName, index*8)
}

func qInv0(elementName string) string {
	return fmt.Sprintf("·q%sInv0(SB)", elementName)
}

type Builder struct {
	path                      string
	elementName               string
	nbWords, nbWordsLastIndex int
	q                         []uint64
}

func NewBuilder(path, elementName string, nbWords int, q []uint64) *Builder {
	return &Builder{path, elementName, nbWords, nbWords - 1, q}
}

func (b *Builder) Build() error {
	f, err := os.Create(b.path)
	if err != nil {
		return err
	}
	defer f.Close()
	asm := bavard.NewAssembly(f)
	asm.Write("#include \"textflag.h\"")

	// mul assign
	if b.nbWords > 6 {
		if err := b.mulLarge(asm); err != nil {
			return err
		}
	} else {
		if err := b.mul(asm); err != nil {
			return err
		}
		// square
		if err := b.square(asm); err != nil {
			return err
		}
		// // sub
		if err := b.sub(asm); err != nil {
			return err
		}
	}

	// // from mont
	if err := b.fromMont(asm); err != nil {
		return err
	}

	// reduce
	if err := b.reduceFn(asm); err != nil {
		return err
	}

	// add
	if err := b.add(asm); err != nil {
		return err
	}

	return nil
}
