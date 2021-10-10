# Wrapper library for CGo and gRPC calls to [dusk-bls12_381-sign](https://github.com/dusk-network/bls12_381-sign)

This library provides wrapper functions which make CGo calls to the BLS 
signature crate linked above, as well as instead running it as a service and 
directing calls through this route.

## Building

To obtain the `dusk-bls12_381-sign` library, run:

```
make build
```

This will create the static C library for you and link it with the Go binary 
via cgo, as well as the microservice binary, and last but not least, 
generates the Go grpc/protobuf which is placed in `bls/proto/bls12381sig.pb.go`.

The generated code is used by the functions in ipc*.go to encode requests 
and decode responses from the Rust gRPC protobuf service that runs to 
provide the API in this mode.

## Testing

```
make test
```

## Benchmarks

### Running

```
make bench
```

### Machine specs

The benchmarks were ran on a 2020 13.3" MacBook Pro.

CPU:

```
$ lscpu
Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
```

RAM:

```
16 GB 3733 MHz LPDDR4X
```

### Results

```
BenchmarkSign
BenchmarkSign-8                      361           3368481 ns/op
BenchmarkVerify
BenchmarkVerify-8                    171           7037573 ns/op
BenchmarkAggregatePk
BenchmarkAggregatePk-8               228           5301390 ns/op
BenchmarkAggregateSig
BenchmarkAggregateSig-8             1123           1063291 ns/op
```
