package main

import (
	"github.com/vishvananda/netlink"
)

func endpointCreate() {
	hostIfName := "myveth0"
	containerIfName := "myveth1"

	veth, err := netlink.LinkByName("myveth0")
	if err == nil && veth != nil {
		return
	}

	veth = &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: hostIfName, TxQLen: 0},
		PeerName:  containerIfName}

	if err := netlink.LinkAdd(veth); err != nil {
		panic(err)
	}

	la := netlink.NewLinkAttrs()
	la.Name = "mycbridge"
	mybridge := &netlink.Bridge{la}

	err = netlink.LinkSetMaster(veth, mybridge)
	if err != nil {
		panic(err)
	}

	err = netlink.LinkSetUp(veth)
	if err != nil {
		panic(err)
	}
}

func endpointSetNs(nsfd int) {
	veth, err := netlink.LinkByName("myveth1")
	if err != nil {
		panic(err)
	}

	err = netlink.LinkSetNsFd(veth, nsfd)
	if err != nil {
		panic(err)
	}
}

func endpointDestroy() {
	hostInt, err := netlink.LinkByName("myveth0")
	if err != nil {
		panic(err)
	}

	err = netlink.LinkDel(hostInt)
	if err != nil {
		panic(err)
	}
}
