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

// Package goff (go finite field) is a unix-like tool that generates fast field arithmetic in Go.
package main

import "github.com/consensys/goff/cmd"

//go:generate go run main.go -m 258664426012969094010652733694893533536393512754914660539884262666720468348340822774968888139573360124440321458177 -o ./internal/example/bls377/ -p bls377 -e Element
//go:generate go run main.go -m 21888242871839275222246405745257275088696311157297823662689037894645226208583 -o ./internal/example/bn256/ -p bn256 -e Element

//go:generate go test -short -count=1 ./internal/example/bn256/
//go:generate go test -short -count=1 ./internal/example/bls377/

func main() {
	cmd.Execute()
}
