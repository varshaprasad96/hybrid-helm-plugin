/*
Copyright 2019 The Kubernetes Authors.
Modifications copyright 2020 The Operator-SDK Authors

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

package scaffolds

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/varshaprasad96/hybrid-helm-plugin/pkg/hybrid/v1alpha1/scaffolds/internal/templates"
	"github.com/varshaprasad96/hybrid-helm-plugin/pkg/hybrid/v1alpha1/scaffolds/internal/templates/hack"
	"sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugins"
)

const (
	// ControllerRuntimeVersion is the kubernetes-sigs/controller-runtime version to be used in the project
	ControllerRuntimeVersion = "v0.8.3"
	// ControllerToolsVersion is the kubernetes-sigs/controller-tools version to be used in the project
	ControllerToolsVersion = "v0.5.0"
)

var _ plugins.Scaffolder = &initScaffolder{}

type initScaffolder struct {
	// fs is the filesystem that will be used by the scaffolder
	fs machinery.Filesystem

	config          config.Config
	boilerplatePath string
	license         string
	owner           string
}

// NewInitScaffolder returns a new plugins.Scaffolder for project initialization operations
func NewInitScaffolder(config config.Config, license, owner string) plugins.Scaffolder {
	return &initScaffolder{
		config:          config,
		boilerplatePath: hack.DefaultBoilerplatePath,
		license:         license,
		owner:           owner,
	}
}

// InjectFS implements Scaffolder
func (s *initScaffolder) InjectFS(fs machinery.Filesystem) {
	s.fs = fs
}

// Scaffold implements scaffolder
func (s *initScaffolder) Scaffold() error {
	fmt.Println("Writing scaffolds for you to edit...")

	// Initialize the machinery.Scaffold that will write the files to disk
	scaffold := machinery.NewScaffold(s.fs,
		machinery.WithDirectoryPermissions(0755),
		machinery.WithFilePermissions(0644),
		machinery.WithConfig(s.config),
	)

	// The boilerplate file needs to be scaffolded as a separate step as it is going to be used
	// by rest of the files, even those scaffolded in this command call.
	bpFile := &hack.Boilerplate{
		License: s.license,
		Owner:   s.owner,
	}

	bpFile.Path = s.boilerplatePath
	if err := scaffold.Execute(bpFile); err != nil {
		return err
	}

	boilerplate, err := afero.ReadFile(s.fs.FS, s.boilerplatePath)
	if err != nil {
		return err
	}

	// Initialize the machinery.Scaffold that will write the files to disk
	scaffold = machinery.NewScaffold(s.fs,
		machinery.WithConfig(s.config),
		machinery.WithBoilerplate(string(boilerplate)),
	)

	return scaffold.Execute(
		&templates.Main{},
		&templates.GitIgnore{},
		&templates.Watches{})
}
