# go-workers-pool

## tests

### For 100 000 units of products 
```
goos: linux
goarch: amd64
pkg: github.com/BogdanYanov/go-workers-pool/worker
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_1_-_1_worker-2                 1000000000               0.721 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_2_-_2_workers-2                1000000000               0.307 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_3_-_10_workers-2               1000000000               0.339 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_4_-_50_workers-2               1000000000               0.275 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_5_-_100_workers-2              1000000000               0.228 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_6_-_500_workers-2              1000000000               0.293 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_7_-_1000_workers-2             1000000000               0.285 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_8_-_5000_workers-2             1000000000               0.290 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_9_-_10000_workers-2            1000000000               0.299 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_10_-_50000_workers-2           1000000000               0.301 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_11_-_100000_workers-2          1000000000               0.294 ns/op
PASS
ok      github.com/BogdanYanov/go-workers-pool/worker   99.293s
```
### For 500 000 units of products
```
goos: linux
goarch: amd64
pkg: github.com/BogdanYanov/go-workers-pool/worker
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_1_-_1_worker-2                 100000000               28.2 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_2_-_2_workers-2                100000000               15.3 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_3_-_10_workers-2               100000000               13.2 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_4_-_50_workers-2               100000000               10.5 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_5_-_100_workers-2              100000000               10.2 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_6_-_500_workers-2              100000000               13.5 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_7_-_1000_workers-2             100000000               14.1 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_8_-_5000_workers-2             100000000               14.3 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_9_-_10000_workers-2            100000000               16.9 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_10_-_50000_workers-2           100000000               14.0 ns/op
BenchmarkWarehouse_Run/Benchmark_SendWork()_case_11_-_100000_workers-2          100000000               14.3 ns/op
PASS
ok      github.com/BogdanYanov/go-workers-pool/worker   19.732s
```
