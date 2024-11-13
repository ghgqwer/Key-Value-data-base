<strong>Benchmark storage results</strong>

BenchmarkGet - 0.0000005 ns/op        0 B/op    

BenchmarkSet - 0.0000448 ns/op        0 B/op    

BenchmarkSetGet - 0.0000242 ns/op        0 B/op 

AllBenchmarks - 0.0000336 ns/op        0 B/op  


<strong>Benchmark cleaner results</strong>
10000 iter: BenchmarkSingleClean-8   	22256586	        54.12 ns/op	       0 B/op	       0 allocs/op
10000 iter: BenchmarkMultiClean-8   	     478	   2279669 ns/op	  720185 B/op	   20001 allocs/op
1000000 iter: BenchmarkSingleClean-8   	20520859	        54.46 ns/op	       0 B/op	       0 allocs/op
1000000 iter: BenchmarkMultiClean-8   	       5	 225649667 ns/op	72000412 B/op	 2000002 allocs/op