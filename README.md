# ZerotierDNS

ztDNS is a dedicated DNS server for a ZeroTier virtual network.

## Overview

ztDNS pulls device names from Zerotier and makes them available by name using either IPv4 assigned addresses or IPv6 assigned addresses.

## Installing

1. First use ```go get``` to install the latest version, or download a precompiled relesase from [https://github.com/uxbh/ztdns/releases](https://github.com/uxbh/ztdns/releases)
```
go get -u github.com/uxbh/ztdns/
go build
```
2. **If you are running on linux**, run ```sudo setcap cap_net_bind_service=+eip ./ztdns``` to enable non-root users to bind privileged ports.
3. Add a new API access token to your user under the account tab at [https://my.zerotier.com](https://my.zerotier.com/).
	If you do not want to store your API access token in the config file you can also run the
	server with the ```env``` command: ```env 'ZTDNS_ZT.API=<<APIToken>>' ./ztdns server```
4. Run ```ztdns mkconfig``` to generate a sample config file.
5. Add your API access token and Network ID, and interface name to the config.
6. Start the server using ```ztdns server```.
7. Add a DNS entry in your ZeroTier members pointing to the member running ztdns.

Once the server is up and running you will be able to resolve names based on the short name and suffix defined in the config file (zt by default) from ZeroTier.
```
dig @serveraddress member.zt A
dig @serveraddress member.zt AAAA
ping member.zt
```

## Contributing

Thanks for considering contributing to the project. We welcome contributions, issues or requests from anyone, and are greatful for any help. Problems or questions? Feel free to open an issue on GitHub.

Please make sure yout contributions adhere to the following guidelines:
 * Code must adhere to the official Go [formating](https://golang.org/doc/effective_go.html#formatting) guidelines  (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Pull requests need to be based on and opened against the `master` branch.
