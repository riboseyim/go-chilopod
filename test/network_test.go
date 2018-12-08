package main

import (
	"fmt"
	"log"
	"net"
	"testing"

	chilopod_net "github.com/riboseyim/go-chilopod/network"
)

func TestDNS(t *testing.T) {

	chilopod_net.Query_dns("riboseyim.com")
	chilopod_net.Query_dns("qq.com")
	chilopod_net.Query_dns("facebook.com")

	host := "riboseyim-qiniu.riboseyim.com"
	cname, _ := net.LookupCNAME(host)
	log.Println("CNAME(域名),host", host, "cname", cname)

	server_ip := "8.8.8.8"
	ptr, _ := net.LookupAddr(server_ip)
	for _, ptrvalue := range ptr {
		fmt.Println("PTR(指针，IP地址的别名),server_ip", server_ip, "ptrvalue", ptrvalue)
	}

}
