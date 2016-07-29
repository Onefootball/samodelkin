package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	errLog  = log.New(os.Stderr, "", log.LstdFlags)
	infoLog = log.New(os.Stdout, "", log.LstdFlags)

	packages = []string{
		"config",
		"container",
		"db",
		"handler",
		"log",
	}

	configFiles = []string{
		"config.json.dist",
		"config.json.prod",
		"config.json.staging",
	}
)

func main() {
	flag.Parse()
	projectPath := flag.Arg(0)
	if !filepath.IsAbs(projectPath) {
		errLog.Fatalf("'%s' is not an ansolute path, i.e. /$GOPATH/github.com/me/myproject", projectPath)
	}

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		errLog.Fatalln(err)
	}

	infoLog.Printf("> Bootstrapping project %s", projectPath)

	for _, pkg := range packages {
		if err := os.MkdirAll(filepath.Join(projectPath, pkg), 0755); err != nil {
			errLog.Fatal(err)
		}

		err := ioutil.WriteFile(
			filepath.Join(projectPath, pkg, pkg+".go"),
			[]byte(fmt.Sprintf("package %s", pkg)),
			0755,
		)

		if err != nil {
			errLog.Fatal(err)
		}
	}

	for _, cfgFile := range configFiles {
		if err := ioutil.WriteFile(filepath.Join(projectPath, cfgFile), []byte("{}"), 0755); err != nil {
			errLog.Fatal(err)
		}
	}

	wd, _ := os.Getwd()
	if err := createMakefile(wd, projectPath); err != nil {
		errLog.Fatal(err)
	}

	if err := createGitignore(wd, projectPath); err != nil {
		errLog.Fatal(err)
	}

	if err := createReadme(projectPath); err != nil {
		errLog.Fatal(err)
	}
}

func createMakefile(wd, projectPath string) error {
	mf, err := ioutil.ReadFile(filepath.Join(wd, "Makefile_sample"))
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(projectPath, "Makefile"), mf, 0755); err != nil {
		return err
	}

	return nil
}

func createGitignore(wd, projectPath string) error {
	gi, err := ioutil.ReadFile(filepath.Join(wd, "gitignore_sample"))
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(wd, ".gitignore"), gi, 0755); err != nil {
		return err
	}

	return nil
}

func createReadme(projectPath string) error {
	if err := ioutil.WriteFile(filepath.Join(projectPath, "README.md"), []byte("modify me"), 0755); err != nil {
		return err
	}

	return nil
}
