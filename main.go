package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

var bindir string

func main() {
	if err := run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("missing program arguments")
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		spec, err := findSpec(dir)
		if err != nil {
			return err
		}
		if spec != nil {
			bindir = path.Join(dir, "bin")
			return spec.exec(os.Args[1:])
		}
		dir = path.Dir(dir)
		if dir == "" || dir == "/" {
			break
		}
	}
	return errors.New("go.tools not found")
}

func findSpec(fpath string) (spec, error) {
	fpath = path.Join(fpath, "go.tools")
	f, err := os.Open(fpath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var s spec
	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		if err := sc.Err(); err != nil {
			return s, err
		}
		s = append(s, sc.Text())
	}
	if s == nil {
		return s, fmt.Errorf("%s is empty", fpath)
	}
	return s, nil
}

type spec []string

func (spec spec) exec(args []string) error {
	for _, spec := range spec {
		def := strings.Split(spec, "@")
		if len(def) != 2 {
			return fmt.Errorf("invalid spec `%s`, expected as: program/path@version", spec)
		}
		pkg, pkgver := def[0], def[1]
		if !strings.HasSuffix(pkg, args[0]) {
			continue
		}
		var vfile *os.File
		var curver string
		var truncate bool
		vpath := path.Join(bindir, args[0]+".version")
		if _, err := os.Stat(vpath); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
			vfile, err = os.OpenFile(vpath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				return err
			}
			defer vfile.Close()
		} else {
			vfile, err = os.OpenFile(vpath, os.O_RDWR, 0)
			if err != nil {
				return err
			}
			defer vfile.Close()
			vdata, err := io.ReadAll(vfile)
			if err != nil {
				return err
			}
			curver = strings.Trim(string(vdata), " \n")
			truncate = true
		}
		exe := path.Join(bindir, args[0])
		if _, err := os.Stat(exe); err != nil || pkgver != curver {
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
			fmt.Println("Installing", spec)
			if err = command("go", "install", spec); err != nil {
				return err
			}
			if truncate {
				if err = vfile.Truncate(0); err != nil {
					return err
				}
				if _, err = vfile.Seek(0, 0); err != nil {
					return err
				}
			}
			if _, err = vfile.WriteString(pkgver); err != nil {
				return err
			}
		}
		return command(exe, args[1:]...)
	}
	return nil
}

func command(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if cmd.Err != nil {
		return cmd.Err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOBIN="+bindir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
