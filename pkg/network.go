package gont

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"
	"syscall"

	"math/rand"
	"net"
	"os"
	"path/filepath"

	"github.com/vishvananda/netns"
)

type Network struct {
	Name     string
	Nodes    map[string]Node
	HostNode *Host
	BasePath string

	Persistent bool

	DefaultOptions Options
}

func HostNode(n *Network) *Host {
	baseNs, err := netns.Get()
	if err != nil {
		return nil
	}

	return &Host{
		BaseNode: BaseNode{
			name: "host",
			Namespace: &Namespace{
				Name:     "base",
				NsHandle: baseNs,
			},
			Network: n,
		},
	}
}

func GetNetworkNames() []string {
	names := []string{}

	nets, err := ioutil.ReadDir(varDir)
	if err != nil {
		return names
	}

	for _, net := range nets {
		if net.IsDir() {
			names = append(names, net.Name())
		}
	}

	sort.Strings(names)

	return names
}

func GetNodeNames(network string) []string {
	names := []string{}

	nodesDir := path.Join(varDir, network, "nodes")

	nets, err := ioutil.ReadDir(nodesDir)
	if err != nil {
		return names
	}

	for _, net := range nets {
		if net.IsDir() {
			names = append(names, net.Name())
		}
	}

	sort.Strings(names)

	return names
}

func GenerateNetworkName() string {
	existing := GetNetworkNames()

	for i := 0; i < 32; i++ {
		random := GetRandomName()

		index := sort.SearchStrings(existing, random)
		if index >= len(existing) || existing[index] != random {
			return random
		}
	}

	index := rand.Intn(len(Names))
	random := Names[index]

	return fmt.Sprintf("%s%d", random, rand.Intn(128)+1)
}

func CleanupAllNetworks() error {
	for _, name := range GetNetworkNames() {
		if err := CleanupNetwork(name); err != nil {
			return err
		}
	}

	return nil
}

func CleanupNetwork(name string) error {
	baseDir := filepath.Join(varDir, name)
	nodesDir := filepath.Join(baseDir, "nodes")

	fis, err := ioutil.ReadDir(nodesDir)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}

		nodeName := fi.Name()
		netNsName := fmt.Sprintf("gont-%s-%s", name, nodeName)

		netns.DeleteNamed(netNsName)
	}

	if err := os.RemoveAll(baseDir); err != nil {
		return err
	}

	return nil
}

func NewNetwork(name string, opts ...Option) (*Network, error) {
	if name == "" {
		name = GenerateNetworkName()
	}

	basePath := filepath.Join(varDir, name)

	n := &Network{
		Name:           name,
		BasePath:       basePath,
		Nodes:          map[string]Node{},
		DefaultOptions: Options{},
	}

	// Apply network specific options
	for _, opt := range opts {
		if nopt, ok := opt.(NetworkOption); ok {
			nopt.Apply(n)
		}
	}

	if stat, err := os.Stat(basePath); err == nil && stat.IsDir() {
		return nil, syscall.EEXIST
	}

	for _, path := range []string{"files", "nodes"} {
		path = filepath.Join(basePath, path)
		if err := os.MkdirAll(path, 0644); err != nil {
			return nil, err
		}
	}

	n.HostNode = HostNode(n)
	if n.HostNode == nil {
		return nil, errors.New("failed to create host node")
	}

	if err := n.UpdateHostsFile(); err != nil {
		return nil, fmt.Errorf("failed to update hosts file: %w", err)
	}

	return n, nil
}

func (n *Network) Teardown() error {
	for _, node := range n.Nodes {
		if err := node.Teardown(); err != nil {
			return err
		}
	}

	if n.BasePath != "" {
		os.RemoveAll(n.BasePath)
	}

	return nil
}

func (n *Network) Close() error {
	if !n.Persistent {
		if err := n.Teardown(); err != nil {
			return err
		}
	}

	return nil
}

func (n *Network) UpdateHostsFile() error {
	fn := filepath.Join(n.BasePath, "files", "etc", "hosts")
	if err := os.MkdirAll(filepath.Dir(fn), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	fmt.Fprintln(f, "# Autogenerated hosts file by Gont")

	hosts := map[string][]string{}

	IPv4loopback := net.IPv4(127, 0, 0, 1)

	hosts[IPv4loopback.String()] = []string{"localhost", "localhost.localdomain", "localhost4", "localhost4.localdomain4"}
	hosts[net.IPv6loopback.String()] = []string{"localhost", "localhost.localdomain", "localhost6", "localhost6.localdomain6"}

	add := func(name string, ip net.IP) {
		addr := ip.String()
		if hosts[addr] == nil {
			hosts[addr] = []string{}
		}

		hosts[addr] = append(hosts[addr], name)
	}

	for _, node := range n.Nodes {
		if host, ok := node.(*Host); ok {
			for _, i := range host.Interfaces {
				for _, a := range i.Addresses {
					add(host.name, a.IP)
					add(host.name+"-"+i.Name, a.IP)
				}
			}
		}
	}

	for addr, names := range hosts {
		fmt.Fprintf(f, "%s %s\n", addr, strings.Join(names, " "))
	}

	if err := f.Sync(); err != nil {
		return err
	}

	return nil
}
