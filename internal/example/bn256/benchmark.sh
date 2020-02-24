#!/bin/bash
go test -c
./bn256.test -test.run=NONE -test.bench="MulAssign" -test.count=5 -test.benchtime=2s -test.cpu=1 .  > bn256.txt 
benchstat -csv -sort=name -norange bn256.txt > bn256.csv 