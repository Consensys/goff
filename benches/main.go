package main

import (
	"bufio"
	"crypto/rand"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/consensys/goff/cmd"
)

//go:generate rm primes.csv
//go:generate go run main.go
func main() {
	parentDir := "./generated/"
	err := os.Mkdir(parentDir, 0700)
	if err != nil {
		log.Println(err)
	}
	primes, err := generatePrimes("primes.csv")
	if err != nil {
		log.Fatal(err)
	}
	packageName := "generated"
	for nbWords, prime := range primes {
		elementName := fmt.Sprintf("Element%02d", nbWords)
		if err := cmd.GenerateFF(packageName, elementName, prime.String(), parentDir, true); err != nil {
			log.Fatal(elementName, err)
		}
	}

}

func generatePrimes(path string) (map[int]*big.Int, error) {
	if _, err := os.Stat(path); err == nil {
		// file exists
		csvFile, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer csvFile.Close()
		reader := csv.NewReader(bufio.NewReader(csvFile))
		toReturn := make(map[int]*big.Int)
		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				return nil, err
			} else if len(line) != 2 {
				return nil, errors.New("expected 2 items per line")
			}
			nbWords, err := strconv.Atoi(strings.ToLower(strings.TrimSpace(line[0])))
			if err != nil {
				return nil, err
			}
			q, _ := big.NewInt(0).SetString(strings.TrimSpace(line[1]), 10)
			toReturn[nbWords] = q
		}
		return toReturn, nil
	}
	// file don't exists
	toReturn := make(map[int]*big.Int)
	csvFile, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	for i := 2; i < 33; i++ {
		j := (i * 64) - 1
		var q *big.Int
		q, _ = rand.Prime(rand.Reader, j)
		record := []string{fmt.Sprintf("%d", i), q.String()}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
		toReturn[i] = q
	}
	writer.Flush()
	return toReturn, nil
}
