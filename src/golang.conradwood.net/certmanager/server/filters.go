package main

import (
	"fmt"
	"strings"
)

func is_public_spam(name string) bool {
	name = strings.ToLower(name)
	if strings.Contains(name, ".proxy.") {
		return true
	}
	if strings.Count(name, "www.") > 1 {
		return true
	}
	if strings.Contains(name, "plesk") {
		return true
	}
	return false

}
func rewrite_host_name(host string) string {
	hostname := strings.ToLower(host)
	if is_public_spam(host) {
		fmt.Printf("rewritten hostname %s to be exactly '%s'\n", host, hostname)
		return "proxy.conradwood.net"
	}
	if strings.Contains(strings.ToLower(hostname), ".proxy.conradwood") {
		hostname = "proxy.conradwood.net"
		fmt.Printf("rewritten hostname %s to be exactly '%s'\n", host, hostname)
	}
	return hostname

}
