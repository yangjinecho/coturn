package main

import (
	"fmt"
	"os/exec"
	"strings"
)

var (
	turnadminpath  = "bin/turnadmin"
	turnserverpath = "bin/turnadmin"
	ldlinuxpath    = "/lib64/ld-linux-x86-64.so.2"
)

type Entry struct {
	Name     string
	Realpath string
}

func ldd(lib string) (paths []Entry, err error) {
	cmd := exec.Command("ldd", lib)
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("ldd: %s => %s", strings.Join(cmd.Args, " "), err.Error())
		return
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		temps := strings.Fields(line)
		if len(temps) == 4 {
			paths = append(paths, Entry{temps[0], temps[2]})
		}
	}
	return
}

func copy(visited map[string]string) (err error) {
	for _, realpath := range visited {
		cmd := exec.Command("cp", "-f", realpath, "lib")
		err = cmd.Run()
		if err != nil {
			err = fmt.Errorf("copy: %s => %s", strings.Join(cmd.Args, " "), err.Error())
			return
		}
	}

	cmd := exec.Command("cp", "-f", ldlinuxpath, "lib/ld-linux.so")
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("copy: %s => %s", strings.Join(cmd.Args, " "), err.Error())
	}
	return
}

func createdb() (err error) {
	cmd := exec.Command("bash", "-c", "sqlite3 turndb/turndb < turndb/schema.sql")
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("createdb: %s => %s", strings.Join(cmd.Args, " "), err.Error())
	}
	return
}

func pack(tarname string) (err error) {
	cmd := exec.Command("tar", "cjf", tarname, "bin", "lib", "examples", "turndb")
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("pack: %s => %s", strings.Join(cmd.Args, " "), err.Error())
	}
	return
}

func main() {
	visited := map[string]string{}

	paths, err := ldd(turnadminpath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, path := range paths {
		visited[path.Name] = path.Realpath
	}

	paths, err = ldd(turnserverpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = copy(visited)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = createdb()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = pack("coturn.tar.bz2")
	if err != nil {
		fmt.Println(err)
		return
	}
}

// LD_LIBRARY_PATH=lib lib/ld-linux.so bin/turnserver -b turndb/turndb
