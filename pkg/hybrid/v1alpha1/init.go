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
	"github.com/spf13/pflag"
	"github.com/varshaprasad96/hybrid-helm-plugin/pkg/hybrid/v1alpha1/scaffolds"
	"sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugin"
)

const (
	groupFlag   = "group"
	versionFlag = "version"
	kindFlag    = "kind"
)

// TODO: Implement the apiSubcommand.
type initSubcommand struct {
	config config.Config

	// For help text
	commandName string

	// Flags
	group   string
	version string
	kind    string
}

var _ plugin.InitSubcommand = &initSubcommand{}

// TODO: bind the same set of flags to the apiSubcommand
func (p *initSubcommand) BindFlags(fs *pflag.FlagSet) {
	fs.SortFlags = false
	fs.StringVar(&p.group, groupFlag, "", "resource Group")
	fs.StringVar(&p.version, versionFlag, "", "resource Version")
	fs.StringVar(&p.kind, kindFlag, "", "resource Kind")
}

func (p *initSubcommand) InjectConfig(c config.Config) error {
	p.config = c
	return nil
}

func (p *initSubcommand) Scaffold(fs machinery.Filesystem) error {
	// TODO: add customizations to config files, as done in helm operator plugin
	scaffolder := scaffolds.NewInitScaffolder(p.config)
	scaffolder.InjectFS(fs)
	return scaffolder.Scaffold()
}
