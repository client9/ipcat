Test of Lua JIT/FFI Performance
----------------------------------


This does a linear scan of the IPCAT datastructure, and does a trivial
check to make sure the loop isn't optimized away.

It's a array of structure of 3 items (int,int,string).

FFI - is a lua iterating over a C datastructure loaded from a shared library.

Lua 1 - iterates over a pure-lua table where each record is "{ int, int, string }" (list-style)

Lua 2 - iterators over a pure-lua table where each record is { ip0=int, ip1=int, owner=string}" (hash-style)

Performance
=============

In seconds.

```
|       | FFI   | Lua 1 | Lua 2 |
|-------|-------|-------|-------|
| jit   |  1.57 | 1.87  | 2.00  |
| nojit | 55.1  | 5.05  | 5.95  |
```

Memory
=============

In KB (1024 bytes kilobytes):

```
|       | FFI*  | Lua 1 | Lua 2 |
|-------|-------|-------|-------|
| jit   | 84    | 605   | 833   |
| nojit | 57    | 602   | 829   |

[*] shared library is 56KB managed externally
```

Raw Data
-----------------------------

```

FFI
../IpDbGenerate-c.py > junk.c
gcc -dynamiclib junk.c -o libjunk.dylib
/usr/bin/time -p luajit -jon iplist0.lua

84.69921875
real         1.60
user         1.59
sys          0.00
/usr/bin/time -p luajit -joff iplist0.lua

56.7685546875
real        53.39
user        53.34
sys          0.02

LUA NATIVE LIST
./IpDbGenerate-lua1.py > iplist1.lua
/usr/bin/time -p luajit -jon iplist1.lua

605.724609375
real         1.78
user         1.78
sys          0.00
/usr/bin/time -p luajit -joff iplist1.lua

602.4306640625
real         5.07
user         5.07
sys          0.00
echo ""

LUA NATIVE HASH
./IpDbGenerate-lua2.py > iplist2.lua
/usr/bin/time -p luajit -jon iplist2.lua

832.869140625
real         2.11
user         2.10
sys          0.00
/usr/bin/time -p luajit -joff iplist2.lua

829.0439453125
real         5.95
user         5.94
sys          0.00

```
