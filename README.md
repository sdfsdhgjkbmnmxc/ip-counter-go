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
BenchmarkCounter/ips_10.txt/naive_(init_buf=64_KiB,_max_buf=4.0_KiB,_cap=0)-10         	  216789	     15735 ns/op	   66288 B/op	      17 allocs/op
BenchmarkCounter/ips_10.txt/naive_(init_buf=64_KiB,_max_buf=16_B,_cap=1024)-10         	  179937	     21838 ns/op	  120441 B/op	      19 allocs/op
BenchmarkCounter/ips_10.txt/netip_map-10                                               	  300775	     11683 ns/op	    5008 B/op	      17 allocs/op
BenchmarkCounter/ips_10.txt/uint32_map-10                                              	  315296	     11191 ns/op	    4416 B/op	       7 allocs/op
BenchmarkCounter/ips_10.txt/uint32_mmap-10                                             	  268306	     13357 ns/op	     528 B/op	       7 allocs/op
BenchmarkCounter/ips_10.txt/uint32_mmap_parallel(w=2)-10                               	  207553	     17190 ns/op	    1112 B/op	      21 allocs/op
BenchmarkCounter/ips_10.txt/uint32_mmap_parallel(w=10)-10                              	  210796	     17127 ns/op	    1112 B/op	      21 allocs/op
BenchmarkCounter/ips_10.txt/uint32_mmap_parallel_limited(w=10,_cs=32_KiB)-10           	  201858	     17502 ns/op	    1256 B/op	      22 allocs/op
BenchmarkCounter/ips_10.txt/uint32_mmap_parallel_limited(w=10,_cs=32_MiB)-10           	  207559	     17509 ns/op	    1256 B/op	      22 allocs/op
BenchmarkCounter/ips_10.txt/bitmap-10                                                  	     661	   5226902 ns/op	536875162 B/op	       5 allocs/op
BenchmarkCounter/ips_10.txt/bitmap_mmap-10                                             	     686	   5237107 ns/op	536871270 B/op	       5 allocs/op
BenchmarkCounter/ips_10.txt/bitmap_mmap_roaring-10                                     	  225213	     15752 ns/op	    6080 B/op	      17 allocs/op
BenchmarkCounter/ips_10.txt/bash_sort-10                                               	     626	   5473207 ns/op	   47302 B/op	     113 allocs/op
BenchmarkCounter/ips_1k.txt/naive_(init_buf=64_KiB,_max_buf=4.0_KiB,_cap=0)-10         	   32026	    113012 ns/op	  190433 B/op	    1024 allocs/op
BenchmarkCounter/ips_1k.txt/naive_(init_buf=64_KiB,_max_buf=16_B,_cap=1024)-10         	   58870	     62197 ns/op	  136281 B/op	    1009 allocs/op
BenchmarkCounter/ips_1k.txt/netip_map-10                                               	   26784	    132787 ns/op	  180225 B/op	    1024 allocs/op
BenchmarkCounter/ips_1k.txt/uint32_map-10                                              	   41215	     87957 ns/op	   42288 B/op	      24 allocs/op
BenchmarkCounter/ips_1k.txt/uint32_mmap-10                                             	   75552	     47639 ns/op	   19368 B/op	       9 allocs/op
BenchmarkCounter/ips_1k.txt/uint32_mmap_parallel(w=2)-10                               	   37720	     95504 ns/op	   58002 B/op	      45 allocs/op
BenchmarkCounter/ips_1k.txt/uint32_mmap_parallel(w=10)-10                              	   34490	     99191 ns/op	   52314 B/op	     101 allocs/op
BenchmarkCounter/ips_1k.txt/uint32_mmap_parallel_limited(w=10,_cs=32_KiB)-10           	   33838	    105668 ns/op	   57968 B/op	      41 allocs/op
BenchmarkCounter/ips_1k.txt/uint32_mmap_parallel_limited(w=10,_cs=32_MiB)-10           	   34070	    109431 ns/op	   57968 B/op	      41 allocs/op
BenchmarkCounter/ips_1k.txt/bitmap-10                                                  	     699	   5119835 ns/op	536875157 B/op	       5 allocs/op
BenchmarkCounter/ips_1k.txt/bitmap_mmap-10                                             	     705	   5087069 ns/op	536871274 B/op	       5 allocs/op
BenchmarkCounter/ips_1k.txt/bitmap_mmap_roaring-10                                     	   16800	    216532 ns/op	  672342 B/op	    1024 allocs/op
BenchmarkCounter/ips_1k.txt/bash_sort-10                                               	     590	   5263624 ns/op	   47283 B/op	     113 allocs/op
BenchmarkCounter/ips_10k.txt/naive_(init_buf=64_KiB,_max_buf=4.0_KiB,_cap=0)-10        	    3962	    895928 ns/op	 1098940 B/op	   10083 allocs/op
BenchmarkCounter/ips_10k.txt/naive_(init_buf=64_KiB,_max_buf=16_B,_cap=1024)-10        	    4009	    847374 ns/op	 1044781 B/op	   10068 allocs/op
BenchmarkCounter/ips_10k.txt/netip_map-10                                              	    2962	   1193847 ns/op	 1472219 B/op	   10083 allocs/op
BenchmarkCounter/ips_10k.txt/uint32_map-10                                             	    4714	    789571 ns/op	  308626 B/op	      83 allocs/op
BenchmarkCounter/ips_10k.txt/uint32_mmap-10                                            	    9259	    403017 ns/op	  152537 B/op	      37 allocs/op
BenchmarkCounter/ips_10k.txt/uint32_mmap_parallel(w=2)-10                              	    5250	    662222 ns/op	  457520 B/op	     132 allocs/op
BenchmarkCounter/ips_10k.txt/uint32_mmap_parallel(w=10)-10                             	    6021	    594799 ns/op	  497045 B/op	     180 allocs/op
BenchmarkCounter/ips_10k.txt/uint32_mmap_parallel_limited(w=10,_cs=32_KiB)-10          	    3876	    923013 ns/op	  457537 B/op	     128 allocs/op
BenchmarkCounter/ips_10k.txt/uint32_mmap_parallel_limited(w=10,_cs=32_MiB)-10          	    4047	    843508 ns/op	  457474 B/op	     128 allocs/op
BenchmarkCounter/ips_10k.txt/bitmap-10                                                 	     626	   5750329 ns/op	536875164 B/op	       5 allocs/op
BenchmarkCounter/ips_10k.txt/bitmap_mmap-10                                            	     633	   5653885 ns/op	536871275 B/op	       5 allocs/op
BenchmarkCounter/ips_10k.txt/bitmap_mmap_roaring-10                                    	    1539	   2306595 ns/op	 6401222 B/op	   10030 allocs/op
BenchmarkCounter/ips_10k.txt/bash_sort-10                                              	     307	  11596557 ns/op	   47295 B/op	     113 allocs/op
BenchmarkCounter/ips_100k.txt/naive_(init_buf=64_KiB,_max_buf=4.0_KiB,_cap=0)-10       	     352	  10101583 ns/op	 8655644 B/op	  100534 allocs/op
BenchmarkCounter/ips_100k.txt/naive_(init_buf=64_KiB,_max_buf=16_B,_cap=1024)-10       	     354	  10033960 ns/op	 8601005 B/op	  100519 allocs/op
BenchmarkCounter/ips_100k.txt/netip_map-10                                             	     276	  12726699 ns/op	12096678 B/op	  100534 allocs/op
BenchmarkCounter/ips_100k.txt/uint32_map-10                                            	     482	   7402303 ns/op	 2439554 B/op	     534 allocs/op
BenchmarkCounter/ips_100k.txt/uint32_mmap-10                                           	     810	   4019251 ns/op	 1218069 B/op	     261 allocs/op
BenchmarkCounter/ips_100k.txt/uint32_mmap_parallel(w=2)-10                             	     558	   6504949 ns/op	 3653899 B/op	     807 allocs/op
BenchmarkCounter/ips_100k.txt/uint32_mmap_parallel(w=10)-10                            	     668	   5360976 ns/op	 3959985 B/op	     915 allocs/op
BenchmarkCounter/ips_100k.txt/uint32_mmap_parallel_limited(w=10,_cs=32_KiB)-10         	     379	   9389066 ns/op	 3654648 B/op	     803 allocs/op
BenchmarkCounter/ips_100k.txt/uint32_mmap_parallel_limited(w=10,_cs=32_MiB)-10         	     423	   8346931 ns/op	 3653845 B/op	     803 allocs/op
BenchmarkCounter/ips_100k.txt/bitmap-10                                                	     256	  14024014 ns/op	536875155 B/op	       5 allocs/op
BenchmarkCounter/ips_100k.txt/bitmap_mmap-10                                           	     271	  13248951 ns/op	536871272 B/op	       5 allocs/op
BenchmarkCounter/ips_100k.txt/bitmap_mmap_roaring-10                                   	     145	  24540188 ns/op	59318838 B/op	   95898 allocs/op
BenchmarkCounter/ips_100k.txt/bash_sort-10                                             	      39	  87189379 ns/op	   47321 B/op	     113 allocs/op
BenchmarkCounter/ips_1m.txt/naive_(init_buf=64_KiB,_max_buf=4.0_KiB,_cap=0)-10         	      16	 192389393 ns/op	127650335 B/op	 1008199 allocs/op
BenchmarkCounter/ips_1m.txt/naive_(init_buf=64_KiB,_max_buf=16_B,_cap=1024)-10         	      16	 194150432 ns/op	127626893 B/op	 1008186 allocs/op
BenchmarkCounter/ips_1m.txt/netip_map-10                                               	      13	 233841288 ns/op	183523874 B/op	 1008198 allocs/op
BenchmarkCounter/ips_1m.txt/uint32_map-10                                              	      33	 100274812 ns/op	38877834 B/op	    8199 allocs/op
BenchmarkCounter/ips_1m.txt/uint32_mmap-10                                             	      46	  72086086 ns/op	19483011 B/op	    4101 allocs/op
BenchmarkCounter/ips_1m.txt/uint32_mmap_parallel(w=2)-10                               	      36	  99161521 ns/op	58363202 B/op	   12313 allocs/op
BenchmarkCounter/ips_1m.txt/uint32_mmap_parallel(w=10)-10                              	      49	  69262452 ns/op	51045884 B/op	   10818 allocs/op
BenchmarkCounter/ips_1m.txt/uint32_mmap_parallel_limited(w=10,_cs=32_KiB)-10           	      21	 150319250 ns/op	58358721 B/op	   12306 allocs/op
BenchmarkCounter/ips_1m.txt/uint32_mmap_parallel_limited(w=10,_cs=32_MiB)-10           	      24	 137316137 ns/op	58348616 B/op	   12306 allocs/op
BenchmarkCounter/ips_1m.txt/bitmap-10                                                  	      34	  97383976 ns/op	536875158 B/op	       5 allocs/op
BenchmarkCounter/ips_1m.txt/bitmap_mmap-10                                             	      38	  89446697 ns/op	536871276 B/op	       5 allocs/op
BenchmarkCounter/ips_1m.txt/bitmap_mmap_roaring-10                                     	      13	 235142699 ns/op	413819184 B/op	  648364 allocs/op
BenchmarkCounter/ips_1m.txt/bash_sort-10                                               	       2	1516016520 ns/op	   48092 B/op	     115 allocs/op
BenchmarkCounter/ips_10m.txt/naive_(init_buf=64_KiB,_max_buf=4.0_KiB,_cap=0)-10        	       2	2792176188 ns/op	1054760160 B/op	10065566 allocs/op
BenchmarkCounter/ips_10m.txt/naive_(init_buf=64_KiB,_max_buf=16_B,_cap=1024)-10        	       2	2707571229 ns/op	1054705832 B/op	10065550 allocs/op
BenchmarkCounter/ips_10m.txt/netip_map-10                                              	       1	3154027000 ns/op	1503485904 B/op	10065569 allocs/op
BenchmarkCounter/ips_10m.txt/uint32_map-10                                             	       3	1467108139 ns/op	311698096 B/op	   65566 allocs/op
BenchmarkCounter/ips_10m.txt/uint32_mmap-10                                            	       3	1142932500 ns/op	155845048 B/op	   32774 allocs/op
BenchmarkCounter/ips_10m.txt/uint32_mmap_parallel(w=2)-10                              	       2	1629851562 ns/op	467539384 B/op	   98351 allocs/op
BenchmarkCounter/ips_10m.txt/uint32_mmap_parallel(w=10)-10                             	       3	1221324528 ns/op	506522536 B/op	  106584 allocs/op
BenchmarkCounter/ips_10m.txt/uint32_mmap_parallel_limited(w=10,_cs=32_KiB)-10          	       2	2306187146 ns/op	467612920 B/op	   98346 allocs/op
BenchmarkCounter/ips_10m.txt/uint32_mmap_parallel_limited(w=10,_cs=32_MiB)-10          	       2	2323456562 ns/op	467539272 B/op	   98346 allocs/op
BenchmarkCounter/ips_10m.txt/bitmap-10                                                 	       4	 947825812 ns/op	536875144 B/op	       5 allocs/op
BenchmarkCounter/ips_10m.txt/bitmap_mmap-10                                            	       4	 859516114 ns/op	536871280 B/op	       5 allocs/op
BenchmarkCounter/ips_10m.txt/bitmap_mmap_roaring-10                                    	       2	1925404354 ns/op	704773424 B/op	 1056727 allocs/op
BenchmarkCounter/ips_10m.txt/bash_sort-10                                              	       1	21980465958 ns/op	   48936 B/op	     118 allocs/op
```
