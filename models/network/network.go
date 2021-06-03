package network

import (
	"fmt"
	"harbored/config"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type Network struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

// Search for networks suitable for use
func NetScan() *[]Network {
	networks := make([]Network, 0)
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if isPrivateIP(ip) {
				// Detect Harbored web-server
				client := http.Client{
					Timeout: time.Millisecond * 400,
				}
				resp, err := client.Get("http://" + strings.Split(addr.String(), "/")[0] + config.Config.ServerPort + "/ping")
				if err != nil {
					fmt.Println(err)
				} else {
					defer resp.Body.Close()
					body, _ := ioutil.ReadAll(resp.Body)
					if string(body) == "pong" {
						network := Network{
							Name: i.Name,
							Ip:   addr.String(),
						}
						networks = append(networks, network)
					}
				}
			}
		}
	}
	return &networks
}

var privateIPBlocks []*net.IPNet

// Initialize privateIPBlocks
func init() {
	for _, cidr := range []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

func isPrivateIP(ip net.IP) bool {
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
