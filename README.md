# ZerotierDNS

ztDNS is a dedicated DNS server for a ZeroTier virtual network.

## Overview

ztDNS pulls device names from Zerotier and makes them available by name using either IPv4 assigned addresses or IPv6 assigned addresses.

## Getting Started

### Traditional

If you prefer the traditional installation route:

#### Requirements

* [Go tools](https://golang.org/doc/install) - if not using a precompiled release

#### Install

1. First use `go get` to install the latest version, or download a precompiled release from [https://github.com/uxbh/ztdns/releases](https://github.com/uxbh/ztdns/releases)
    ``` bash
    go get -u github.com/uxbh/ztdns/
    go build
    ```
2. **If you are running on Linux**, run `sudo setcap cap_net_bind_service=+eip ./ztdns` to enable non-root users to bind privileged ports. On other operating systems, the program may need to be run as an administrator.

3. Add a new API access token to your user under the account tab at [https://my.zerotier.com](https://my.zerotier.com/).
    If you do not want to store your API access token in the configuration file you can also run the
    server with the `env` command: `env 'ZTDNS_ZT.API=<<APIToken>>' ./ztdns server`
4. Run `ztdns mkconfig` to generate a sample configuration file.
5. Add your API access token, Network names and IDs, and interface name to the configuration. Make sure you call ifconfig to determine your zerotier interface name. It won't always be zt0.
6. Start the server using `ztdns server`.
7. Add a DNS entry in your ZeroTier members pointing to the member running ztdns.

Once the server is up and running you will be able to resolve names based on the short name and suffix defined in the configuration file (zt by default) from ZeroTier.

```bash
dig @serveraddress member.domain.zt A
dig @serveraddress member.domain.zt AAAA
ping member.domain.zt
```

### Service

If you want to create a service so this starts on boot for Ubuntu, first add a bash script which spins up the server. I called mine `start-ztdns-server`:

```bash
#!/bin/sh
/path/to/ztdns server
```

Then add `ztdns.service` to `/etc/systemd/system/`. Make sure whatever you set `WorkingDirectory` to contains the .ztdns.toml configuration file.

```bash
[Unit]
Description=Zerotier DNS Server
[Service]
User=<user_name>
# The configuration file application.properties should be here:
#change this to your workspace
WorkingDirectory=/path/containing/ztdns_config/
#path to executable.
#executable is a bash script which calls jar file
ExecStart=/path/to/start-ztdns-server
SuccessExitStatus=143
TimeoutStopSec=10
Restart=on-failure
RestartSec=5
[Install]
WantedBy=multi-user.target
```

Then run systemctl enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable ztdns.service
sudo systemctl start ztdns.service
```

If you want to stop the service

```bash
sudo systemctl stop ztdns.service
sudo systemctl disable ztdns.service
```

### Docker

If you prefer to run the server with Docker:

#### Docker Requirements

* [Docker](https://docs.docker.com/install/)
* [Docker Compose](https://docs.docker.com/compose/install/)

#### Docker Install

1. Clone or download this repo
1. Create a `.ztdns.toml` file in the main directory by copying the `.ztdns.toml.example` file.
1. Add your API access token, Network ID, and interface name to the newly created configuration file.
1. By default it will be bound to port 5356 on the host, that can be changed to standard DNS port 53 by modifying the `docker-compose.yml` file. *You must be running Docker with root permissions in order to bind the privileged port properly.*
1. Run `docker-compose up` to start the server.
1. Add a DNS entry in your ZeroTier members pointing to the member running ztdns.

Once the server is up and running you will be able to resolve names based on the short name, domain and suffix defined in the configuration file (zt by default) from ZeroTier.

```bash
# remove -p 5356 if running on port 53
dig @127.0.0.1 -p 5356 member.domain.zt A
dig @127.0.0.1 -p 5356 member.domain.zt AAAA
ping member.domain.zt
```

## Contributing

Thanks for considering contributing to the project. We welcome contributions, issues or requests from anyone, and are grateful for any help. Problems or questions? Feel free to open an issue on GitHub.

Please make sure your contributions adhere to the following guidelines:

* Code must adhere to the official Go [formating](https://golang.org/doc/effective_go.html#formatting) guidelines  (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
* Pull requests need to be based on and opened against the `master` branch.
