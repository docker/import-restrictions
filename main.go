/*
   Copyright 2020 Docker, Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type goPackage struct {
	ImportPath string   `yaml:",omitempty"`
	Deps       []string `yaml:",omitempty"`
}

type importRestrictions struct {
	Path             string   `yaml:"path,omitempty"`
	ForbiddenImports []string `yaml:"forbiddenImports,omitempty"`
}

func main() {
	app := cli.App{
		Name:            "import-restrictions",
		Usage:           "Restrict imports in your go project",
		ArgsUsage:       "config-file",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "configuration",
				Aliases: []string{"c"},
				Value:   "import-restrictions.yaml",
			},
		},
		Action: func(clix *cli.Context) error {
			return run(clix.String("configuration"))
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(configFile string) error {
	cfg, err := loadConfig(configFile)
	if err != nil {
		return err
	}

	var importErrors *multierror.Error

	for _, dir := range cfg {
		dirImports, err := getDirDeps(dir.Path)
		if err != nil {
			return err
		}

		for _, dirImport := range dirImports {
			for _, dependency := range dirImport.Deps {
				if stringSliceContains(dir.ForbiddenImports, dependency) {
					importErrors = multierror.Append(importErrors, fmt.Errorf("Forbidden import %q in package %s", dependency, dirImport.ImportPath))
					importErrors = multierror.Append(importErrors, fmt.Errorf("Forbidden import %q in package %s", dependency, dirImport.ImportPath))
				}
			}
		}
	}

	if importErrors != nil {
		importErrors.ErrorFormat = formatErrors
	}

	return importErrors.ErrorOrNil()
}

func formatErrors(errs []error) string {
	messages := make([]string, len(errs))
	for i, err := range errs {
		messages[i] = "* " + err.Error()
	}
	return strings.Join(messages, "\n")
}

func loadConfig(cfg string) ([]importRestrictions, error) {
	config, err := ioutil.ReadFile(cfg)
	if err != nil {
		return nil, err
	}

	var ir []importRestrictions
	if err := yaml.Unmarshal(config, &ir); err != nil {
		return nil, err
	}

	return ir, nil
}

func getDirDeps(dir string) ([]goPackage, error) {
	cmd := exec.Command("go", "list", "-json", fmt.Sprintf("%s...", dir))
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(errors.Wrap(err, string(stdout)), "go list")
	}

	dec := json.NewDecoder(bytes.NewReader(stdout))
	var packages []goPackage
	for dec.More() {
		var pkg goPackage
		if err := dec.Decode(&pkg); err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

func stringSliceContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
