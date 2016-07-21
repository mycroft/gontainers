package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"net"
)

func createRoutes(veth netlink.Link) []*netlink.Route {
	gw, _, err := net.ParseCIDR("172.19.80.1/32")
	if err != nil {
		panic(err)
	}

	routes := make([]*netlink.Route, 1)

	routes[0] = &netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: veth.Attrs().Index,
		Gw:        gw,
	}

	return routes
}

func routingUp() {
	veth, err := netlink.LinkByName("myveth1")
	if err != nil {
		panic(err)
	}

	err = netlink.LinkSetUp(veth)
	if err != nil {
		panic(err)
	}

	addr, _ := netlink.ParseAddr("172.19.80.2/24")
	err = netlink.AddrAdd(veth, addr)
	if err != nil {
		panic(err)
	}

	routes := createRoutes(veth)

	for _, route := range routes {
		fmt.Println("Adding route", route)
		err := netlink.RouteAdd(route)
		if err != nil {
			fmt.Println(err)
			// panic(err)
		}
	}
}

func routingDown() {
	veth, err := netlink.LinkByName("myveth1")
	if err != nil {
		panic(err)
	}

	routes := createRoutes(veth)

	for _, route := range routes {
		fmt.Println("Removing route", route)
		err := netlink.RouteDel(route)
		if err != nil {
			// panic(err)
			fmt.Println(err)
		}
	}

	err = netlink.LinkSetDown(veth)
	if err != nil {
		panic(err)
	}
}
