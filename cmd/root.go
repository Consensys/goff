// Copyright 2019 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"math/bits"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/consensys/bavard"
	"github.com/consensys/goff/internal/templates/asm"
	"github.com/consensys/goff/internal/templates/element"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "goff",
	Short:   "goff generates arithmetic operations for any moduli",
	Run:     cmdGenerate,
	Version: Version,
}

// flags
var (
	fModulus     string
	fOutputDir   string
	fPackageName string
	fElementName string
	fBenches     bool
)

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVarP(&fElementName, "element", "e", "", "name of the generated struct and file")
	rootCmd.PersistentFlags().StringVarP(&fModulus, "modulus", "m", "", "field modulus (base 10)")
	rootCmd.PersistentFlags().StringVarP(&fOutputDir, "output", "o", "", "destination path to create output files")
	rootCmd.PersistentFlags().StringVarP(&fPackageName, "package", "p", "", "package name in generated files")
	rootCmd.PersistentFlags().BoolVarP(&fBenches, "benches", "b", false, "set to true to generate montgomery multiplication (CIOS, FIPS, noCarry) benchmarks")

	if bits.UintSize != 64 {
		panic("goff only supports 64bits architectures")
	}
}

func cmdGenerate(cmd *cobra.Command, args []string) {
	fmt.Println()
	fmt.Println("running goff version", Version)
	fmt.Println()

	// parse flags
	if err := parseFlags(cmd); err != nil {
		_ = cmd.Usage()
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(-1)
	}

	// generate code
	if err := GenerateFF(fPackageName, fElementName, fModulus, fOutputDir, fBenches, false); err != nil {
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(-1)
	}
}

func GenerateFF(packageName, elementName, modulus, outputDir string, benches bool, noCollidingNames bool) error {
	// compute field constants
	F, err := newField(packageName, elementName, modulus, benches, noCollidingNames)
	if err != nil {
		return err
	}

	// source file templates
	src := []string{
		element.Base,
		element.Add,
		element.Sub,
		element.Reduce,
		element.Exp,
		element.FromMont,
		element.Conv,
		element.MulCIOS,
		element.MulFIPS,
		element.MulNoCarry,
		element.Sqrt,
	}

	// test file templates
	tst := []string{
		element.MulCIOS,
		element.MulFIPS,
		element.MulNoCarry,
		element.SquareNoCarryTemplate,
		element.Reduce,
		element.Test,
	}

	// output files
	eName := strings.ToLower(elementName)

	pathSrc := filepath.Join(outputDir, eName+".go")
	pathSrcArith := filepath.Join(outputDir, "arith.go")
	pathTest := filepath.Join(outputDir, eName+"_test.go")

	bavardOpts := []func(*bavard.Bavard) error{
		bavard.Apache2("ConsenSys AG", 2020),
		bavard.Package(F.PackageName, "contains field arithmetic operations"),
		bavard.GeneratedBy(fmt.Sprintf("goff (%s)", Version)),
	}

	// generate source file
	if err := bavard.Generate(pathSrc, src, F, bavardOpts...); err != nil {
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

	if F.ASM { // max words without having to deal with spilling
		// generate mul.s
		{
			asm := []string{
				asm.Mul,
			}
			pathMulAsm := filepath.Join(outputDir, eName+"_mul_amd64.s")
			if err := bavard.Generate(pathMulAsm, asm, F,
				bavard.GeneratedBy(fmt.Sprintf("goff (%s)", Version)),
				bavard.Import(false),
				bavard.Format(true),
				bavard.Funcs(template.FuncMap{
					"reg":    reg,
					"maxReg": maxReg,
				})); err != nil {
				return err
			}

			// generate mul_amd64.go
			src := []string{
				element.MontgomeryMultiplicationAMD64,
			}
			pathSrc := filepath.Join(outputDir, eName+"_mul_amd64.go")
			if err := bavard.Generate(pathSrc, src, F, bavardOpts...); err != nil {
				return err
			}
		}

		if F.NoCarrySquare {
			// generate square.s
			// TODO this doesn't work well  for now as it's slower than mul
			asm := []string{
				asm.Square,
			}
			pathMulAsm := filepath.Join(outputDir, eName+"_square_amd64.s")
			if err := bavard.Generate(pathMulAsm, asm, F,
				bavard.GeneratedBy(fmt.Sprintf("goff (%s)", Version)),
				bavard.Import(false),
				bavard.Format(true),
				bavard.Funcs(template.FuncMap{
					"reg":    reg,
					"maxReg": maxReg,
				})); err != nil {
				return err
			}

			// generate square_amd64.go
			src := []string{
				element.SquareCIOSNoCarryAMD64,
			}
			pathSrc := filepath.Join(outputDir, eName+"_square_amd64.go")
			if err := bavard.Generate(pathSrc, src, F, bavardOpts...); err != nil {
				return err
			}
		}

	}

	{
		// generate mul.go
		src := []string{
			element.MontgomeryMultiplication,
			element.MulCIOS,
			element.MulNoCarry,
			element.Reduce,
		}
		pathSrc := filepath.Join(outputDir, eName+"_mul.go")
		bavardOptsCpy := make([]func(*bavard.Bavard) error, len(bavardOpts))
		copy(bavardOptsCpy, bavardOpts)
		if F.ASM {
			bavardOptsCpy = append(bavardOptsCpy, bavard.BuildTag("!amd64"))
		}
		if err := bavard.Generate(pathSrc, src, F, bavardOptsCpy...); err != nil {
			return err
		}
	}
	{
		// generate square.go
		src := []string{
			element.SquareCIOSNoCarry,
			element.SquareNoCarryTemplate,
			element.MulNoCarry,
			element.Reduce,
		}
		pathSrc := filepath.Join(outputDir, eName+"_square.go")
		bavardOptsCpy := make([]func(*bavard.Bavard) error, len(bavardOpts))
		copy(bavardOptsCpy, bavardOpts)
		if F.ASM && F.NoCarrySquare {
			bavardOptsCpy = append(bavardOptsCpy, bavard.BuildTag("!amd64"))
		}
		if err := bavard.Generate(pathSrc, src, F, bavardOptsCpy...); err != nil {
			return err
		}
	}

	return nil
}

var registers = []string{ // AX and DX are reserved
	"CX",
	"BX",
	"BP",
	"SI",
	"DI",
	"R8",
	"R9",
	"R10",
	"R11",
	"R12",
	"R13",
	"R14",
	"R15",
}

func maxReg() int {
	return len(registers)
}

func reg(i int, offset ...int) string {
	if len(offset) == 1 {
		i += offset[0]
	}
	return registers[i]
}

func parseFlags(cmd *cobra.Command) error {
	if fModulus == "" ||
		fOutputDir == "" ||
		fPackageName == "" ||
		fElementName == "" {
		return errMissingArgument
	}

	// clean inputs
	fOutputDir = filepath.Clean(fOutputDir)
	fPackageName = strings.ToLower(fPackageName)

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
