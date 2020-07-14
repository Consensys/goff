package asm

import (
	"fmt"
	"os"

	"github.com/consensys/bavard"
)

const smallModulus = 6

func qAt(index int, elementName string) string {
	return fmt.Sprintf("·q%s+%d(SB)", elementName, index*8)
}

func qInv0(elementName string) string {
	return fmt.Sprintf("·q%sInv0(SB)", elementName)
}

type builder struct {
	path                      string
	elementName               string
	nbWords, nbWordsLastIndex int
	q                         []uint64
}

// NewBuilder returns a builder object to help generated assembly code for some operations
func NewBuilder(path, elementName string, nbWords int, q []uint64) *builder {
	return &builder{path, elementName, nbWords, nbWords - 1, q}
}

// Build ...
func (b *builder) Build(noCarrySquare bool) error {
	f, err := os.Create(b.path)
	if err != nil {
		return err
	}
	defer f.Close()
	asm := bavard.NewAssembly(f)
	asm.Write("#include \"textflag.h\"")

	if b.nbWords > smallModulus {
		// mul
		// fills up all available registers
		if err := b.mulLarge(asm); err != nil {
			return err
		}
	} else {
		// mul
		if err := b.mul(asm); err != nil {
			return err
		}
		// square
		if noCarrySquare {
			if err := b.square(asm); err != nil {
				return err
			}
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

	// double
	if err := b.double(asm); err != nil {
		return err
	}

	return nil
}
