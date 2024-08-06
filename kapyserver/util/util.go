package util

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/urfave/cli"
	apinet "k8s.io/apimachinery/pkg/util/net"
)

func MustGetEnv(key string) string {
	v := GetEnv(key)
	if v == "" {
		log.Fatalf("env var not set: %s", key)
	}

	return v
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

// GetHostnameAndIPs takes a node name and list of IPs, usually from CLI args.
// If set, these are used to return the node's name and addresses. If not set,
// the system hostname and primary interface addresses are returned instead.
func GetHostnameAndIPs(name string, nodeIPs []string) (string, []net.IP, error) {
	ips := []net.IP{}
	if len(nodeIPs) == 0 {
		hostIP, err := apinet.ChooseHostInterface()
		if err != nil {
			return "", nil, err
		}
		ips = append(ips, hostIP)
		// If IPv6 it's an IPv6 only node
		if hostIP.To4() != nil {
			hostIPv6, err := apinet.ResolveBindAddress(net.IPv6loopback)
			if err == nil && !hostIPv6.Equal(hostIP) {
				ips = append(ips, hostIPv6)
			}
		}
	} else {
		var err error
		ips, err = ParseStringSliceToIPs(nodeIPs)
		if err != nil {
			return "", nil, fmt.Errorf("invalid node-ip: %w", err)
		}
	}

	if name == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return "", nil, err
		}
		name = hostname
	}

	// Use lower case hostname to comply with kubernetes constraint:
	// https://github.com/kubernetes/kubernetes/issues/71140
	name = strings.ToLower(name)

	return name, ips, nil
}

// ParseStringSliceToIPs converts slice of strings that in turn can be lists of comma separated unparsed IP addresses
// into a single slice of net.IP, it returns error if at any point parsing failed
func ParseStringSliceToIPs(s cli.StringSlice) ([]net.IP, error) {
	var ips []net.IP
	for _, unparsedIP := range s {
		for _, v := range strings.Split(unparsedIP, ",") {
			ip := net.ParseIP(v)
			if ip == nil {
				return nil, fmt.Errorf("invalid ip format '%s'", v)
			}
			ips = append(ips, ip)
		}
	}

	return ips, nil
}
