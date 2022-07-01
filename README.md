
----------------------------------------

# **`goff` migrated to [`gnark-crypto`](https://github.com/ConsenSys/gnark-crypto) (/field/goff)**


----------------------------------------


## Fast finite field arithmetic in Golang
[![License](https://img.shields.io/badge/license-Apache%202-blue)](LICENSE)  [![Go Report Card](https://goreportcard.com/badge/github.com/consensys/goff)](https://goreportcard.com/badge/github.com/consensys/goff)

`goff` (go **f**inite **f**ield) is a unix-like tool that generates fast field arithmetic in Go.

We introduced `goff` [in this article](https://hackmd.io/@gnark/goff): the project came from the need to have performant field operations in Go.
For most moduli, `goff` outperforms `math/big` and optimized libraries written in C++ or Rust.

In particular, `goff` modular multiplication is blazingly fast. ["Faster big-integer modular multiplication for most moduli"](https://hackmd.io/@gnark/modular_multiplication) explains the algorithmic optimization we discovered and implemented, and presents some [benchmarks](https://github.com/ConsenSys/gnark-crypto#benchmarks).

Actively developed and maintained by the team (gnark@consensys.net) behind:
* [`gnark`: a framework to execute (and verify) algorithms in zero-knowledge](https://github.com/consensys/gnark) 

## Warning
**`goff` has not been audited and is provided as-is, use at your own risk. In particular, `goff` makes no security guarantees such as constant time implementation or side-channel attack resistance.**

`goff` generates code optimized for 64bits architectures. It generates optimized assembly for moduli matching the `NoCarry` condition on `amd64` which support `ADX/MULX` instructions. Other targets have a fallback generic Go code. 

__Since v0.4.0, `goff`'s code has been migrated into [`gnark-crypto`](https://github.com/consensys/gnark-crypto). This repo contains the unix-like tool only.__

<img style="display: block;margin: auto;" width="80%"
src="banner_goff.png">

## Getting started

### Go version

`goff` is tested with the last 2 major releases of Go (1.15 and 1.16).

### Install `goff`

```bash
# dependencies
go get golang.org/x/tools/cmd/goimports
go get github.com/klauspost/asmfmt/cmd/asmfmt

# goff
go get github.com/consensys/goff
```

## Usage

### Generated API

Example [API doc](https://pkg.go.dev/github.com/consensys/gnark-crypto/ecc/bn254/fr)

### API -- go.mod (recommended)

At the root of your repo:
```bash
# note that code has been migrated in gnark-crypto since v0.4.0
go get github.com/consensys/gnark-crypto
``` 

then in a `main.go`  (that can be called using a `go:generate` workflow):

```go
import (
  "github.com/consensys/gnark-crypto/field/generator"
  "github.com/consensys/gnark-crypto/field"

fp, _ = field.NewField("fp", "Element", fpModulus)

generator.GenerateFF(fp, "fp"))
```


### Command line interface

```bash
goff

running goff version v0.4.0

Usage:
  goff [flags]

Flags:
  -e, --element string   name of the generated struct and file
  -h, --help             help for goff
  -m, --modulus string   field modulus (base 10)
  -o, --output string    destination path to create output files
  -p, --package string   package name in generated files
  -v, --version          version for goff
```

### `goff` -- a short example

Running 
```bash
goff -m 21888242871946452262085832188824287194645226208583 -o ./bn256/ -p bn256 -e Element
```

outputs the `.go` and `.s` files in `./bn256/`



The generated type has an API that's similar with `big.Int`

Example API signature
```go 
// Mul z = x * y mod q
func (z *Element) Mul(x, y *Element) *Element 
```

and can be used like so:

```go 
var a, b Element
a.SetUint64(2)
b.SetString("984896738")

a.Mul(a, b)

a.Sub(a, a)
 .Add(a, b)
 .Inv(a)
 
b.Exp(b, 42)
b.Neg(b)
```

### Build tags

`goff` generate optimized assembly for `amd64` target. 

For the `Mul` operation, using `ADX` instructions and `ADOX/ADCX` result in a significant performance gain. 

The "default" target `amd64` checks if the running architecture supports these instruction, and reverts to generic path if not. This check adds a branch and forces the function to reserve some bytes on the frame to store the argument to call `_mulGeneric` .

`goff` output can be compiled with `amd64_adx` flag which omits this check. Will crash if the platform running the binary doesn't support the `ADX` instructions (roughly, before 2016). 


## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our [code of conduct](CODE_OF_CONDUCT.md), and the process for submitting pull requests to us.
Get in touch: zkteam@consensys.net

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/consensys/goff/tags). 


## License

This project is licensed under the Apache 2 License - see the [LICENSE](LICENSE) file for details
