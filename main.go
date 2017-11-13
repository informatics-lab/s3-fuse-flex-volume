package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
)

func makeResponse(status, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}
}

/// Return status
func Init() interface{} {
	resp := makeResponse("Success", "No Initialization required")
	resp["capabilities"] = map[string]interface{}{
		"attach": false,
	}
	return resp
}

/// If NFS hasn't been mounted yet, mount!
/// If mounted, bind mount to appropriate place.
func Mount(target string, options map[string]string) interface{} {
	opts := strings.Split(options["mountOptions"], ",")
	sort.Strings(opts)
	sortedOpts := strings.Join(opts, ",")

	subPath := options["subPath"]
	createIfNecessary := options["createIfNecessary"] == "true"
	createModeUint64, err := strconv.ParseUint(options["createMode"], 0, 32)
	createMode := os.FileMode(createModeUint64)

	sharePath := options["share"]
	fmt.Printf("target:%s, share:%s, options:%v", target, sharePath, options)
	//createUid := strconv.Atoi(options["createUid"])
	//createGid := strconv.Atoi(options["createGid"])

	mountPath := fmt.Sprintf("/mnt/nfsflexvolume/%s/options/%s", sharePath, sortedOpts)

	os.MkdirAll(mountPath, 0755)

	mountCmd := exec.Command("mount", "-t", "nfs4", sharePath, mountPath, "-o", sortedOpts)
	out, err := mountCmd.CombinedOutput()
	if err != nil {
		return makeResponse("Failure", fmt.Sprintf("%s: %s", err.Error(), out))
	}

	srcPath := path.Join(mountPath, subPath)

	if createIfNecessary {
		fmt.Printf("test %v", createIfNecessary)
		err := os.MkdirAll(srcPath, createMode)
		if err != nil {
			return makeResponse("Failure", fmt.Sprintf("Could not create subPath: %s", err.Error()))
		}
	}

	bindMountCmd := exec.Command("mount", "--bind", path.Join(mountPath, subPath), target)
	out, err = bindMountCmd.CombinedOutput()
	if err != nil {
		return makeResponse("Failure", fmt.Sprintf("%s: %s", err.Error(), out))
	}

	return makeResponse("Success", "Mount completed!")
}

func Unmount(mountPath string) interface{} {
	umountCmd := exec.Command("umount", mountPath)
	out, err := umountCmd.CombinedOutput()
	if err != nil {
		return makeResponse("Failure", fmt.Sprintf("%s: %s", err.Error(), out))
	}
	return makeResponse("Success", "Successfully unmounted")
}

func printJSON(data interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", string(jsonBytes))
}

func main() {
	switch action := os.Args[1]; action {
	case "init":
		printJSON(Init())
	case "mount":
		optsString := os.Args[3]
		opts := make(map[string]string)
		json.Unmarshal([]byte(optsString), &opts)
		printJSON(Mount(os.Args[2], opts))
	case "unmount":
		printJSON(Unmount(os.Args[2]))
	default:
		printJSON(makeResponse("Not supported", fmt.Sprintf("Operation %s is not supported", action)))
	}

}
