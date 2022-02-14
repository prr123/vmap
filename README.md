# vmap
faster golang map 

bench mark tests:
  + BenchmarkVmap-12        42433309                27.82 ns/op
  + BenchmarkGoMap-12       28941366                42.53 ns/op
  + BenchmarkVmapNH-12      79078124                18.53 ns/op
  + BenchmarkGoMapNH-12     54449022                27.59 ns/op

The bench marks of Vmap and GoMap are based on read operations with valid keys.
The bench marks of VmapNH and GoMapNH are based on read operations with invalid keys.
The tests were run with maps of 500 keys.
The speed improvement rests on using a hash function and a look-up table. Presumably the golang map fuction relies on iterating through the map to see whether a given key is in the map. The trade-off is a larger memory foot print. 
The hash function is based on DJB2 (for an explanation see for example: https://theartincode.stanis.me/008-djb2/).

The current implementation is based on uint16 (65535) table. This should allow up to 1000 keys, maybe even 10,000 keys with minimum collisions.
