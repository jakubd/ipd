# ipd

Go CLI app and library wrapper for Maxmind database lookups.

# Usage

Lookup one IP:

```shell
$ ipd 8.8.8.8
8.8.8.8,US,GOOGLE,AS15169,good_ip
```

Lookup list of IPs

Via pipe:

```shell
$ cat ips.txt | ipd pipe 
8.8.8.8,US,GOOGLE,AS15169,good_ip
8.8.4.4,US,GOOGLE,AS15169,good_ip
```

# Setup/Install

Currently, only Linux is supported. 

You need to download the maxmind databases yourself by setting up an account and downloading the libraries 
yourself from [Maxmind](https://dev.maxmind.com/geoip/geoip2/geolite2/)

It is recommended to manage the databases with [geoipupdate](https://github.com/maxmind/geoipupdate) it is currently in
[this contrib debian repos](https://packages.debian.org/buster/geoipupdate) so you can install with:

```shell
sudo apt install geoipupdate
```

You need to have both the country and ASN Geolite databases.

# License

[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)