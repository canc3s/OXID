# OXID

通过windows的DCOM接口进行网卡进行信息枚举，无需认证，只要目标的135端口开放即可获得信息。可以有效提高内网渗透的效率，定位多网卡主机。

```
Usage of ./OXID:
  -i string
    	single ip address
  -n string
    	CIDR notation of a network
  -t int
    	thread num (default 2000)
  -time duration
    	timeout on connection, in seconds (default 2ns)

./OXID -i 192.168.1.1
./OXID -i 192.168.1.1/24
```

> 很久之前练习写的工具
