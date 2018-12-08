package network

import (
	"fmt"
	"net"
)

func Query_dns(domain string) {
	iprecords, _ := net.LookupIP(domain)
	for _, ip := range iprecords {
		fmt.Println("A(主机地址),domain", domain, "ip", ip)
	}

	nameserver, _ := net.LookupNS(domain)
	for _, ns := range nameserver {
		fmt.Println("NS(域名服务器),domain", domain, "nameserver", ns)
	}

	mxrecords, _ := net.LookupMX(domain)
	for _, mx := range mxrecords {
		fmt.Println("MX(邮件交换),domain", domain, "mx.host", mx.Host, "mx.pref", mx.Pref)
	}

	txtrecords, _ := net.LookupTXT(domain)

	for _, txt := range txtrecords {
		fmt.Println("TXT(文本标识),domain", domain, "txt", txt)
	}
}
