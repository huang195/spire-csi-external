package procs

import (
    "os/exec"
    "fmt"
    "strings"
    "strconv"
    "io/ioutil"
    "encoding/json"
)

type Proc struct {
    Pid     int
    Cmdline string
}

func GetPeerProcs(ppid int) ([]Proc, error) {

    var procs []Proc

    procFiles, err := ioutil.ReadDir("/proc")
   	if err != nil {
		return nil, err
	}

    for _, proc := range procFiles {
        if pidNum, err := strconv.Atoi(proc.Name()); err == nil {

            statusFile := fmt.Sprintf("/proc/%d/status", pidNum)
            statusContent, err := ioutil.ReadFile(statusFile)
            if err != nil {
                continue
            }

            for _, line := range strings.Split(string(statusContent), "\n") {
                if strings.HasPrefix(line, "PPid:") {
                    fields := strings.Fields(line)
                    if len(fields) == 2 {
                        childPPID, err := strconv.Atoi(fields[1])
                        if err == nil && childPPID == ppid {
                            // If the PPID matches, this is a child process
                            cmdline, _ := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pidNum))
                            p := Proc{Pid: pidNum, Cmdline: string(cmdline)}
                            procs = append(procs, p)
                        }
                    }
                }
            }
        }
    }

    return procs, nil
}

func GetPodSandboxID(podName, podNamespace string) (string, error) {
    cmd := exec.Command("crictl", "pods", "--name", podName, "--namespace", podNamespace, "-o", "json")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }

    sandboxID, err := extractSandboxIDFromJSON(output)
    return sandboxID, err
}

func extractSandboxIDFromJSON(input []byte) (string, error) {
    var result map[string]interface{}

    if err := json.Unmarshal(input, &result); err != nil {
        return "", err
    }

    items, ok := result["items"].([]map[string]interface{})
    if !ok {
        return "", fmt.Errorf("Unexpected json result returned from crictl")
    }

    if len(items) != 1 {
        return "", fmt.Errorf("Unexpected number of pods returned from crictl")
    }

    if id, ok := items[0]["id"].(string); ok {
        return id, nil
    }

    return "", fmt.Errorf("Cannot find id field of the returned pod")
}


