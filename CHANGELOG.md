<a name="unreleased"></a>
## [Unreleased]

### Build
- added per branch workflow, faster on non dev/master
- tentatively, added windows target
- fix yml ci file, maybe.
- fix yml ci file, maybe.
- fix yml ci file, maybe.
- wip, shorter workflow on dev branches
- remove circleCI and try github slack integration
- revert version.go to correct value
- disable shallow clone in build workflow
- disable shallow clone in build workflow
- fix yml syntax
- testing error on version check in CI
- add codeql-analysis
- running go vet, gofmt check etc.
- run 32bit tests on example folder only, generating from 32bit machine not supported
- added no_adx and 32bits path in CI config file
- run CI on macOS-latest
- add mod cache + test with race option
- skipping version generated test while setting up github actions
- removing the generate step for now
- fix syntax err in go.yml
- added git dep install in go.yml
- experimenting with github actions
- experimenting with github actions

### Fix
- field generation for LexicographicallyLargest failed for some modulus


<a name="v0.3.10"></a>
## [v0.3.10] - 2021-02-01

<a name="v0.3.9"></a>
## [v0.3.9] - 2021-01-04
### Hotfix
- element.Neg(0) != 0 in some code path, fix wrong template in asm


<a name="v0.3.8"></a>
## [v0.3.8] - 2020-12-22

<a name="v0.3.7"></a>
## [v0.3.7] - 2020-12-16
### Circleci
- test all asm code path with noadx

### SetBytes
- uses sync.Pool to avoid allocating big.Int.

### Tests
- added generated fuzz-test methods for each API, factorizing templating around tests


<a name="v0.3.6"></a>
## [v0.3.6] - 2020-10-19

<a name="v0.3.5"></a>
## [v0.3.5] - 2020-09-28

<a name="v0.3.4"></a>
## [v0.3.4] - 2020-09-23

<a name="v0.3.3"></a>
## [v0.3.3] - 2020-09-04
### Circleci
- test only generated files on 32bit platform, not code generation itself
- testing 32bit compile path

### E2
- mul and square e2 now check for supportAdx and call generic version if not present, enabling wrapper inlining

### Element
- remove use of unsafe in ToBigInt method

### Wip
- cross platform builds -- 32bit compatibility
- cross platform builds -- added reduce to ops_generic.go for test to compile
- made golint and staticcheck happier. fixing bad cross platform build


<a name="v0.3.2"></a>
## [v0.3.2] - 2020-08-31
### Asm
- code cleanup. added (wip) some asm generators for tower of extension

### Circleci
- failling one test on purpose
- failling one test on purpose
- experimenting new workflow
- experimenting new workflow
- experimenting new workflow

### Element
- re-reverted IsZero to use OR instead of == . ASM code shows this generated version has no jump, better for speculative execution

### Pull Requests
- Merge pull request [#24](https://github.com/consensys/goff/issues/24) from ConsenSys/develop


<a name="v0.3.1"></a>
## [v0.3.1] - 2020-07-14
### Pull Requests
- Merge pull request [#23](https://github.com/consensys/goff/issues/23) from ConsenSys/develop


<a name="v0.3.0"></a>
## [v0.3.0] - 2020-06-22
### Cleanup
- remove asm double, mul fips, benches and other polluting templates

### Pull Requests
- Merge pull request [#20](https://github.com/consensys/goff/issues/20) from ConsenSys/develop
- Merge pull request [#19](https://github.com/consensys/goff/issues/19) from ConsenSys/handy_functions


<a name="v0.2.2"></a>
## [v0.2.2] - 2020-04-12
### Pull Requests
- Merge pull request [#18](https://github.com/consensys/goff/issues/18) from ConsenSys/develop


<a name="v0.2.1"></a>
## [v0.2.1] - 2020-04-08
### Pull Requests
- Merge pull request [#17](https://github.com/consensys/goff/issues/17) from ConsenSys/asm
- Merge pull request [#16](https://github.com/consensys/goff/issues/16) from ConsenSys/develop
- Merge pull request [#15](https://github.com/consensys/goff/issues/15) from ConsenSys/fast_conversion


<a name="v0.2.0"></a>
## [v0.2.0] - 2020-04-07
### Base
- added sync and strconv imports

### Circleci
- remove long test in benches (timeout)

### Element
- SetBigInt -- added fast path without memalloc when 0 < big.Int <= q
- added One() and FromInterface() methods

### Pull Requests
- Merge pull request [#14](https://github.com/consensys/goff/issues/14) from ConsenSys/develop
- Merge pull request [#13](https://github.com/consensys/goff/issues/13) from ConsenSys/asm


<a name="v0.1.0"></a>
## v0.1.0 - 2020-03-06
### Documentation
- highlighted lack of constant time operation in generated code

### Pull Requests
- Merge pull request [#2](https://github.com/consensys/goff/issues/2) from alexeykiselev/fix-set-random


[Unreleased]: https://github.com/consensys/goff/compare/v0.3.10...HEAD
[v0.3.10]: https://github.com/consensys/goff/compare/v0.3.9...v0.3.10
[v0.3.9]: https://github.com/consensys/goff/compare/v0.3.8...v0.3.9
[v0.3.8]: https://github.com/consensys/goff/compare/v0.3.7...v0.3.8
[v0.3.7]: https://github.com/consensys/goff/compare/v0.3.6...v0.3.7
[v0.3.6]: https://github.com/consensys/goff/compare/v0.3.5...v0.3.6
[v0.3.5]: https://github.com/consensys/goff/compare/v0.3.4...v0.3.5
[v0.3.4]: https://github.com/consensys/goff/compare/v0.3.3...v0.3.4
[v0.3.3]: https://github.com/consensys/goff/compare/v0.3.2...v0.3.3
[v0.3.2]: https://github.com/consensys/goff/compare/v0.3.1...v0.3.2
[v0.3.1]: https://github.com/consensys/goff/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/consensys/goff/compare/v0.2.2...v0.3.0
[v0.2.2]: https://github.com/consensys/goff/compare/v0.2.1...v0.2.2
[v0.2.1]: https://github.com/consensys/goff/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/consensys/goff/compare/v0.1.0...v0.2.0
