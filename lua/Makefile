
all: lua-ffi lua-native-list lua-native-hash

lua-ffi:
	@echo "FFI"
	../IpDbGenerate-c.py > junk.c
	gcc -dynamiclib junk.c -o libjunk.dylib
	/usr/bin/time -p luajit -jon iplist0.lua
	/usr/bin/time -p luajit -joff iplist0.lua

lua-native-list:
	@echo ""
	@echo "LUA NATIVE LIST"
	./IpDbGenerate-lua1.py > iplist1.lua
	/usr/bin/time -p luajit -jon iplist1.lua
	/usr/bin/time -p luajit -joff iplist1.lua

lua-native-hash:
	echo ""
	@echo "LUA NATIVE HASH"
	./IpDbGenerate-lua2.py > iplist2.lua
	/usr/bin/time -p luajit -jon iplist2.lua
	/usr/bin/time -p luajit -joff iplist2.lua

clean:
	rm -f iplist1.lua iplist2.lua
	rm -f junk.c
	rm -f libjunk.*
	rm -rf *~
