package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/consensys/bavard"
	"github.com/consensys/goff/asm/amd64"
	"github.com/consensys/goff/field"
	"github.com/consensys/goff/internal/templates/element"
)

// GenerateFF will generate go (and .s) files in outputDir for modulus (in base 10)
func GenerateFF(packageName, elementName, modulus, outputDir string, noCollidingNames bool) error {
	// compute field constants
	F, err := field.NewField(packageName, elementName, modulus, noCollidingNames, Version)
	if err != nil {
		return err
	}

	// source file templates
	src := []string{
		element.Base,
		element.Reduce,
		element.Exp,
		element.Conv,
		element.MulCIOS,
		element.MulNoCarry,
		element.Sqrt,
		element.Inverse,
	}

	// test file templates
	tst := []string{
		element.MulCIOS,
		element.MulNoCarry,
		element.Reduce,
		element.Test,
	}

	// output files
	eName := strings.ToLower(elementName)

	pathSrc := filepath.Join(outputDir, eName+".go")
	pathSrcArith := filepath.Join(outputDir, "arith.go")
	pathTest := filepath.Join(outputDir, eName+"_test.go")

	// remove old format generated files
	oldFiles := []string{"_mul.go", "_mul_amd64.go", "_mul_amd64.s",
		"_square.go", "_square_amd64.go", "_ops_decl.go", "_square_amd64.s", "_ops_amd64.go"}
	for _, of := range oldFiles {
		os.Remove(filepath.Join(outputDir, eName+of))
	}

	bavardOpts := []func(*bavard.Bavard) error{
		bavard.Apache2("ConsenSys Software Inc.", 2020),
		bavard.Package(F.PackageName),
		bavard.GeneratedBy(fmt.Sprintf("goff (%s)", Version)),
		bavard.Funcs(template.FuncMap{"toTitle": strings.Title}),
	}
	optsWithPackageDoc := append(bavardOpts, bavard.Package(F.PackageName, "contains field arithmetic operations for modulus "+F.Modulus))

	// generate source file
	if err := bavard.Generate(pathSrc, src, F, optsWithPackageDoc...); err != nil {
		return err
	}
	// generate arithmetics source file
	if err := bavard.Generate(pathSrcArith, []string{element.Arith}, F, bavardOpts...); err != nil {
		return err
	}

	// generate test file
	if err := bavard.Generate(pathTest, tst, F, bavardOpts...); err != nil {
		return err
	}

	// if we generate assembly code
	if F.ASM {
		// generate ops.s
		{
			pathSrc := filepath.Join(outputDir, eName+"_ops_amd64.s")
			f, err := os.Create(pathSrc)
			if err != nil {
				return err
			}
			defer f.Close()
			ffamd64 := amd64.NewFFAmd64(f, F)
			if err := ffamd64.Generate(); err != nil {
				return err
			}

		}

	}

	{
		// generate ops_amd64.go
		src := []string{
			element.OpsAMD64,
		}
		pathSrc := filepath.Join(outputDir, eName+"_ops_amd64.go")
		if err := bavard.Generate(pathSrc, src, F, bavardOpts...); err != nil {
			return err
		}
	}

	{
		// generate ops.go
		src := []string{
			element.OpsNoAsm,
			element.MulCIOS,
			element.MulNoCarry,
			element.Reduce,
		}
		pathSrc := filepath.Join(outputDir, eName+"_ops_noasm.go")
		bavardOptsCpy := make([]func(*bavard.Bavard) error, len(bavardOpts))
		copy(bavardOptsCpy, bavardOpts)
		if F.ASM {
			bavardOptsCpy = append(bavardOptsCpy, bavard.BuildTag("!amd64"))
		}
		if err := bavard.Generate(pathSrc, src, F, bavardOptsCpy...); err != nil {
			return err
		}
	}

	return nil
}
