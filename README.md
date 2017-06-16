# ZerotierDNS

ztDNS is a dedicated DNS server for a ZeroTier virtual network.  

## Overview 

ztDNS pulls device names from Zerotier and makes them available by name using either IPv4 assigned addresses or IPv6 assigned addresses.  

## Installing

1. First use ```go get``` to install the latest version, or download a precompiled relesase from [https://gitlab.com/uxbh/ztdns/tags](https://gitlab.com/uxbh/ztdns/tags)  
```
go get -u gitlab.com/uxbh/ztdns/
```
2. Add a new API access token to your user under the account tab at [https://my.zerotier.com](https://my.zerotier.com/).  
	If you do not want to store your API access token in the config file you can also run the  
	server with the ```env``` command: ```env 'ZTDNS_ZT.API=<<APIToken>>' ./ztdns server```
1. Run ```ztdns mkconfig``` to generate a sample config file.  
1. Add your API access token and Network ID, and interface name to the config.  
1. Start the server using ```ztdns server```.  
1. Add a DNS entry in your ZeroTier members pointing to the member running ztdns.  

Once the server is up and running you will be able to resolve names based on the short name and suffix defined in the config file (zt by default) from ZeroTier.  
```
dig @serveraddress member.zt A
dig @serveraddress member.zt AAAA
ping member.zt
```


## TODO

1. [ ] 1st Release
1. [X] Nicer logging
1. [X] Mkconfig command  
1. [X] Update DNSDatabase with zt devices  
1. [X] Improve Documentation  
1. [X] Get listen IP by interface
1. [X] Finish Readme  