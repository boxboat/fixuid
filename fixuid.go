package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	config "github.com/go-ozzo/ozzo-config"
)

const ranFile = "/var/run/fixuid.ran"

var logger = log.New(os.Stderr, "", 0)
var quietFlag = flag.Bool("q", false, "quiet mode")

func main() {
	runtime.GOMAXPROCS(1)
	logger.SetPrefix("fixuid: ")
	flag.Parse()

	// development warning
	logInfo("fixuid should only ever be used on development systems. DO NOT USE IN PRODUCTION")

	argsWithoutProg := flag.Args()
	// detect what user we are running as
	runtimeUIDInt := os.Getuid()
	runtimeUID := strconv.Itoa(runtimeUIDInt)
	runtimeGIDInt := os.Getgid()
	runtimeGID := strconv.Itoa(runtimeGIDInt)

	// only run once on the system
	if _, err := os.Stat(ranFile); !os.IsNotExist(err) {
		logInfo("already ran on this system; will not attempt to change UID/GID")
		exitOrExec(runtimeUIDInt, runtimeGIDInt, argsWithoutProg)
	}

	// check that script is running as root
	if os.Geteuid() != 0 {
		logger.Fatalln(`fixuid is not running as root, ensure that the following criteria are met:
        - fixuid binary is owned by root: 'chown root:root /path/to/fixuid'
        - fixuid binary has the setuid bit: 'chmod u+s /path/to/fixuid'
        - NoNewPrivileges is disabled in container security profile
        - volume containing fixuid binary does not have the 'nosuid' mount option`)
	}

	// load config from /etc/fixuid/config.[json|toml|yaml|yml]
	rootConfig := config.New()
	configError := errors.New("could not find config at /etc/fixuid/config.[json|toml|yaml|yml]")
	var filePath string
	for _, fileName := range [...]string{"config.json", "config.toml", "config.yaml", "config.yml"} {
		filePath = path.Join("/etc/fixuid", fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			configError = rootConfig.Load(filePath)
			if configError != nil {
				logInfo("error when loading configuration file " + filePath)
			} else {
				break
			}
		}
	}
	if configError != nil {
		logger.Fatalln(configError)
	}

	// validate the container user from the config
	containerUser := rootConfig.GetString("user")
	if containerUser == "" {
		logger.Fatalln("cannot find key 'user' in configuration file " + filePath)
	}
	containerUID, containerUIDError := findUID(containerUser)
	if containerUIDError != nil {
		logger.Fatalln(containerUIDError)
	}
	if containerUID == "" {
		logger.Fatalln("user '" + containerUser + "' does not exist")
	}
	containerUIDInt, err := strconv.Atoi(containerUID)
	if err != nil {
		logger.Fatal(err)
	}
	containerUIDUint32 := uint32(containerUIDInt)

	// validate the container group from the config
	containerGroup := rootConfig.GetString("group")
	if containerGroup == "" {
		logger.Fatalln("cannot find key 'group' in configuration file " + filePath)
	}
	containerGID, containerGIDError := findGID(containerGroup)
	if containerGIDError != nil {
		logger.Fatalln(containerGIDError)
	}
	if containerGID == "" {
		logger.Fatalln("group '" + containerGroup + "' does not exist")
	}
	containerGIDInt, err := strconv.Atoi(containerGID)
	if err != nil {
		logger.Fatal(err)
	}
	containerGIDUint32 := uint32(containerGIDInt)

	// validate the paths from the config
	var paths []string
	err = rootConfig.Configure(&paths, "paths")
	if err != nil {
		switch err.(type) {
		case *config.ConfigPathError:
			paths = append(paths, "/")
		default:
			logger.Fatalln("key 'paths' is malformed; should be an array of strings in configuration file " + filePath)
		}
	}

	// declare uid/gid vars and
	var oldUID, newUID, oldGID, newGID string
	needChown := false

	// decide if need to change UIDs
	existingUser, existingUserError := findUser(runtimeUID)
	if existingUserError != nil {
		logger.Fatalln(existingUserError)
	}
	if existingUser == "" {
		logInfo("updating user '" + containerUser + "' to UID '" + runtimeUID + "'")
		needChown = true
		oldUID = containerUID
		newUID = runtimeUID
	} else {
		oldUID = ""
		newUID = ""
		if existingUser == containerUser {
			logInfo("runtime UID '" + runtimeUID + "' already matches container user '" + containerUser + "' UID")
		} else {
			logInfo("runtime UID '" + runtimeUID + "' matches existing user '" + existingUser + "'; not changing UID")
			needChown = true
		}
	}

	// deicide if need to change GIDs
	existingGroup, existingGroupError := findGroup(runtimeGID)
	if existingGroupError != nil {
		logger.Fatalln(existingGroupError)
	}
	if existingGroup == "" {
		logInfo("updating group '" + containerGroup + "' to GID '" + runtimeGID + "'")
		needChown = true
		oldGID = containerGID
		newGID = runtimeGID
	} else {
		oldGID = ""
		newGID = ""
		if existingGroup == containerGroup {
			logInfo("runtime GID '" + runtimeGID + "' already matches container group '" + containerGroup + "' GID")
		} else {
			logInfo("runtime GID '" + runtimeGID + "' matches existing group '" + existingGroup + "'; not changing GID")
			needChown = true
		}
	}

	// update /etc/passwd if necessary
	if oldUID != newUID || oldGID != newGID {
		err := updateEtcPasswd(containerUser, oldUID, newUID, oldGID, newGID)
		if err != nil {
			logger.Fatalln(err)
		}
	}

	// update /etc/group if necessary
	if oldGID != newGID {
		err := updateEtcGroup(containerGroup, oldGID, newGID)
		if err != nil {
			logger.Fatalln(err)
		}
	}

	// search entire filesystem and chown containerUID:containerGID to runtimeUID:runtimeGID
	if needChown {

		// proccess /proc/mounts
		mounts, err := parseProcMounts()
		if err != nil {
			logger.Fatalln(err)
		}

		// store the current mountpoint
		var mountpoint string

		// this function is called for every file visited
		visit := func(filePath string, fileInfo os.FileInfo, err error) error {

			// an error to lstat or filepath.readDirNames
			// see https://github.com/boxboat/fixuid/issues/4
			if err != nil {
				logInfo("error when visiting " + filePath)
				logInfo(err)
				return nil
			}

			// stat file to determine UID and GID
			sys, ok := fileInfo.Sys().(*syscall.Stat_t)
			if !ok {
				logInfo("cannot stat " + filePath)
				return filepath.SkipDir
			}

			// prevent recursing into mounts
			if findMountpoint(filePath, mounts) != mountpoint {
				if sys.Uid == containerUIDUint32 && sys.Gid == containerGIDUint32 {
					logInfo("skipping mounted path " + filePath)
				}
				if fileInfo.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			// only chown if file is containerUID:containerGID
			if sys.Uid == containerUIDUint32 && sys.Gid == containerGIDUint32 {
				logInfo("chown " + filePath)
				err := syscall.Lchown(filePath, runtimeUIDInt, runtimeGIDInt)
				if err != nil {
					logInfo("error changing owner of " + filePath)
					logInfo(err)
				}
				return nil
			}
			return nil
		}

		for _, path := range paths {
			// stat the path to ensure it exists
			_, err := os.Stat(path)
			if err != nil {
				logInfo("error accessing path: " + path)
				logInfo(err)
				continue
			}
			mountpoint = findMountpoint(path, mounts)

			logInfo("recursively searching path " + path)
			filepath.Walk(path, visit)
		}

	}

	// mark the script as ran
	if err := ioutil.WriteFile(ranFile, []byte{}, 0644); err != nil {
		logger.Fatalln(err)
	}

	// if the existing HOME directory is "/", change it to the user's home directory
	existingHomeDir := os.Getenv("HOME")
	if existingHomeDir == "/" {
		homeDir, homeDirErr := findHomeDir(runtimeUID)
		if homeDirErr == nil && homeDir != "" && homeDir != "/" {
			if len(argsWithoutProg) > 0 {
				os.Setenv("HOME", homeDir)
			} else {
				fmt.Println(`export HOME="` + strings.Replace(homeDir, `"`, `\"`, -1) + `"`)
			}
		}
	}

	// all done
	exitOrExec(runtimeUIDInt, runtimeGIDInt, argsWithoutProg)
}

func logInfo(v ...interface{}) {
	if !*quietFlag {
		logger.Println(v...)
	}
}

func exitOrExec(runtimeUIDInt int, runtimeGIDInt int, argsWithoutProg []string) {
	if len(argsWithoutProg) > 0 {
		// exec mode - de-escalate privileges and exec new process
		binary, err := exec.LookPath(argsWithoutProg[0])
		if err != nil {
			logger.Fatalln(err)
		}

		// de-escalate the user back to the original
		if err := syscall.Setreuid(runtimeUIDInt, runtimeUIDInt); err != nil {
			logger.Fatalln(err)
		}
		// de-escalate the group back to the original
		if err := syscall.Setregid(runtimeGIDInt, runtimeGIDInt); err != nil {
			logger.Fatalln(err)
		}

		// exec new process
		env := os.Environ()
		if err := syscall.Exec(binary, argsWithoutProg, env); err != nil {
			logger.Fatalln(err)
		}
	}

	// nothing to exec; exit the program
	os.Exit(0)
}

func searchColonDelimetedFile(filePath string, search string, searchOffset int, returnOffset int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), ":")
		if cols[searchOffset] == search {
			return cols[returnOffset], nil
		}
	}
	return "", nil
}

func findUID(user string) (string, error) {
	return searchColonDelimetedFile("/etc/passwd", user, 0, 2)
}

func findUser(uid string) (string, error) {
	return searchColonDelimetedFile("/etc/passwd", uid, 2, 0)
}

func findHomeDir(uid string) (string, error) {
	return searchColonDelimetedFile("/etc/passwd", uid, 2, 5)
}

func findGID(group string) (string, error) {
	return searchColonDelimetedFile("/etc/group", group, 0, 2)
}

func findGroup(gid string) (string, error) {
	return searchColonDelimetedFile("/etc/group", gid, 2, 0)
}

func updateEtcPasswd(user string, oldUID string, newUID string, oldGID string, newGID string) error {
	// user:pass:uid:gid:comment:home dir:shell
	file, err := os.Open("/etc/passwd")
	if err != nil {
		return err
	}

	newLines := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), ":")
		if oldUID != "" && newUID != "" && cols[0] == user && cols[2] == oldUID {
			cols[2] = newUID
		}
		if oldGID != "" && newGID != "" && cols[3] == oldGID {
			cols[3] = newGID
		}
		newLines += strings.Join(cols, ":") + "\n"
	}
	file.Close()

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := ioutil.WriteFile("/etc/passwd", []byte(newLines), 0644); err != nil {
		return err
	}

	return nil
}

func updateEtcGroup(group string, oldGID string, newGID string) error {
	// group:pass:gid:users
	file, err := os.Open("/etc/group")
	if err != nil {
		return err
	}

	newLines := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), ":")
		if oldGID != "" && newGID != "" && cols[0] == group && cols[2] == oldGID {
			cols[2] = newGID
		}
		newLines += strings.Join(cols, ":") + "\n"
	}
	file.Close()

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := ioutil.WriteFile("/etc/group", []byte(newLines), 0644); err != nil {
		return err
	}

	return nil
}

func parseProcMounts() (map[string]bool, error) {
	// device mountpoint type options dump fsck
	// spaces appear as \040
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}

	mounts := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cols := strings.Fields(scanner.Text())
		if len(cols) >= 2 {
			mounts[filepath.Clean(strings.Replace(cols[1], "\\040", " ", -1))] = true
		}
	}
	file.Close()

	return mounts, nil
}

func findMountpoint(path string, mounts map[string]bool) string {
	path = filepath.Clean(path)
	var lastPath string
	for path != lastPath {
		if _, ok := mounts[path]; ok {
			return path
		}
		lastPath = path
		path = filepath.Dir(path)
	}
	return "/"
}
