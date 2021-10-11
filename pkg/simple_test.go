package gont_test

import (
	"testing"

	g "github.com/stv0g/gont/pkg"
	o "github.com/stv0g/gont/pkg/options"
)

// TestPing performs and end-to-end ping test
// between two hosts on a switched topology
//
//  h1 <-> sw <-> h2
func TestPing(t *testing.T) {
	var (
		err    error
		n      *g.Network
		sw     *g.Switch
		h1, h2 *g.Host
	)

	if n, err = g.NewNetwork(nname, opts...); err != nil {
		t.Errorf("Failed to create network: %s", err)
		t.FailNow()
	}
	defer n.Close()

	if sw, err = n.AddSwitch("sw"); err != nil {
		t.Errorf("Failed to create switch: %s", err)
		t.FailNow()
	}

	if h1, err = n.AddHost("h1",
		o.Interface("veth0", sw,
			o.AddressIPv4(10, 0, 0, 1, 24)),
	); err != nil {
		t.Errorf("Failed to create host: %s", err)
		t.FailNow()
	}

	if h2, err = n.AddHost("h2",
		o.Interface("veth0", sw,
			o.AddressIPv4(10, 0, 0, 2, 24)),
	); err != nil {
		t.Errorf("Failed to create host: %s", err)
		t.FailNow()
	}

	if err := g.TestConnectivity(h1, h2); err != nil {
		t.Errorf("Failed to test connectivity")
	}
}

// TestPingDirect performs and end-to-end ping test between two
// directly connected hosts
//
// h1 <-> h2
func TestPingDirect(t *testing.T) {
	var (
		err    error
		n      *g.Network
		h1, h2 *g.Host
	)

	if n, err = g.NewNetwork(nname, opts...); err != nil {
		t.Errorf("Failed to create network: %s", err)
		t.FailNow()
	}
	defer n.Close()

	if h1, err = n.AddHost("h1"); err != nil {
		t.Errorf("Failed to create host: %s", err)
		t.FailNow()
	}

	if h2, err = n.AddHost("h2"); err != nil {
		t.Errorf("Failed to create host: %s", err)
		t.FailNow()
	}

	if err := n.AddLink(
		o.Interface("veth0", h1,
			o.AddressIPv4(10, 0, 0, 1, 24)),
		o.Interface("veth0", h2,
			o.AddressIPv4(10, 0, 0, 2, 24)),
	); err != nil {
		t.Errorf("Failed to connect hosts: %s", err)
		t.FailNow()
	}

	if _, _, err = h1.Run("cat", "/etc/hosts"); err != nil {
		t.Errorf("Failed to show /etc/hosts file")
	}

	if err := g.TestConnectivity(h1, h2); err != nil {
		t.Errorf("Failed to test connectivity")
	}
}

// TestPingMultiHop performs and end-to-end ping test
// through a routed multi-hop topology
//
//  h1 <-> sw1 <-> r1 <-> sw2 <-> h2
func TestPingMultiHop(t *testing.T) {
	var (
		err      error
		n        *g.Network
		sw1, sw2 *g.Switch
		h1, h2   *g.Host
	)

	if n, err = g.NewNetwork(nname, opts...); err != nil {
		t.Errorf("Failed to create network: %s", err)
		t.FailNow()
	}
	defer n.Close()

	if sw1, err = n.AddSwitch("sw1"); err != nil {
		t.Errorf("Failed to add switch: %s", err)
		t.FailNow()
	}

	if sw2, err = n.AddSwitch("sw2"); err != nil {
		t.Errorf("Failed to add switch: %s", err)
		t.FailNow()
	}

	if h1, err = n.AddHost("h1",
		o.GatewayIPv4(10, 0, 1, 1),
		o.Interface("veth0", sw1,
			o.AddressIPv4(10, 0, 1, 2, 24),
		),
	); err != nil {
		t.Errorf("Failed to add host: %s", err)
		t.FailNow()
	}

	if h2, err = n.AddHost("h2",
		o.GatewayIPv4(10, 0, 2, 1),
		o.Interface("veth0", sw2,
			o.AddressIPv4(10, 0, 2, 2, 24)),
	); err != nil {
		t.Errorf("Failed to add host: %s", err)
		t.FailNow()
	}

	if _, err := n.AddRouter("r1",
		o.Interface("veth0", sw1,
			o.AddressIPv4(10, 0, 1, 1, 24),
		),
		o.Interface("veth1", sw2,
			o.AddressIPv4(10, 0, 2, 1, 24),
		),
	); err != nil {
		t.Errorf("Failed to add router: %s", err)
		t.FailNow()
	}

	if err := g.TestConnectivity(h1, h2); err != nil {
		t.FailNow()
	}

	if err := h1.Traceroute(h2); err != nil {
		t.Fail()
	}
}
