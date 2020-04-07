package examples

//go:generate go run ../../../main.go -m 258664426012969094010652733694893533536393512754914660539884262666720468348340822774968888139573360124440321458177 -o ../../../examples/bls377/ -p bls377 -e Element
//go:generate go run ../../../main.go -m 21888242871839275222246405745257275088696311157297823662689037894645226208583 -o ../../../examples/bn256/ -p bn256 -e Element

//go:generate go test -short -count=1 ../../../examples/bn256/
//go:generate go test -short -count=1 ../../../examples/bls377/
