Test of Lua JIT/FFI Performance
----------------------------------


This does a linear scan of the IPCAT datastructure, and does a trivial
check to make sure the loop isn't optimized away.

It's a array of structure of 3 items (int,int,string).

FFI - is a lua iterating over a C datastructure loaded from a shared library.

Lua 1 - iterates over a pure-lua table where each record is "{ int, int, string }" (list-style)

Lua 2 - iterators over a pure-lua table where each record is { ip0=int, ip1=int, owner=string}" (hash-style)

```
|       | FFI   | Lua 1 | Lua 2 |
|-------|-------|-------|-------|
| jit   |  1.57 | 1.87  | 2.00  |
| nojit | 55.1  | 5.05  | 5.95  |
```

Raw Data
-----------------------------

```
FFI
../IpDbGenerate-c.py > junk.c
gcc -dynamiclib junk.c -o libjunk.dylib
/usr/bin/time -p luajit -jon iplist0.lua
real         1.57
user         1.56
sys          0.00
/usr/bin/time -p luajit -joff iplist0.lua
real        53.11
user        53.09
sys          0.02

LUA NATIVE LIST
./IpDbGenerate-lua1.py > iplist1.lua
/usr/bin/time -p luajit -jon iplist1.lua
real         1.87
user         1.77
sys          0.00
/usr/bin/time -p luajit -joff iplist1.lua
real         5.05
user         5.05
sys          0.00
echo ""

LUA NATIVE HASH
./IpDbGenerate-lua2.py > iplist2.lua
/usr/bin/time -p luajit -jon iplist2.lua
real         2.00
user         2.00
sys          0.00
/usr/bin/time -p luajit -joff iplist2.lua
real         5.95
user         5.95
sys          0.00
```
