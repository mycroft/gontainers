package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// FS related functions
func makeFsPrivate() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	orig_fileinfo, err := os.Stat(path)
	orig_dev := orig_fileinfo.Sys().(*syscall.Stat_t).Dev

	for path != "/" {
		dir := filepath.Dir(path)

		fileinfo, _ := os.Stat(dir)
		new_dev := fileinfo.Sys().(*syscall.Stat_t).Dev

		if new_dev != orig_dev {
			break
		}

		path = dir
	}

	fmt.Println("Making", path, "rec private.")

	must(syscall.Mount("none", path, "", syscall.MS_PRIVATE|syscall.MS_REC, ""))
}

func prepareOverlay(orig string, dest string) string {
	// sudo mount -t overlay overlay -olowerdir=./orig,upperdir=./new,workdir=./work rootfs
	// mount("overlay", "/home/mycroft/dev/go-snippets/containers/fs/rootfs", "overlay", MS_MGC_VAL, "lowerdir=fs/orig,upperdir=fs/new"...) = 0

	// create temporary directories
	temp_dir, err := ioutil.TempDir("/tmp", "container")
	must(err)

	tmp_upper := fmt.Sprintf("%s/upper", temp_dir)
	must(os.Mkdir(tmp_upper, 0700))
	tmp_work := fmt.Sprintf("%s/work", temp_dir)
	must(os.Mkdir(tmp_work, 0700))

	options := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", orig, tmp_upper, tmp_work)

	err = syscall.Mount("overlay", dest, "overlay", syscall.MS_MGC_VAL, options)
	if err != nil {
		fmt.Printf("Could not mount overlay at %s\n", dest)
		panic(err)
	}

	return temp_dir
}

func unmountOverlay(dest string, temp_path string) {
	current_path, _ := os.Getwd()
	os.Chdir(dest)

	fmt.Println(dest)
	cmd := exec.Command("/usr/bin/tar", "cvf", "/tmp/rootfs.tar", ".")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	os.Chdir(current_path)

	os.RemoveAll(temp_path)

	must(syscall.Unmount("fs/rootfs", syscall.MNT_DETACH))
}

func Mkdev(major int64, minor int64) int {
	return int(((minor & 0xfff00) << 12) | ((major & 0xfff) << 8) | (minor & 0xff))
}
