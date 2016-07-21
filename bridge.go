package main

import (
	"github.com/vishvananda/netlink"
)

// Bridge related functions
func bridgeCreate() {
	la := netlink.NewLinkAttrs()
	la.Name = "mycbridge"
	mybridge := &netlink.Bridge{la}
	err := netlink.LinkAdd(mybridge)
	if err != nil {
		panic(err)
	}

	addr, _ := netlink.ParseAddr("172.19.80.1/24")
	err = netlink.AddrAdd(mybridge, addr)
	if err != nil {
		panic(err)
	}

	err = netlink.LinkSetUp(mybridge)
	if err != nil {
		panic(err)
	}
}

func bridgeDestroy() {
	bridgeInt, err := netlink.LinkByName("mycbridge")
	if err != nil {
		panic(err)
	}

	err = netlink.LinkSetDown(bridgeInt)
	if err != nil {
		panic(err)
	}

	err = netlink.LinkDel(bridgeInt)
	if err != nil {
		panic(err)
	}
}
