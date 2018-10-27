package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
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

func isMountPoint(path string) bool {
	cmd := exec.Command("mountpoint", path)
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

/// If goofys hasn't been mounted yet, mount!
/// If mounted, bind mount to appropriate place.
func Mount(target string, options map[string]string) interface{} {

	bucket := options["bucket"]
	subPath := options["subPath"]
	dirMode, ok := options["dirMode"]
	if !ok {
		dirMode = "0755"
	}
	fileMode, ok := options["fileMode"]
	if !ok {
		fileMode = "0644"
	}

	args := []string{
		"-o", "allow_other",
		"--dir-mode", dirMode,
		"--file-mode", fileMode,
	}

	if accessKey, ok := options["access-key"]; ok {
		args = append(args, "--access-key", accessKey)
	}
	if secretKey, ok := options["secret-key"]; ok {
		args = append(args, "--secret-key", secretKey)
	}

	mountPath := path.Join("/mnt/goofys", bucket)

	args = append(args, bucket, mountPath)

	if !isMountPoint(mountPath) {
		os.MkdirAll(mountPath, 0755)
		mountCmd := exec.Command("goofys", args...)
		mountCmd.Start()
	}

	srcPath := path.Join(mountPath, subPath)

	// Create subpath if it does not exist
	intDirMode, _ := strconv.ParseUint(dirMode, 8, 32)
	os.MkdirAll(srcPath, os.FileMode(intDirMode))

	// Now we rmdir the target, and then make a symlink to it!
	err := os.Remove(target)
	if err != nil {
		return makeResponse("Failure", err.Error())
	}

	err = os.Symlink(srcPath, target)

	return makeResponse("Success", "Mount completed!")
}

func Unmount(target string) interface{} {
	err := os.Remove(target)
	if err != nil {
		return makeResponse("Failure", err.Error())
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
