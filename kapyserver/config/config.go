package config

import (
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"time"

	"github.com/k3s-io/k3s/pkg/server"
	"github.com/kapycluster/corpy/kapyserver/util"
	"github.com/kapycluster/corpy/types"

	daemonsconfig "github.com/k3s-io/k3s/pkg/daemons/config"
	apinet "k8s.io/apimachinery/pkg/util/net"
	kubeapiserverflag "k8s.io/component-base/cli/flag"
	utilsnet "k8s.io/utils/net"
)

type ServerConfig struct {
	server.Config
	LBAddress   string
	ClusterCIDR string
	ServiceCIDR string
}

const (
	defaultClusterCIDR   = "10.11.0.0/16"
	defaultServiceCIDR   = "10.19.0.0/16"
	defaultClusterDomain = "cluster.local"
	defaultNodePortRange = "30000-32767"
)

var ContainerRuntimeReady = make(chan struct{})

// NewServerConfig sets defaults and creates a new ServerConfig
func NewServerConfig() (*ServerConfig, error) {
	config := &ServerConfig{}
	config.DisableAgent = true
	config.ControlConfig.Token = util.MustGetEnv(types.KapyServerToken)
	config.ControlConfig.KubeConfigOutput = util.MustGetEnv(types.KapyServerKubeConfigPath)
	config.ControlConfig.AdvertiseIP = util.MustGetEnv(types.KapyServerAdvertiseIP)
	config.LBAddress = util.MustGetEnv(types.KapyServerLoadBalancerAddress)
	config.ClusterCIDR = util.GetEnv(types.KapyServerClusterCIDR)
	config.ServiceCIDR = util.GetEnv(types.KapyServerServiceCIDR)
	config.ControlConfig.DataDir = util.GetEnv(types.KapyServerDataDir)
	config.ControlConfig.Datastore.Endpoint = util.MustGetEnv(types.KapyServerDatastore)
	config.ControlConfig.Datastore.NotifyInterval = 5 * time.Second
	config.ControlConfig.BindAddress = config.ControlConfig.AdvertiseIP
	config.ControlConfig.FlannelBackend = "wireguard-native"

	config.ControlConfig.HTTPSPort = 443
	config.ControlConfig.SupervisorPort = config.ControlConfig.HTTPSPort
	config.SupervisorPort = config.ControlConfig.HTTPSPort

	if config.ControlConfig.DataDir == "" {
		config.ControlConfig.DataDir = "/data"
	}

	if config.ClusterCIDR == "" {
		config.ClusterCIDR = defaultClusterCIDR
	}

	if config.ServiceCIDR == "" {
		config.ServiceCIDR = defaultServiceCIDR
	}

	_, nodeIPs, err := util.GetHostnameAndIPs("", []string{})
	if err != nil {
		return nil, err
	}
	for _, ip := range nodeIPs {
		config.ControlConfig.SANs = append(config.ControlConfig.SANs, ip.String())
	}

	config.ControlConfig.ServerNodeName = "kapy-server"
	config.ControlConfig.SANs = append(
		config.ControlConfig.SANs,
		"127.0.0.1",
		"::1",
		"localhost",
		config.LBAddress,
		config.ControlConfig.AdvertiseIP,
	)

	_, parsedClusterCIDR, err := net.ParseCIDR(config.ClusterCIDR)
	if err != nil {
		return nil, fmt.Errorf("invalid cluster CIDR: %w", err)
	}
	config.ControlConfig.ClusterIPRange = parsedClusterCIDR
	config.ControlConfig.ClusterIPRanges = []*net.IPNet{parsedClusterCIDR}

	_, parsedServiceCIDR, err := net.ParseCIDR(config.ServiceCIDR)
	if err != nil {
		return nil, fmt.Errorf("invalid service CIDR: %w", err)
	}
	config.ControlConfig.ServiceIPRange = parsedServiceCIDR
	config.ControlConfig.ServiceIPRanges = []*net.IPNet{parsedServiceCIDR}

	config.ControlConfig.ClusterDomain = "cluster.local"
	config.ControlConfig.ServiceNodePortRange, err = apinet.ParsePortRange(defaultNodePortRange)
	if err != nil {
		return nil, fmt.Errorf("invalid node port range: %w", err)
	}

	apiServerServiceIP, err := utilsnet.GetIndexedIP(config.ControlConfig.ServiceIPRange, 1)
	if err != nil {
		return nil, err
	}

	config.ControlConfig.SANs = append(config.ControlConfig.SANs, apiServerServiceIP.String())
	// Use the 10th IP of the service range as the cluster DNS IP
	clusterDNS, err := utilsnet.GetIndexedIP(config.ControlConfig.ServiceIPRange, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to configure cluster DNS address: %w", err)
	}

	config.ControlConfig.ClusterDNS = clusterDNS
	config.ControlConfig.ClusterDNSs = []net.IP{clusterDNS}

	// XXX: What does this do?
	config.ControlConfig.EgressSelectorMode = "cluster"

	config.ControlConfig.DefaultLocalStoragePath = filepath.Join(config.ControlConfig.DataDir, "storage")
	config.ControlConfig.DisableServiceLB = true
	config.ControlConfig.DisableCCM = true
	config.ControlConfig.Disables = map[string]bool{
		"metrics-server": true,
		"traefik":        true,
	}

	tlsCipherSuites := []string{
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
	}
	config.ControlConfig.ExtraAPIArgs = append(
		config.ControlConfig.ExtraAPIArgs,
		"tls-cipher-suites="+strings.Join(tlsCipherSuites, ","),
	)
	config.ControlConfig.CipherSuites = tlsCipherSuites
	config.ControlConfig.TLSCipherSuites, err = kubeapiserverflag.TLSCipherSuites(tlsCipherSuites)
	if err != nil {
		return nil, fmt.Errorf("invalid tls cipher suites: %w", err)
	}

	config.ControlConfig.Runtime = daemonsconfig.NewRuntime(ContainerRuntimeReady)

	return config, nil
}
