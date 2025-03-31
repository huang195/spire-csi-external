package cgroups

import(
    "io"
    "os"
    "strconv"
    "strings"
    "fmt"
    //"io/ioutil"
    "unicode"
    "path/filepath"

    "k8s.io/apimachinery/pkg/types"
    "github.com/spiffe/spire/pkg/agent/common/cgroups"
)

const (
    cgroupPathTemplate = "/sys/fs/cgroup/kubelet.slice/kubelet-kubepods.slice/kubelet-kubepods-besteffort.slice/kubelet-kubepods-besteffort-pod%s.slice/cri-containerd-0000000000000000000000000000000000000000000000000000000000000000.scope"
)

func CreateFakeCgroup(uid string) error {
    podUID := canonicalizePodUID(uid)

    dirName := fmt.Sprintf(cgroupPathTemplate, podUID)

    err := os.Mkdir(dirName, 0755)
    if err != nil {
        return err
    }

    return nil
}

func DeleteFakeCgroup(uid string) error {
    podUID := canonicalizePodUID(uid)

    dirName := fmt.Sprintf(cgroupPathTemplate, podUID)

    err := os.Remove(dirName)
    if err != nil {
        return err
    }

    return nil
}

func EnterCgroup(pid int, path string) error {
    err := os.WriteFile(path, []byte(strconv.Itoa(pid)), 0644)
    if err != nil {
        return err
    }

    return nil
}

func GetPodProcsPath(uid string) string {
    podUID := canonicalizePodUID(uid)

    dirName := fmt.Sprintf(cgroupPathTemplate, podUID)

    return fmt.Sprintf("%s/cgroup.procs", dirName)
}

func GetMyCgroupProcsPath() (string, error) {
    pid := os.Getpid()

    cgroups, err := cgroups.GetCgroups(int32(pid), dirFS("/"))
    if err != nil {
        return "", err
    }

    for _, cgroup := range cgroups {
        //TODO: We are just going to use the first entry for now

        controllerList := cgroup.ControllerList
        groupPath := cgroup.GroupPath

        substrings := strings.SplitN(controllerList, "=", 2)
        if len(substrings) == 2 {
            controllerList= substrings[1]
        }

        return fmt.Sprintf("/sys/fs/cgroup/%s%s/cgroup.procs", controllerList, groupPath), nil
    }
    return "", fmt.Errorf("cannot find cgroup for the current process")
}

// canonicalizePodUID converts a Pod UID, as represented in a cgroup path, into
// a canonical form. Practically this means that we convert any punctuation to
// dashes, which is how the UID is represented within Kubernetes.
func canonicalizePodUID(uid string) types.UID {
    return types.UID(strings.Map(func(r rune) rune {
        if unicode.IsPunct(r) {
            r = '_'
        }
        return r
    }, uid))
}

type dirFS string

func (d dirFS) Open(p string) (io.ReadCloser, error) {
    return os.Open(filepath.Join(string(d), p))
}
