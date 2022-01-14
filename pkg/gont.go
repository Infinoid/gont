package gont

import (
	"errors"
	"os"
	"path"

	"github.com/vishvananda/netns"
	"kernel.org/pub/linux/libs/security/libcap/cap"
)

const (
	hostsFile = "/etc/hosts"
	netnsDir  = "/var/run/netns/"
	varDir    = "/var/run/gont"

	loopbackInterfaceName = "lo"
	bridgeInterfaceName   = "br"
)

func CheckCaps() error {
	c := cap.GetProc()
	if v, err := c.GetFlag(cap.Effective, cap.NET_ADMIN); err != nil || !v {
		return errors.New("missing NET_ADMIN capabilities")
	}
	return nil
}

// Identify returns the network and node name
// if the current process is running in a network netspace created by Gont
func Identify() (string, string, error) {
	curHandle, err := netns.Get()
	if err != nil {
		return "", "", err
	}

	for _, network := range NetworkNames() {
		for _, node := range NodeNames(network) {
			f := path.Join("/var/run/gont", network, "nodes", node, "ns", "net")

			handle, err := netns.GetFromPath(f)
			if err != nil {
				return "", "", err
			}

			if curHandle.Equal(handle) {
				return network, node, nil
			}
		}
	}

	return "", "", os.ErrNotExist
}

// TestConnectivity performs ICMP ping tests between all pairs of nodes in the network
func TestConnectivity(hosts ...*Host) error {
	for _, a := range hosts {
		for _, b := range hosts {
			if a != b {
				if _, err := a.Ping(b); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
