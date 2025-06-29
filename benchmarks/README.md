### Benchmarks against `sync.WaitGroup` (in `group_WaitGroup_bench_test.go`)

```
goos: darwin
goarch: arm64
pkg: github.com/asmsh/sema/benchmarks
cpu: Apple M2
BenchmarkWaitGroupUncontended
BenchmarkWaitGroupUncontended/Uncontended
BenchmarkWaitGroupUncontended/Uncontended-8         	444735992	         2.664 ns/op	       0 B/op	       0 allocs/op
BenchmarkWaitGroupUncontended/Uncontended-HighParallelism
BenchmarkWaitGroupUncontended/Uncontended-HighParallelism-8         	462382617	         2.561 ns/op	       0 B/op	       0 allocs/op
BenchmarkSemaGroupUncontended
BenchmarkSemaGroupUncontended/Uncontended
BenchmarkSemaGroupUncontended/Uncontended-8                         	27039561	        44.51 ns/op	       0 B/op	       0 allocs/op
BenchmarkSemaGroupUncontended/Uncontended-HighParallelism
BenchmarkSemaGroupUncontended/Uncontended-HighParallelism-8         	26653957	        44.75 ns/op	       0 B/op	       0 allocs/op
BenchmarkWaitGroupAddDone
BenchmarkWaitGroupAddDone/no_work
BenchmarkWaitGroupAddDone/no_work-8                                 	17774079	        67.54 ns/op	       0 B/op	       0 allocs/op
BenchmarkWaitGroupAddDone/with_work
BenchmarkWaitGroupAddDone/with_work-8                               	13681306	        86.09 ns/op	       0 B/op	       0 allocs/op
BenchmarkWaitGroupAddDone/with_work-HighParallelism
BenchmarkWaitGroupAddDone/with_work-HighParallelism-8               	11923569	        84.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkSemaGroupAddDone
BenchmarkSemaGroupAddDone/no_work
BenchmarkSemaGroupAddDone/no_work-8                                 	 4933724	       245.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkSemaGroupAddDone/with_work
BenchmarkSemaGroupAddDone/with_work-8                               	 4297702	       286.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkSemaGroupAddDone/with_work-HighParallelism
BenchmarkSemaGroupAddDone/with_work-HighParallelism-8               	 4066167	       292.5 ns/op	       3 B/op	       0 allocs/op
BenchmarkWaitGroupWait
BenchmarkWaitGroupWait/no_work
BenchmarkWaitGroupWait/no_work-8                                    	1000000000	         0.4215 ns/op	       0 B/op	       0 allocs/op
BenchmarkWaitGroupWait/with_work
BenchmarkWaitGroupWait/with_work-8                                  	185626333	         6.062 ns/op	       0 B/op	       0 allocs/op
BenchmarkWaitGroupWait/with_work-HighParallelism
BenchmarkWaitGroupWait/with_work-HighParallelism-8                  	188749502	         6.254 ns/op	       0 B/op	       0 allocs/op
BenchmarkSemaGroupWait
BenchmarkSemaGroupWait/no_work
BenchmarkSemaGroupWait/no_work-8                                    	47804797	        23.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkSemaGroupWait/with_work
BenchmarkSemaGroupWait/with_work-8                                  	24355755	        49.41 ns/op	       0 B/op	       0 allocs/op
BenchmarkSemaGroupWait/with_work-HighParallelism
BenchmarkSemaGroupWait/with_work-HighParallelism-8                  	23944462	        49.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkWaitGroupActuallyWait
BenchmarkWaitGroupActuallyWait/#00
BenchmarkWaitGroupActuallyWait/#00-8                                	10681387	       115.0 ns/op	      32 B/op	       2 allocs/op
BenchmarkWaitGroupActuallyWait/HighParallelism
BenchmarkWaitGroupActuallyWait/HighParallelism-8                    	 8999552	       122.8 ns/op	      32 B/op	       2 allocs/op
BenchmarkSemaGroupActuallyWait
BenchmarkSemaGroupActuallyWait/#00
BenchmarkSemaGroupActuallyWait/#00-8                                	 7234860	       187.4 ns/op	     175 B/op	       2 allocs/op
BenchmarkSemaGroupActuallyWait/HighParallelism
BenchmarkSemaGroupActuallyWait/HighParallelism-8                    	 9562087	       128.6 ns/op	     176 B/op	       3 allocs/op
PASS
```

> Note: the `HighParallelism` version of the tests calls `*testing.B.SetParallelism` with a big number.

### Benchmarks against `semaphore.Weighted` (in `group_Weighted_bench_test.go`)

```
goos: darwin
goarch: arm64
pkg: github.com/asmsh/sema/benchmarks
cpu: Apple M2
BenchmarkNewSeq
BenchmarkNewSeq/Weighted-1
BenchmarkNewSeq/Weighted-1-8         	64114237	        18.62 ns/op	      80 B/op	       1 allocs/op
BenchmarkNewSeq/semChan-1
BenchmarkNewSeq/semChan-1-8          	58379600	        20.32 ns/op	     112 B/op	       1 allocs/op
BenchmarkNewSeq/sema.Group-1
BenchmarkNewSeq/sema.Group-1-8       	28960964	        43.07 ns/op	     160 B/op	       2 allocs/op
BenchmarkNewSeq/Weighted-128
BenchmarkNewSeq/Weighted-128-8       	64233522	        18.79 ns/op	      80 B/op	       1 allocs/op
BenchmarkNewSeq/semChan-128
BenchmarkNewSeq/semChan-128-8        	62760547	        20.75 ns/op	     112 B/op	       1 allocs/op
BenchmarkNewSeq/sema.Group-128
BenchmarkNewSeq/sema.Group-128-8     	26161530	        43.46 ns/op	     160 B/op	       2 allocs/op
BenchmarkAcquireSeq
BenchmarkAcquireSeq/Weighted-acquire-1-1-1
BenchmarkAcquireSeq/Weighted-acquire-1-1-1-8         	56787709	        20.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-tryAcquire-1-1-1
BenchmarkAcquireSeq/Weighted-tryAcquire-1-1-1-8      	70184672	        17.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-acquire-1-1-1
BenchmarkAcquireSeq/semChan-acquire-1-1-1-8          	44740088	        26.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-tryAcquire-1-1-1
BenchmarkAcquireSeq/semChan-tryAcquire-1-1-1-8       	43468286	        27.19 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-acquire-1-1-1
BenchmarkAcquireSeq/sema.Group-acquire-1-1-1-8       	60758466	        19.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-tryAcquire-1-1-1
BenchmarkAcquireSeq/sema.Group-tryAcquire-1-1-1-8    	66868198	        18.27 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-acquire-2-1-1
BenchmarkAcquireSeq/Weighted-acquire-2-1-1-8         	57632158	        20.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-tryAcquire-2-1-1
BenchmarkAcquireSeq/Weighted-tryAcquire-2-1-1-8      	71248684	        16.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-acquire-2-1-1
BenchmarkAcquireSeq/semChan-acquire-2-1-1-8          	42636093	        27.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-tryAcquire-2-1-1
BenchmarkAcquireSeq/semChan-tryAcquire-2-1-1-8       	41260804	        29.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-acquire-2-1-1
BenchmarkAcquireSeq/sema.Group-acquire-2-1-1-8       	59411086	        19.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-tryAcquire-2-1-1
BenchmarkAcquireSeq/sema.Group-tryAcquire-2-1-1-8    	64465298	        18.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-acquire-16-1-1
BenchmarkAcquireSeq/Weighted-acquire-16-1-1-8        	56219694	        21.19 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-tryAcquire-16-1-1
BenchmarkAcquireSeq/Weighted-tryAcquire-16-1-1-8     	70196648	        17.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-acquire-16-1-1
BenchmarkAcquireSeq/semChan-acquire-16-1-1-8         	40661769	        31.21 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-tryAcquire-16-1-1
BenchmarkAcquireSeq/semChan-tryAcquire-16-1-1-8      	39260431	        30.64 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-acquire-16-1-1
BenchmarkAcquireSeq/sema.Group-acquire-16-1-1-8      	59444317	        20.21 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-tryAcquire-16-1-1
BenchmarkAcquireSeq/sema.Group-tryAcquire-16-1-1-8   	63143910	        18.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-acquire-128-1-1
BenchmarkAcquireSeq/Weighted-acquire-128-1-1-8       	58138008	        20.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-tryAcquire-128-1-1
BenchmarkAcquireSeq/Weighted-tryAcquire-128-1-1-8    	69257736	        17.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-acquire-128-1-1
BenchmarkAcquireSeq/semChan-acquire-128-1-1-8        	40666765	        29.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-tryAcquire-128-1-1
BenchmarkAcquireSeq/semChan-tryAcquire-128-1-1-8     	40026404	        29.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-acquire-128-1-1
BenchmarkAcquireSeq/sema.Group-acquire-128-1-1-8     	62855058	        18.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-tryAcquire-128-1-1
BenchmarkAcquireSeq/sema.Group-tryAcquire-128-1-1-8  	63912934	        17.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-acquire-2-2-1
BenchmarkAcquireSeq/Weighted-acquire-2-2-1-8         	55351615	        21.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-tryAcquire-2-2-1
BenchmarkAcquireSeq/Weighted-tryAcquire-2-2-1-8      	68423340	        17.02 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-acquire-2-2-1
BenchmarkAcquireSeq/semChan-acquire-2-2-1-8          	21834556	        54.08 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-tryAcquire-2-2-1
BenchmarkAcquireSeq/semChan-tryAcquire-2-2-1-8       	21684639	        54.95 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-acquire-2-2-1
BenchmarkAcquireSeq/sema.Group-acquire-2-2-1-8       	61023801	        19.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-tryAcquire-2-2-1
BenchmarkAcquireSeq/sema.Group-tryAcquire-2-2-1-8    	64455343	        18.52 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-acquire-16-2-8
BenchmarkAcquireSeq/Weighted-acquire-16-2-8-8        	 7945699	       150.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-tryAcquire-16-2-8
BenchmarkAcquireSeq/Weighted-tryAcquire-16-2-8-8     	 9197974	       129.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-acquire-16-2-8
BenchmarkAcquireSeq/semChan-acquire-16-2-8-8         	 2693203	       443.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-tryAcquire-16-2-8
BenchmarkAcquireSeq/semChan-tryAcquire-16-2-8-8      	 2647851	       454.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-acquire-16-2-8
BenchmarkAcquireSeq/sema.Group-acquire-16-2-8-8      	 8075656	       148.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-tryAcquire-16-2-8
BenchmarkAcquireSeq/sema.Group-tryAcquire-16-2-8-8   	 7835320	       149.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-acquire-128-2-64
BenchmarkAcquireSeq/Weighted-acquire-128-2-64-8      	  978645	      1223 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-tryAcquire-128-2-64
BenchmarkAcquireSeq/Weighted-tryAcquire-128-2-64-8   	 1000000	      1027 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-acquire-128-2-64
BenchmarkAcquireSeq/semChan-acquire-128-2-64-8       	  328615	      3594 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-tryAcquire-128-2-64
BenchmarkAcquireSeq/semChan-tryAcquire-128-2-64-8    	  319842	      3676 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-acquire-128-2-64
BenchmarkAcquireSeq/sema.Group-acquire-128-2-64-8    	  994155	      1175 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-tryAcquire-128-2-64
BenchmarkAcquireSeq/sema.Group-tryAcquire-128-2-64-8 	  997916	      1178 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-acquire-2-1-2
BenchmarkAcquireSeq/Weighted-acquire-2-1-2-8         	30426123	        39.20 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-tryAcquire-2-1-2
BenchmarkAcquireSeq/Weighted-tryAcquire-2-1-2-8      	36450252	        33.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-acquire-2-1-2
BenchmarkAcquireSeq/semChan-acquire-2-1-2-8          	21353325	        56.13 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-tryAcquire-2-1-2
BenchmarkAcquireSeq/semChan-tryAcquire-2-1-2-8       	20759044	        57.29 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-acquire-2-1-2
BenchmarkAcquireSeq/sema.Group-acquire-2-1-2-8       	36670378	        32.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-tryAcquire-2-1-2
BenchmarkAcquireSeq/sema.Group-tryAcquire-2-1-2-8    	36588200	        32.88 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-acquire-16-8-2
BenchmarkAcquireSeq/Weighted-acquire-16-8-2-8        	30643285	        39.31 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-tryAcquire-16-8-2
BenchmarkAcquireSeq/Weighted-tryAcquire-16-8-2-8     	36339043	        33.06 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-acquire-16-8-2
BenchmarkAcquireSeq/semChan-acquire-16-8-2-8         	 2704630	       441.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-tryAcquire-16-8-2
BenchmarkAcquireSeq/semChan-tryAcquire-16-8-2-8      	 2694669	       443.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-acquire-16-8-2
BenchmarkAcquireSeq/sema.Group-acquire-16-8-2-8      	36998781	        32.04 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-tryAcquire-16-8-2
BenchmarkAcquireSeq/sema.Group-tryAcquire-16-8-2-8   	36525048	        32.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-acquire-128-64-2
BenchmarkAcquireSeq/Weighted-acquire-128-64-2-8      	30605348	        39.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/Weighted-tryAcquire-128-64-2
BenchmarkAcquireSeq/Weighted-tryAcquire-128-64-2-8   	36082899	        33.08 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-acquire-128-64-2
BenchmarkAcquireSeq/semChan-acquire-128-64-2-8       	  338746	      3541 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/semChan-tryAcquire-128-64-2
BenchmarkAcquireSeq/semChan-tryAcquire-128-64-2-8    	  336674	      3550 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-acquire-128-64-2
BenchmarkAcquireSeq/sema.Group-acquire-128-64-2-8    	37181728	        32.32 ns/op	       0 B/op	       0 allocs/op
BenchmarkAcquireSeq/sema.Group-tryAcquire-128-64-2
BenchmarkAcquireSeq/sema.Group-tryAcquire-128-64-2-8 	37119571	        32.74 ns/op	       0 B/op	       0 allocs/op
PASS
```
