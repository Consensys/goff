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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/consensys/goff/cmd/internal/template/element"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "goff",
	Short:   "goff generates arithmetic operations for any moduli",
	Run:     cmdGenerate,
	Version: buildString(),
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
}

func cmdGenerate(cmd *cobra.Command, args []string) {
	fmt.Println()
	if Version != "" {
		fmt.Println("running goff version", Version)
	} else {
		fmt.Println("/!\\ running goff in DEV mode /!\\")
	}

	fmt.Println()

	// parse flags
	if err := parseFlags(cmd); err != nil {
		_ = cmd.Usage()
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(-1)
	}

	// generate code
	if err := GenerateFF(fPackageName, fElementName, fModulus, fOutputDir, fBenches); err != nil {
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(-1)
	}
}

func GenerateFF(packageName, elementName, modulus, outputDir string, benches bool) error {
	// compute field constants
	F, err := newField(packageName, elementName, modulus, benches)
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
		element.MontgomeryMultiplication,
		element.Sqrt,
	}

	if F.NoCarrySquare {
		src = append(src, element.SquareCIOSNoCarry)
	} else {
		src = append(src, element.MontSquareCIOS)
	}

	// test file templates
	tst := []string{
		element.MulCIOS,
		element.MulFIPS,
		element.MulNoCarry,
		element.Reduce,
		element.Test,
	}

	// output files
	eName := strings.ToLower(elementName)

	pathSrc := filepath.Join(outputDir, eName+".go")
	pathSrcArith := filepath.Join(outputDir, "arith.go")
	pathTest := filepath.Join(outputDir, eName+"_test.go")

	// generate source file
	if err := generateCode(pathSrc, src, F); err != nil {
		return err
	}
	// generate arithmetics source file
	if err := generateCode(pathSrcArith, []string{element.Arith}, F); err != nil {
		return err
	}

	// generate test file
	if err := generateCode(pathTest, tst, F); err != nil {
		return err
	}

	return nil
}

func generateCode(output string, templates []string, F *field) error {
	// create output file
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	fmt.Printf("generating %-70s\n", output)

	// parse templates
	tmpl := template.Must(template.New("").
		Funcs(helpers()).
		Parse(aggregate(templates)))

	// execute template
	if err = tmpl.Execute(file, F); err != nil {
		file.Close()
		return err
	}
	file.Close()

	// run goformat to prettify output source
	if err := exec.Command("gofmt", "-s", "-w", output).Run(); err != nil {
		return err
	}
	if err := exec.Command("goimports", "-w", output).Run(); err != nil {
		return err
	}
	return nil
}

func aggregate(values []string) string {
	var sb strings.Builder
	for _, v := range values {
		sb.WriteString(v)
	}
	return sb.String()
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
