package main

import (
	"fmt"
	"github.com/vishvananda/netns"
	"os"
	"os/exec"
	"syscall"
)

func parent() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Pid:", os.Getpid())

	endpointCreate()

	ns, err := netns.GetFromPid(os.Getpid())
	if err != nil {
		panic(err)
	}

	must(netns.Set(ns))

	temp_path := prepareOverlay("fs/orig", "fs/rootfs")

	if err := cmd.Start(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	path := fmt.Sprintf("/proc/%d/ns/net", cmd.Process.Pid)
	tmpfile, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	nsfd := int(tmpfile.Fd())
	endpointSetNs(nsfd)

	if err := cmd.Wait(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	endpointDestroy()

	unmountOverlay("fs/rootfs", temp_path)

}

func child() {
	// Create temp file for NS
	path := "/proc/self/ns/net"

	must(syscall.Mount("fs/rootfs", "fs/rootfs", "", syscall.MS_BIND, ""))
	must(os.MkdirAll("fs/rootfs/oldrootfs", 0700))
	must(syscall.PivotRoot("fs/rootfs", "fs/rootfs/oldrootfs"))
	must(os.Chdir("/"))

	must(syscall.Unmount("/oldrootfs", syscall.MNT_DETACH))
	must(os.Remove("/oldrootfs"))
	must(syscall.Mount("proc", "/proc", "proc", 0, ""))

	// Some devices
	syscall.Mknod("/dev/null", 0666, Mkdev(int64(1), int64(3)))
	syscall.Mknod("/dev/zero", 0666, Mkdev(int64(1), int64(5)))
	syscall.Mknod("/dev/random", 0666, Mkdev(int64(1), int64(8)))
	syscall.Mknod("/dev/urandom", 0666, Mkdev(int64(1), int64(9)))

	fmt.Println("Pid:", os.Getpid())

	ns, err := netns.GetFromPath(path)
	if err != nil {
		fmt.Println("cant find ns")
	}
	must(netns.Set(ns))

	routingUp()

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	routingDown()
}
