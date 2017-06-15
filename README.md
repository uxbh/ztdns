# ZerotierDNS

ztDNS is a dedicated DNS server for a ZeroTier virtual network.  

## Overview 

ztDNS pulls device names from Zerotier and makes them available by name using either IPv4 assigned addresses or IPv6 assigned addresses.  

## Installing

Using ztDNS is easy. First use ```go get``` to install the latest version. 
```
go get -u gitlab.com/uxbh/ztdns/
```
Alternatively run a precompiled version from TBD.  

Next in your ZeroTier network members add a DNS entry pointing to the member running ztdns.  


## TODO

1. [ ] Mkconfig command  
1. [ ] Update DNSDatabase with zt devices  
1. [ ] Improve Documentation  
1. [ ] Finish Readme  
  

...