package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"

    "github.com/go-logr/logr"
	"github.com/go-logr/zapr"
    "go.uber.org/zap"

    "github.com/spiffe/spiffe-csi/pkg/server"

    "github.com/huang195/spire-csi/pkg/driver"
)

var (
    nodeIDFlag               = flag.String("node-id", "", "Kubernetes Node ID. If unset, the node ID is obtained from the environment (i.e., -node-id-env)")
    nodeIDEnvFlag            = flag.String("node-id-env", "MY_NODE_NAME", "Envvar from which to obtain the node ID. Overridden by -node-id.")
    csiSocketPathFlag        = flag.String("csi-socket-path", "/csi-identity/csi.sock", "Path to the CSI socket")
    pluginNameFlag           = flag.String("plugin-name", "csi-identity.spiffe.io", "Plugin name to register")
    workloadAPISocketDirFlag = flag.String("workload-api-socket-dir", "", "Path to the Workload API socket directory")
)

func main() {
    flag.Usage = func() {
        fmt.Fprintln(os.Stderr, "spiffe-csi-driver provides an ephemeral inline CSI volume containing SPIRE workload identities")
        fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage:\n")
        flag.PrintDefaults()
    }
    flag.Parse()

    var log logr.Logger
    zapLog, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to set up logger: %v", err)
		os.Exit(1)
	}
	log = zapr.NewLogger(zapLog)

    nodeID := getNodeIDFromFlags()

    log.Info("Starting.",
        "nodeID", nodeID,
        "workloadAPISocketDir", *workloadAPISocketDirFlag,
        "csiSocketPathFlag", *csiSocketPathFlag)

    driver, err := driver.New(driver.Config{
		Log:                  log,
		NodeID:               nodeID,
		PluginName:           *pluginNameFlag,
		WorkloadAPISocketDir: *workloadAPISocketDirFlag,
	})
    if err != nil {
		log.Error(err, "Failed to create driver")
		os.Exit(1)
	}

    serverConfig := server.Config{
		Log:           log,
		CSISocketPath: *csiSocketPathFlag,
		Driver:        driver,
	}

    // Keep another process in the same cgroup so we don't lose the cgroup when we hop out of it
    sleepCmd := exec.Command("sleep", "infinity")
    err = sleepCmd.Start()

    if err := server.Run(serverConfig); err != nil {
		log.Error(err, "Failed to serve")
		os.Exit(1)
	}

    log.Info("Done")
}

func getNodeIDFromFlags() string {
	nodeID := os.Getenv(*nodeIDEnvFlag)
	if *nodeIDFlag != "" {
		nodeID = *nodeIDFlag
	}
	return nodeID
}
