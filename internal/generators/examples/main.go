package main

import (
	"path/filepath"
	"sync"

	"github.com/consensys/goff/field"
	"github.com/consensys/goff/generator"
)

//go:generate go run main.go
//go:generate go test -short -count=1 ../../../examples/...
func main() {
	var wg sync.WaitGroup
	for _, fData := range []struct {
		modulus string
		label   string
	}{
		{label: "bn256", modulus: "21888242871839275222246405745257275088696311157297823662689037894645226208583"},
		{label: "bls381", modulus: "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787"},
		{label: "bls377", modulus: "258664426012969094010652733694893533536393512754914660539884262666720468348340822774968888139573360124440321458177"},
		{label: "bw761", modulus: "6891450384315732539396789682275657542479668912536150109513790160209623422243491736087683183289411687640864567753786613451161759120554247759349511699125301598951605099378508850372543631423596795951899700429969112842764913119068299"},
	} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dir := filepath.Join("../../../examples/", fData.label)
			F, _ := field.NewField("fp", "Element", fData.modulus)
			if err := generator.GenerateFF(F, dir); err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
}
