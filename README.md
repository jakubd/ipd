# ipd

Go CLI app and library wrapper for Maxmind database lookups.

# Usage

Lookup one IP:

```shell
$ ipd 8.8.8.8
8.8.8.8, US, GOOGLE, AS15169, good_ip
```

Lookup list of IPs

Via pipe:

```shell
 $ cat ips.txt | ipd pipe
8.8.8.8, US, GOOGLE, AS15169, good_ip
8.8.4.4, US, GOOGLE, AS15169, good_ip
1.1.1.1, AU, CLOUDFLARENET, AS13335, good_ip
```

Can optionally show links to common intel services with `-i` flag:

```shell
 $ ipd -i 8.8.8.8        
8.8.8.8, US, GOOGLE, AS15169, good_ip,  https://censys.io/ipv4/8.8.8.8,  https://www.shodan.io/host/8.8.8.8,  https://bgp.he.net/AS15169
```

Can take both domain/URL input if the `-r` flag is set.

```shell
 $ cat ips.txt | ipd pipe -r
https://freebsd.org, 96.47.72.84, US, NYINTERNET, AS11403, good_ip
one.one.one.one, 1.1.1.1, AU, CLOUDFLARENET, AS13335, good_ip
8.8.4.4, 8.8.4.4, US, GOOGLE, AS15169, good_ip
```

# Setup/Install

Currently, only Linux with GeoLite databases is supported. 

You need to download the maxmind databases yourself by setting up an account and downloading the libraries 
yourself from [Maxmind](https://dev.maxmind.com/geoip/geoip2/geolite2/)

Neccesary files are: `GeoLite2-ASN.mmdb` and `GeoLite2-ASN.mmdb`

It is recommended to manage the databases with [geoipupdate](https://github.com/maxmind/geoipupdate) it is currently in
[this contrib debian repos](https://packages.debian.org/buster/geoipupdate) so you can install with:

```shell
sudo apt install geoipupdate
```

You should put the databases in `/var/lib/GeoIP` directory as both `ipd` and `geoipupdate` use this directory. 

# License

[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)