# rasa/idbenchmark

Benchmark ID generators using Badger, Bolt, Bbolt, MySQL (ISAM & InnoDB) & Redis

## Getting Started

go get github.com/rasa/idbenchmark

## Running the tests

```shell
cd idbenchmark
make test
```

On a Windows 7 x64 box, Intel i7-3820QM @ 2.7Ghz with RAID5 SSDs, the `GOMAXPROCS=1` results are :
```
BenchmarkBadgerNoSync65536          	20000000	        57.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkBadger65536                	20000000	        58.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkBadgerNoSync4096           	20000000	        63.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkBadger4096                 	20000000	        63.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkBadger256                  	10000000	       142 ns/op	       8 B/op	       0 allocs/op
BenchmarkBadgerNoSync256            	10000000	       151 ns/op	       8 B/op	       0 allocs/op
BenchmarkBadger64                   	 5000000	       399 ns/op	      34 B/op	       1 allocs/op
BenchmarkBadgerNoSync64             	 5000000	       402 ns/op	      34 B/op	       1 allocs/op
BenchmarkBadger1                    	  100000	     21811 ns/op	    2180 B/op	      75 allocs/op
BenchmarkBadgerNoSync1              	  100000	     21931 ns/op	    2216 B/op	      75 allocs/op
BenchmarkRedis                      	   20000	     61003 ns/op	     360 B/op	      11 allocs/op
BenchmarkBbolt                      	   10000	    167509 ns/op	    6180 B/op	      40 allocs/op
BenchmarkBolt                       	   10000	    170709 ns/op	    6093 B/op	      37 allocs/op
BenchmarkMysqlInsert                	   10000	    173309 ns/op	      64 B/op	       3 allocs/op
BenchmarkMysqlReplace               	   10000	    177710 ns/op	      64 B/op	       3 allocs/op
BenchmarkMysqlUpdate                	   10000	    181610 ns/op	      64 B/op	       3 allocs/op
BenchmarkMysqlUpdateLimit1          	   10000	    182610 ns/op	      64 B/op	       3 allocs/op
BenchmarkInnoDBReplace              	    5000	    371821 ns/op	      65 B/op	       3 allocs/op
BenchmarkInnoDBInsert               	    3000	    395022 ns/op	      66 B/op	       3 allocs/op
BenchmarkInnoDBUpdate               	    5000	    396822 ns/op	      65 B/op	       3 allocs/op
BenchmarkBadgerParallelNoSync1000   	    5000	    414823 ns/op	   11077 B/op	     379 allocs/op
BenchmarkInnoDBUpdateLimit1         	    3000	    423690 ns/op	      66 B/op	       3 allocs/op
BenchmarkBadgerParallel100000       	   10000	    552131 ns/op	     220 B/op	       7 allocs/op
BenchmarkBadgerParallelNoSync10000  	   10000	    573932 ns/op	    2214 B/op	      75 allocs/op
BenchmarkBadgerParallelNoSync100000 	   10000	    585533 ns/op	     220 B/op	       7 allocs/op
BenchmarkBadgerParallel10000        	   10000	    605634 ns/op	    2214 B/op	      75 allocs/op
BenchmarkBadgerParallel1000         	   10000	    809446 ns/op	   22160 B/op	     759 allocs/op
BenchmarkBadgerParallel100          	    5000	   1394279 ns/op	  110804 B/op	    3800 allocs/op
BenchmarkBadgerParallelNoSync100    	    5000	   1409280 ns/op	  110804 B/op	    3800 allocs/op
BenchmarkBadgerParallelNoSync10     	    2000	   4537259 ns/op	  443220 B/op	   15200 allocs/op
BenchmarkBadgerParallelNoSync1      	     500	  11134637 ns/op	 1108042 B/op	   38000 allocs/op
BenchmarkBadgerParallel10           	    5000	  11852077 ns/op	 1207285 B/op	   38470 allocs/op
BenchmarkRedisParallel              	     200	  12755729 ns/op	   72000 B/op	    2200 allocs/op
BenchmarkBoltParallel               	     100	  16660953 ns/op	  609220 B/op	    3701 allocs/op
BenchmarkBboltParallel              	     100	  16690954 ns/op	  617854 B/op	    4000 allocs/op
BenchmarkMysqlInsertParallel        	     100	  18241043 ns/op	    6967 B/op	     307 allocs/op
BenchmarkMysqlReplaceParallel       	     100	  18281045 ns/op	    6967 B/op	     307 allocs/op
BenchmarkMysqlUpdateLimit1Parallel  	     100	  18311047 ns/op	    6967 B/op	     307 allocs/op
BenchmarkMysqlUpdateParallel        	     100	  19241101 ns/op	    6967 B/op	     307 allocs/op
BenchmarkBadgerParallel1            	    1000	  22103264 ns/op	 2216127 B/op	   76003 allocs/op
BenchmarkInnoDBReplaceParallel      	     100	  40122294 ns/op	    6967 B/op	     307 allocs/op
BenchmarkInnoDBUpdateLimit1Parallel 	     100	  41522375 ns/op	    6967 B/op	     307 allocs/op
BenchmarkInnoDBUpdateParallel       	     100	  42512431 ns/op	    6967 B/op	     307 allocs/op
BenchmarkInnoDBInsertParallel       	     200	  77854453 ns/op	   13331 B/op	     607 allocs/op
ok  	github.com/rasa/idbenchmark	275.639s
```
With `GOMAXPROCS=8`, the results are:
```
BenchmarkBadgerNoSync65536-8          	20000000	        61.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkBadger65536-8                 	20000000	        61.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkBadger4096-8                 	20000000	        66.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkBadgerNoSync4096-8            	20000000	        66.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkBadger256-8                  	10000000	       143 ns/op	       8 B/op	       0 allocs/op
BenchmarkBadgerNoSync256-8             	10000000	       144 ns/op	       8 B/op	       0 allocs/op
BenchmarkBadgerNoSync64-8             	 3000000	       403 ns/op	      33 B/op	       1 allocs/op
BenchmarkBadger64-8                    	 3000000	       410 ns/op	      33 B/op	       1 allocs/op
BenchmarkBadger1-8                    	   50000	     29121 ns/op	    2212 B/op	      75 allocs/op
BenchmarkBadgerNoSync1-8               	   50000	     29181 ns/op	    2213 B/op	      75 allocs/op
BenchmarkRedis-8                      	   20000	     60653 ns/op	     360 B/op	      11 allocs/op
BenchmarkBbolt-8                       	   10000	    172309 ns/op	    6186 B/op	      40 allocs/op
BenchmarkBolt-8                       	   10000	    172409 ns/op	    6097 B/op	      37 allocs/op
BenchmarkMysqlInsert-8                 	   10000	    182810 ns/op	      64 B/op	       3 allocs/op
BenchmarkMysqlUpdate-8                	   10000	    186810 ns/op	      64 B/op	       3 allocs/op
BenchmarkMysqlUpdateLimit1-8           	   10000	    191410 ns/op	      64 B/op	       3 allocs/op
BenchmarkMysqlReplace-8               	   10000	    233913 ns/op	      64 B/op	       3 allocs/op
BenchmarkInnoDBReplace-8               	    3000	    393689 ns/op	      66 B/op	       3 allocs/op
BenchmarkInnoDBInsert-8               	    5000	    393822 ns/op	      65 B/op	       3 allocs/op
BenchmarkInnoDBUpdate-8                	    3000	    416023 ns/op	      66 B/op	       3 allocs/op
BenchmarkInnoDBUpdateLimit1-8         	    3000	    417023 ns/op	      66 B/op	       3 allocs/op
BenchmarkBadgerParallel100-8           	    3000	    911718 ns/op	   63634 B/op	    2250 allocs/op
BenchmarkBadgerParallelNoSync100-8    	    3000	    919719 ns/op	   63631 B/op	    2250 allocs/op
BenchmarkBadgerParallel100000-8        	   10000	   1626093 ns/op	     211 B/op	       7 allocs/op
BenchmarkBadgerParallelNoSync100000-8 	   10000	   1632593 ns/op	     211 B/op	       7 allocs/op
BenchmarkBadgerParallel10000-8         	   10000	   1638793 ns/op	    2119 B/op	      74 allocs/op
BenchmarkBadgerParallelNoSync10000-8  	   10000	   1720898 ns/op	    2120 B/op	      74 allocs/op
BenchmarkBadgerParallel1000-8          	   10000	   1764200 ns/op	   21203 B/op	     750 allocs/op
BenchmarkBadgerParallelNoSync1000-8   	   10000	   1813503 ns/op	   21203 B/op	     750 allocs/op
BenchmarkBadgerParallel10-8            	    2000	   5657323 ns/op	  441685 B/op	   15184 allocs/op
BenchmarkBadgerParallelNoSync10-8     	    2000	   5672324 ns/op	  441395 B/op	   15181 allocs/op
BenchmarkRedisParallel-8               	     500	  11668667 ns/op	  180153 B/op	    5500 allocs/op
BenchmarkMysqlUpdateLimit1Parallel-8  	     100	  13590777 ns/op	    7481 B/op	     310 allocs/op
BenchmarkMysqlReplaceParallel-8        	     100	  13710784 ns/op	    7488 B/op	     310 allocs/op
BenchmarkBadgerParallelNoSync1-8      	     500	  13802789 ns/op	 1106501 B/op	   37984 allocs/op
BenchmarkBadgerParallel1-8             	     500	  14080805 ns/op	 1106239 B/op	   37981 allocs/op
BenchmarkMysqlInsertParallel-8        	     100	  14100806 ns/op	    7492 B/op	     310 allocs/op
BenchmarkMysqlUpdateParallel-8         	     100	  14100807 ns/op	    7506 B/op	     310 allocs/op
BenchmarkInnoDBReplaceParallel-8      	     100	  14630837 ns/op	   60894 B/op	     614 allocs/op
BenchmarkBoltParallel-8                	     100	  17541003 ns/op	  610020 B/op	    3701 allocs/op
BenchmarkBboltParallel-8              	     100	  17861021 ns/op	  618771 B/op	    4001 allocs/op
BenchmarkInnoDBInsertParallel-8        	     200	  27776588 ns/op	  128425 B/op	    1252 allocs/op
BenchmarkInnoDBUpdateLimit1Parallel-8 	     100	  36402082 ns/op	   10059 B/op	     337 allocs/op
BenchmarkInnoDBUpdateParallel-8        	     100	  38352194 ns/op	    9873 B/op	     334 allocs/op
ok  	github.com/rasa/idbenchmark	263.905s
```

## Authors

* **Ross Smith II** - *Initial work* - [@rasa](https://github.com/rasa)

See also the list of [contributors](https://github.com/rasa/idbenchmark/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
