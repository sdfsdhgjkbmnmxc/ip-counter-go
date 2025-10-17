# IP Address Counter

![Benchmarks](https://github.com/sdfsdhgjkbmnmxc/ip-counter-go/workflows/Benchmarks/badge.svg)

Solution for [Ecwid IP Address Counter Task](https://github.com/Ecwid/new-job/blob/master/IP-Addr-Counter-GO.md)

Download test file:

```bash
make testdata-large
```

Run benchmarks:

```bash
make bench
```

### Benchmarks

```
cpu: Apple M1 Pro
BenchmarkCounter/ips_1k.txt/naive_(init_buf=65536,_max_buf=4096,_cap=0)-10         	   55280	    108543 ns/op	  190433 B/op	    1024 allocs/op
BenchmarkCounter/ips_1k.txt/naive_(init_buf=65536,_max_buf=16,_cap=1024)-10        	  100873	     59255 ns/op	  136280 B/op	    1009 allocs/op
BenchmarkCounter/ips_1k.txt/netip_map-10                                           	   47198	    128312 ns/op	  180225 B/op	    1024 allocs/op
BenchmarkCounter/ips_1k.txt/uint32_map-10                                          	   67312	     87777 ns/op	   42288 B/op	      24 allocs/op
BenchmarkCounter/ips_1k.txt/uint32_mmap-10                                         	  124393	     47692 ns/op	   19368 B/op	       9 allocs/op
BenchmarkCounter/ips_1k.txt/uint32_mmap_parallel(workers=2)-10                     	   84382	     71482 ns/op	   38809 B/op	      24 allocs/op
BenchmarkCounter/ips_1k.txt/uint32_mmap_parallel(workers=8)-10                     	   78555	     77599 ns/op	   39737 B/op	      60 allocs/op
BenchmarkCounter/ips_1k.txt/bitmap-10                                              	    1173	   5098662 ns/op	536875157 B/op	       5 allocs/op
BenchmarkCounter/ips_1k.txt/bash_sort-10                                           	    1150	   5149207 ns/op	   47281 B/op	     113 allocs/op
BenchmarkCounter/ips_10k.txt/naive_(init_buf=65536,_max_buf=4096,_cap=0)-10        	    6702	    902085 ns/op	 1098937 B/op	   10083 allocs/op
BenchmarkCounter/ips_10k.txt/naive_(init_buf=65536,_max_buf=16,_cap=1024)-10       	    6997	    851916 ns/op	 1044781 B/op	   10068 allocs/op
BenchmarkCounter/ips_10k.txt/netip_map-10                                          	    5013	   1202126 ns/op	 1472217 B/op	   10083 allocs/op
BenchmarkCounter/ips_10k.txt/uint32_map-10                                         	    7746	    767142 ns/op	  308626 B/op	      83 allocs/op
BenchmarkCounter/ips_10k.txt/uint32_mmap-10                                        	   13142	    395007 ns/op	  152536 B/op	      37 allocs/op
BenchmarkCounter/ips_10k.txt/uint32_mmap_parallel(workers=2)-10                    	   13296	    447957 ns/op	  305147 B/op	      80 allocs/op
BenchmarkCounter/ips_10k.txt/uint32_mmap_parallel(workers=8)-10                    	   15140	    391646 ns/op	  306350 B/op	     104 allocs/op
BenchmarkCounter/ips_10k.txt/bitmap-10                                             	     849	   5918012 ns/op	536875149 B/op	       5 allocs/op
BenchmarkCounter/ips_10k.txt/bash_sort-10                                          	     523	  11482028 ns/op	   47282 B/op	     113 allocs/op
BenchmarkCounter/ips_100k.txt/naive_(init_buf=65536,_max_buf=4096,_cap=0)-10       	     568	  10745606 ns/op	 8655347 B/op	  100534 allocs/op
BenchmarkCounter/ips_100k.txt/naive_(init_buf=65536,_max_buf=16,_cap=1024)-10      	     568	  10455756 ns/op	 8601393 B/op	  100519 allocs/op
BenchmarkCounter/ips_100k.txt/netip_map-10                                         	     451	  13653807 ns/op	12096542 B/op	  100534 allocs/op
BenchmarkCounter/ips_100k.txt/uint32_map-10                                        	     790	   7717497 ns/op	 2439571 B/op	     534 allocs/op
BenchmarkCounter/ips_100k.txt/uint32_mmap-10                                       	    1454	   4108268 ns/op	 1218076 B/op	     261 allocs/op
BenchmarkCounter/ips_100k.txt/uint32_mmap_parallel(workers=2)-10                   	    1353	   4450944 ns/op	 2436106 B/op	     528 allocs/op
BenchmarkCounter/ips_100k.txt/uint32_mmap_parallel(workers=8)-10                   	    1736	   3513128 ns/op	 2437246 B/op	     552 allocs/op
BenchmarkCounter/ips_100k.txt/bitmap-10                                            	     380	  14220547 ns/op	536875154 B/op	       5 allocs/op
BenchmarkCounter/ips_100k.txt/bash_sort-10                                         	      62	  89654347 ns/op	   47293 B/op	     113 allocs/op
BenchmarkCounter/ips_1m.txt/naive_(init_buf=65536,_max_buf=4096,_cap=0)-10         	      25	 206723215 ns/op	127635449 B/op	 1008198 allocs/op
BenchmarkCounter/ips_1m.txt/naive_(init_buf=65536,_max_buf=16,_cap=1024)-10        	      26	 211576559 ns/op	127562332 B/op	 1008182 allocs/op
BenchmarkCounter/ips_1m.txt/netip_map-10                                           	      24	 243131212 ns/op	183519896 B/op	 1008197 allocs/op
BenchmarkCounter/ips_1m.txt/uint32_map-10                                          	      55	 105829439 ns/op	38872133 B/op	    8198 allocs/op
BenchmarkCounter/ips_1m.txt/uint32_mmap-10                                         	      75	  80517168 ns/op	19482992 B/op	    4101 allocs/op
BenchmarkCounter/ips_1m.txt/uint32_mmap_parallel(workers=2)-10                     	      58	  90709661 ns/op	38966582 B/op	    8208 allocs/op
BenchmarkCounter/ips_1m.txt/uint32_mmap_parallel(workers=8)-10                     	      85	  67049670 ns/op	38967283 B/op	    8232 allocs/op
BenchmarkCounter/ips_1m.txt/bitmap-10                                              	      58	  98946769 ns/op	536875145 B/op	       5 allocs/op
BenchmarkCounter/ips_1m.txt/bash_sort-10                                           	       4	1581060125 ns/op	   47806 B/op	     115 allocs/op
BenchmarkCounter/ips_10m.txt/naive_(init_buf=65536,_max_buf=4096,_cap=0)-10        	       2	2882983875 ns/op	1054759984 B/op	10065565 allocs/op
BenchmarkCounter/ips_10m.txt/naive_(init_buf=65536,_max_buf=16,_cap=1024)-10       	       2	2847174396 ns/op	1054705856 B/op	10065550 allocs/op
BenchmarkCounter/ips_10m.txt/netip_map-10                                          	       2	3262138166 ns/op	1503485584 B/op	10065565 allocs/op
BenchmarkCounter/ips_10m.txt/uint32_map-10                                         	       4	1529641052 ns/op	311698032 B/op	   65566 allocs/op
BenchmarkCounter/ips_10m.txt/uint32_mmap-10                                        	       5	1201089292 ns/op	155845009 B/op	   32773 allocs/op
BenchmarkCounter/ips_10m.txt/uint32_mmap_parallel(workers=2)-10                    	       4	1527997031 ns/op	311690024 B/op	   65552 allocs/op
BenchmarkCounter/ips_10m.txt/uint32_mmap_parallel(workers=8)-10                    	       5	1085800525 ns/op	311707790 B/op	   65578 allocs/op
BenchmarkCounter/ips_10m.txt/bitmap-10                                             	       6	 952043819 ns/op	536875144 B/op	       5 allocs/op
BenchmarkCounter/ips_10m.txt/bash_sort-10                                          	       1	22393380417 ns/op	   48936 B/op	     118 allocs/op
```
