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

package e2e

import (
	"os"
	"testing"

	"gotest.tools/v3/icmd"
)

func TestMain(m *testing.M) {
	_ = os.Chdir("..")
	os.Exit(m.Run())
}

func TestGood(t *testing.T) {
	cmd := icmd.Command("import-restrictions", "--configuration", "./e2e/testdata/config/ok.yaml")
	result := icmd.RunCmd(cmd)
	result.Assert(t, icmd.Expected{
		ExitCode: 0,
	})
}

func TestBad(t *testing.T) {
	cmd := icmd.Command("import-restrictions", "--configuration", "./e2e/testdata/config/bad.yaml")
	result := icmd.RunCmd(cmd)
	result.Assert(t, icmd.Expected{
		ExitCode: 1,
		Err:      `* Forbidden import "bytes" in package github.com/docker/import-restrictions`,
	})
}
