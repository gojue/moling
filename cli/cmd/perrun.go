// Copyright 2025 CFC4N <cfc4n.cs@gmail.com>. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Repository: https://github.com/gojue/moling

package cmd

import (
	"github.com/spf13/cobra"
	"path/filepath"
)

// mlsCommandPreFunc is a pre-run function for the MoLing command.
func mlsCommandPreFunc(cmd *cobra.Command, args []string) error {
	err := CreateDirectory(mlConfig.BasePath)
	if err != nil {
		return err
	}
	for _, dirName := range mlDirectories {
		err = CreateDirectory(filepath.Join(mlConfig.BasePath, dirName))
		if err != nil {
			return err
		}
	}
	return nil
}
