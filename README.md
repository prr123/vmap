# vmap
faster golang map 

bench mark tests:
BenchmarkVmap-12        42433309                27.82 ns/op
BenchmarkGoMap-12       28941366                42.53 ns/op
BenchmarkVmapNH-12      79078124                18.53 ns/op
BenchmarkGoMapNH-12     54449022                27.59 ns/op

The bech marks of Vmap and GoMap are based on read operations with valid keys.
VmapNH and GoMapNH are based on read operations with invalid keys.
