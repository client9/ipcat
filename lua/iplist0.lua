local ffi = require("ffi")
ffi.cdef[[
typedef struct ip_range_owner {
  uint32_t ip0;
  uint32_t ip1;
  const char* owner;
} ip_range_owner_t;

ip_range_owner_t iplist[];

int iplist_len;
]]

clib = ffi.load('junk');

function scan_list(x)
    local clib = clib;
    local iplist = clib.iplist;
    local imax = clib.iplist_len;
    for i=0,imax do
        local arec = iplist[i];
        if x >= arec.ip0 and x <= arec.ip1 then
           io.write('found it');
           break;
        end
    end
end


function loopit(imax)
    for i=0,imax do
        scan_list(999999999);
    end
end

loopit(100000);


