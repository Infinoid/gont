package options

import (
	"fmt"
	"net"

	"github.com/google/nftables/expr"
	g "github.com/stv0g/gont/pkg"
	"github.com/stv0g/gont/pkg/options/filters"
	nl "github.com/vishvananda/netlink"
)

func Route(network net.IPNet, gw net.IP) g.Route {
	return g.Route{
		Route: nl.Route{
			Dst: &network,
			Gw:  gw,
		},
	}
}

func DefaultGatewayIPv4(a, b, c, d byte) g.Route {
	return Route(g.DefaultIPv4Mask, net.IPv4(a, b, c, d))
}

func DefaultGatewayIP(fmts string, args ...any) g.Route {
	str := fmt.Sprintf(fmts, args...)

	gw := net.ParseIP(str)
	isV4 := gw.To4() != nil

	if isV4 {
		return Route(g.DefaultIPv4Mask, gw)
	}

	return Route(g.DefaultIPv6Mask, gw)
}

func Filter(h g.FilterHook, stmts ...filters.Statement) g.FilterRule {
	r := g.FilterRule{
		Hook:  h,
		Exprs: []expr.Any{},
	}

	for _, stmt := range stmts {
		r.Exprs = append(r.Exprs, stmt...)
	}

	return r
}
