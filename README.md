# ZerotierDNS

ztDNS is a dedicated DNS server for a ZeroTier virtual network.  

## Overview 

ztDNS pulls device names from Zerotier and makes them available by name using either IPv4 assigned addresses or IPv6 assigned addresses.  

## Installing

Using ztDNS is easy. First use ```go get``` to install the latest version. 
```
go get -u gitlab.com/uxbh/ztdns/
```
Alternatively run a precompiled release from [https://gitlab.com/uxbh/ztdns/tags](https://gitlab.com/uxbh/ztdns/tags).  

Next in your ZeroTier network members add a DNS entry pointing to the member running ztdns.  

## TODO

1. [ ] 1st Release
1. [ ] Nicer logging
1. [X] Mkconfig command  
1. [X] Update DNSDatabase with zt devices  
1. [ ] Improve Documentation  
1. [ ] Get listen IP by interface
1. [ ] Finish Readme  