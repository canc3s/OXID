package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	buffer1, _ = hex.DecodeString("05000b03100000004800000001000000b810b810000000000100000000000100c4fefc9960521b10bbcb00aa0021347a00000000045d888aeb1cc9119fe808002b10486002000000")
	buffer2, _ = hex.DecodeString("050000031000000018000000010000000000000000000500")
	begin, _   = hex.DecodeString("0700")
	end, _     = hex.DecodeString("00000900")
)

func getAddres(ip string, timeout time.Duration) {

	conn, err := net.DialTimeout("tcp", ip+":135", time.Second*timeout)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second * timeout))
	conn.Write(buffer1)
	reply := make([]byte, 1024)
	if n, err := conn.Read(reply); err != nil || n != 60 {
		return
	}

	conn.Write(buffer2)
	n, err := conn.Read(reply)
	if err != nil || n == 0 {
		return
	}
	start := bytes.Index(reply, begin)
	last := bytes.LastIndex(reply, end)

	datas := bytes.Split(reply[start:last], begin)
	fmt.Println("--------------------------------------\r\n[*] Retrieving network interface of", ip)
	for i := range datas {
		if i < 2 {
			continue
		}
		address := bytes.ReplaceAll(datas[i], []byte{0}, []byte{})
		fmt.Println("Address:", string(address))
	}
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func main() {

	host := flag.String("i", "", "single ip address")
	thread := flag.Int("t", 2000, "thread num")
	timeout := flag.Duration("time", 2, "timeout on connection, in seconds")
	netCIDR := flag.String("n", "", "CIDR notation of a network")
	flag.Parse()

	if *host == "" && *netCIDR == "" {
		flag.Usage()
	}

	if *host != "" {
		getAddres(*host, *timeout)
		return
	}

	c := make(chan struct{}, *thread)

	if *netCIDR != "" && *host == "" {
		ip, ipNet, err := net.ParseCIDR(*netCIDR)
		if err != nil {
			fmt.Println("invalid CIDR")
			return
		}
		var wg sync.WaitGroup

		for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
			wg.Add(1)
			go func(ip string) {
				c <- struct{}{}
				defer wg.Done()
				getAddres(ip, *timeout)
				<-c
			}(ip.String())
		}

		wg.Wait()
	}
}
