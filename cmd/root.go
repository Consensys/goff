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

	"github.com/consensys/bavard"
	"github.com/consensys/goff/asm"
	"github.com/consensys/goff/internal/templates/e2"
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
)

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVarP(&fElementName, "element", "e", "", "name of the generated struct and file")
	rootCmd.PersistentFlags().StringVarP(&fModulus, "modulus", "m", "", "field modulus (base 10)")
	rootCmd.PersistentFlags().StringVarP(&fOutputDir, "output", "o", "", "destination path to create output files")
	rootCmd.PersistentFlags().StringVarP(&fPackageName, "package", "p", "", "package name in generated files")

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
	if err := GenerateFF(fPackageName, fElementName, fModulus, fOutputDir, false); err != nil {
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(-1)
	}
}

// GenerateFF2 will generate go (and .s) files in outputDir for E2 (field extension)
// modulus (in base 10)
func GenerateFF2(packageName, elementName, modulus, outputDir string) error {

	// compute field constants
	_F, err := newField(packageName, elementName, modulus, false)
	if err != nil {
		return err
	}

	type tData struct {
		*field
		BN256  bool
		BLS381 bool
	}
	F := &tData{field: _F}

	// TODO make this special curve business go away in gurvy.
	specialCurve := asm.NONE
	if modulus == "21888242871839275222246405745257275088696311157297823662689037894645226208583" {
		specialCurve = asm.BN256
		F.BN256 = true
	} else if modulus == "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787" {
		specialCurve = asm.BLS381
		F.BLS381 = true
	}

	// output files
	eName := strings.ToLower(elementName)

	pathSrc := filepath.Join(outputDir, eName+"_amd64.go")

	// source file templates
	src := []string{
		e2.Base,
	}

	bavardOpts := []func(*bavard.Bavard) error{
		bavard.Apache2("ConsenSys AG", 2020),
		bavard.Package(F.PackageName),
		bavard.GeneratedBy(fmt.Sprintf("goff (%s)", Version)),
	}

	// generate source file
	if err := bavard.Generate(pathSrc, src, F, bavardOpts...); err != nil {
		return err
	}

	// generate assembly
	{
		pathAsm := filepath.Join(outputDir, eName+"_amd64.s")
		f, err := os.Create(pathAsm)
		if err != nil {
			return err
		}
		defer f.Close()
		builder := asm.NewBuilder(f, F.ElementName, F.NbWords, F.Q, F.NoCarrySquare)

		if err := builder.GenerateTowerAssembly(specialCurve); err != nil {
			return err
		}

	}

	return nil
}

// GenerateFF will generate go (and .s) files in outputDir for modulus (in base 10)
func GenerateFF(packageName, elementName, modulus, outputDir string, noCollidingNames bool) error {
	// compute field constants
	F, err := newField(packageName, elementName, modulus, noCollidingNames)
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
		"_square.go", "_square_amd64.go", "_square_amd64.s", "_ops_amd64.go"}
	for _, of := range oldFiles {
		os.Remove(filepath.Join(outputDir, eName+of))
	}

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

	// if we generate assembly code
	if F.ASM {
		// generate ops.s
		{
			pathMulAsm := filepath.Join(outputDir, eName+"_ops_amd64.s")
			f, err := os.Create(pathMulAsm)
			if err != nil {
				return err
			}
			defer f.Close()
			builder := asm.NewBuilder(f, F.ElementName, F.NbWords, F.Q, F.NoCarrySquare)
			if err := builder.GenerateAssembly(); err != nil {
				return err
			}

		}

	}

	{
		// generate ops_decl.go
		src := []string{
			element.Ops,
			element.Reduce,
			element.MulCIOS,
			element.MulNoCarry,
		}
		pathSrc := filepath.Join(outputDir, eName+"_ops_decl.go")
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
