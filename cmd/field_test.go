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
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"testing"
)

// integration test will create modulus for various field sizes and run tests

func TestIntegration(t *testing.T) {
	parentDir := "./internal/tests/integration"
	os.RemoveAll(parentDir)
	err := os.Mkdir(parentDir, 0700)
	defer os.RemoveAll(parentDir)
	if err != nil {
		t.Fatal(err)
	}

	var bits []int
	if testing.Short() {
		for i := 128; i <= 448; i += 64 {
			bits = append(bits, i-3, i-2, i-1, i)
		}
	} else {
		for i := 120; i < 704; i++ {
			bits = append(bits, i)
		}
	}

	moduli := make(map[string]string)
	for _, i := range bits {
		var q *big.Int
		var nbWords int
		if i%64 == 0 {
			q, _ = rand.Prime(rand.Reader, i)
			moduli[fmt.Sprintf("e_cios_%04d", i)] = q.String()
		} else {
			for {
				q, _ = rand.Prime(rand.Reader, i)
				nbWords = len(q.Bits())
				const B = (^uint64(0) >> 1) - 1
				// TODO platform specific here
				if uint64(q.Bits()[nbWords-1]) <= B {
					break
				}
			}
			moduli[fmt.Sprintf("e_nocarry_%04d", i)] = q.String()
		}

	}

	for elementName, modulus := range moduli {
		// generate field
		if err := GenerateFF("integration", elementName, modulus, parentDir, false); err != nil {
			t.Fatal(elementName, err)
		}
	}

	// run go test
	cmd := exec.Command("go", "test", "./"+parentDir)
	out, err := cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		t.Fatal(err)
	}

}
