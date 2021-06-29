// Copyright 2020 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/varshaprasad96/hybrid-helm-plugin/pkg/hybrid/v1alpha1/scaffolds"
	"sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugin"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugin/util"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugins/golang"
)

// TODO: Implement the apiSubcommand.
type initSubcommand struct {
	config config.Config

	// For help text
	commandName string

	// boilerplate options
	license string
	owner   string

	// go config options
	repo string
}

var _ plugin.InitSubcommand = &initSubcommand{}

// UpdateContext define plugin context
// TODO: modify this
func (p *initSubcommand) UpdateMetadata(cliMeta plugin.CLIMetadata, subcmdMeta *plugin.SubcommandMetadata) {
	subcmdMeta.Description = `Initialize a new project including the following files:
	- a "go.mod" with project dependencies
	- a "PROJECT" file that stores project configuration
	- a "Makefile" with several useful make targets for the project
	- several YAML files for project deployment under the "config" directory
	- a "main.go" file that creates the manager that will run the project controllers
  `
	subcmdMeta.Examples = fmt.Sprintf(`  # Initialize a new project with your domain and name in copyright
	$ %[1]s init --plugins=%[2]s --domain=example.com --owner "Your Name"

	# Initialize a new project defining a specific project version
	%[1]s init --plugins=%[2]s --project-version 3
`, cliMeta.CommandName, pluginKey)

	p.commandName = cliMeta.CommandName
}

// TODO: bind the same set of flags to the apiSubcommand
func (p *initSubcommand) BindFlags(fs *pflag.FlagSet) {
	fs.SortFlags = false

	// project args
	fs.StringVar(&p.repo, "repo", "", "name to use for go module (e.g., github.com/user/repo), "+
		"defaults to the go package of the current working directory.")

	// boilerplate args
	fs.StringVar(&p.license, "license", "apache2",
		"license to use to boilerplate, may be one of 'apache2', 'none'")
	fs.StringVar(&p.owner, "owner", "", "owner to add to the copyright")

}

func (p *initSubcommand) InjectConfig(c config.Config) error {
	p.config = c

	// Try to guess repository if flag is not set
	if p.repo == "" {
		repoPath, err := golang.FindCurrentRepo()
		if err != nil {
			return fmt.Errorf("error finding current repository: %v", err)
		}
		p.repo = repoPath
	}

	if err := p.config.SetRepository(p.repo); err != nil {
		return err
	}
	return nil
}

// TODO: Pre-scaffold check to verify if the right Go version and directory is used
// needs to be added from Kubebuilder.
func (p *initSubcommand) Scaffold(fs machinery.Filesystem) error {
	// TODO: add customizations to config files, as done in helm operator plugin
	scaffolder := scaffolds.NewInitScaffolder(p.config, p.license, p.owner)
	scaffolder.InjectFS(fs)
	err := scaffolder.Scaffold()
	if err != nil {
		return err
	}

	// Ensure that we are pinning the controller-runtime version
	// xref: https://github.com/kubernetes-sigs/kubebuilder/issues/997
	err = util.RunCmd("Get controller runtime", "go", "get",
		"sigs.k8s.io/controller-runtime@"+scaffolds.ControllerRuntimeVersion)
	if err != nil {
		return err
	}
	return nil
}

func (p *initSubcommand) PostScaffold() error {
	err := util.RunCmd("Update dependencies", "go", "mod", "tidy")
	if err != nil {
		return err
	}

	return nil
}
