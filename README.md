# Fast finite field arithmetic in Golang
[![Gitter](https://badges.gitter.im/consensys-gnark/community.svg)](https://gitter.im/consensys-gnark/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge) [![License](https://img.shields.io/badge/license-Apache%202-blue)](LICENSE)  [![Go Report Card](https://goreportcard.com/badge/github.com/consensys/goff)](https://goreportcard.com/badge/github.com/consensys/goff)

`goff` (go **f**inite **f**ield) is a fast field arithmetic library in Go.

We introduced `goff` [in this article](https://hackmd.io/@zkteam/goff): the project came from the need to have performant field operations in Go.
For most moduli, `goff` outperforms `math/big` and optimized libraries written in C++ or Rust.

In particular, `goff` modular multiplication is blazingly fast. ["Faster big-integer modular multiplication for most moduli"](https://hackmd.io/@zkteam/modular_multiplication) explains the algorithmic optimization we discovered and implemented, and presents some [benchmarks](https://github.com/ConsenSys/gurvy#benchmarks).

Actively developed and maintained by the team (zkteam@consensys.net) behind:
* [`gnark`: a framework to execute (and verify) algorithms in zero-knowledge](https://github.com/consensys/gnark) 
* [`gurvy`: elliptic curve cryptography (+pairing) library](https://github.com/consensys/gurvy)

## Warning
**`goff` has not been audited and is provided as-is, use at your own risk. In particular, `goff` makes no security guarantees such as constant time implementation or side-channel attack resistance.**

`goff` generates code optimized for 64bits architectures. It generates optimized assembly for moduli matching the `NoCarry` condition on `amd64` which support `ADX/MULX` instructions. Other targets have a fallback generic Go code. 

<img style="display: block;margin: auto;" width="80%"
src="banner_goff.png">

## Getting started

### Go version

`goff` is tested with the last 2 major releases of Go (1.14 and 1.15).

### Install `goff`

```bash
# dependencies
go get golang.org/x/tools/cmd/goimports

# goff
go get github.com/consensys/goff
```

## Usage

### Generated API

Example [API doc](https://pkg.go.dev/github.com/consensys/goff/examples/bn256)

### API -- go.mod (recommended)

At the root of your repo:
```bash
# goff
go get github.com/consensys/goff
``` 

then in a `main.go`  (that can be called using a `go:generate` workflow):

```
goff.GenerateFF(packageName, structName, modulus, destinationPath, false)
```


### Command line interface

```bash
goff

running goff version v0.3.4

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
goff -m 21888242871...94645226208583 -o ./bn256/ -p bn256 -e Element
```

outputs three `.go` files in `./bn256/`
* `element.go`
* `element_test.go`
* `arith.go`

The generated type has an API that's similar with `big.Int`

Example API signature
```go 
// Mul z = x * y mod q
func (z *Element) Mul(x, y *Element) *Element {
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

### Benchmarks


```bash
# for BN256 or BLS377
cd examples/bls377 # or cd examples/bn256

go test -c

./bls377.test -test.run=NONE -test.bench="." -test.count=10 -test.benchtime=1s -test.cpu=1 . | tee bls377.txt

benchstat bls377.txt 
```


## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our [code of conduct](CODE_OF_CONDUCT.md), and the process for submitting pull requests to us.
Get in touch: zkteam@consensys.net

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/consensys/goff/tags). 


## License

This project is licensed under the Apache 2 License - see the [LICENSE](LICENSE) file for details
