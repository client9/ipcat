
all: lua-ffi lua-native-list lua-native-hash

lua-ffi:
	../IpDbGenerate-c.py > junk.c
	gcc -dynamiclib junk.c -o libjunk.dylib
	/usr/bin/time -p luajit -jon iplist0.lua
	/usr/bin/time -p luajit -joff iplist0.lua

lua-native-list:
	./IpDbGenerate-lua.py > iplist1.lua
	/usr/bin/time -p luajit -jon iplist1.lua
	/usr/bin/time -p luajit -joff iplist1.lua

lua-native-hash:
	./IpDbGenerate-lua2.py > iplist2.lua
	/usr/bin/time -p luajit -jon iplist2.lua
	/usr/bin/time -p luajit -joff iplist2.lua
