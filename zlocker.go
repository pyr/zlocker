package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var (
	sessionTimeout   = flag.Int("t", 60, "Session timeout")
	waitPeriod       = flag.Int("w", 1, "Wait period before releasing lock")
	zookeeperCluster = flag.String("z", "", "Address of zookeeper cluster")
	lockName         = flag.String("l", "", "Name of zookeeper lock to request")
	flagVersion      = flag.Bool("v", false, "Display version and exit")

	version string
)

func unlock(lock *zk.Lock) {
	if lock.Unlock() != nil {
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Printf("%s %s\nGo version: %s (%s)\n",
			path.Base(os.Args[0]),
			version,
			runtime.Version(),
			runtime.Compiler,
		)
		os.Exit(0)
	}

	logger := log.New(os.Stderr, "zlocker: ", 0)
	if len(*zookeeperCluster) == 0 || len(*lockName) == 0 {
		logger.Fatal("need a host and lock name")
	}

	servers := strings.Split(*zookeeperCluster, ",")
	cmdline := strings.Join(flag.Args(), " ")
	if len(cmdline) == 0 {
		logger.Fatal("nothing to run, exiting")
	}

	cluster, _, err := zk.Connect(
		servers,
		time.Second*time.Duration(*sessionTimeout),
		zk.WithLogger(logger))
	if err != nil {
		logger.Fatal("cannot connect to cluster")
	}
	defer cluster.Close()

	acl := zk.WorldACL(zk.PermAll)
	lock := zk.NewLock(cluster, *lockName, acl)
	if err = lock.Lock(); err != nil {
		logger.Fatal("cannot lock, exiting")
		os.Exit(1)
	}
	defer unlock(lock)
	cmd := exec.Command("/bin/sh", "-c", cmdline)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		os.Exit(1)
	}
	if *waitPeriod > 0 {
		fmt.Println("command finished, sleeping")
		time.Sleep(time.Second * time.Duration(*waitPeriod))
	}
	os.Exit(0)
}
