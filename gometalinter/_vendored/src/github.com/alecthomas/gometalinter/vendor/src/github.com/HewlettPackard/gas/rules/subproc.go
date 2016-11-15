// (c) Copyright 2016 Hewlett Packard Enterprise Development LP
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

package rules

import (
	gas "github.com/HewlettPackard/gas/core"
	"go/ast"
	"regexp"
	"strings"
)

type Subprocess struct {
	pattern *regexp.Regexp
}

func (r *Subprocess) Match(n ast.Node, c *gas.Context) (*gas.Issue, error) {
	if node := gas.MatchCall(n, r.pattern); node != nil {
		// call with variable command or arguments
		for _, arg := range node.Args {
			if _, test := arg.(*ast.BasicLit); !test {
				// TODO: try to resolve the symbol ...
				what := "Subprocess launching with variable."
				return gas.NewIssue(c, n, what, gas.High, gas.High), nil
			}
		}

		// call with partially qualified command
		if str, err := gas.GetString(node.Args[0]); err == nil {
			if !strings.HasPrefix(str, "/") {
				what := "Subprocess launching with partial path."
				return gas.NewIssue(c, n, what, gas.Medium, gas.High), nil
			}
		}

		what := "Subprocess launching should be audited."
		return gas.NewIssue(c, n, what, gas.Low, gas.High), nil
	}
	return nil, nil
}

func NewSubproc() (r gas.Rule, n ast.Node) {
	r = &Subprocess{
		pattern: regexp.MustCompile(`^exec.Command$`),
	}
	n = (*ast.CallExpr)(nil)
	return
}
