package inspectors

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/keyval-dev/odigos/common"
	procdiscovery "github.com/keyval-dev/odigos/procdiscovery/pkg/process"
)

type dotnetInspector struct{}

const (
	aspnet = "ASPNET"
	dotnet = "DOTNET"
)

var dotNet = &dotnetInspector{}

func (d *dotnetInspector) Inspect(p *procdiscovery.Details) (common.ProgrammingLanguage, bool) {
	data, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/environ", p.ProcessID))
	if err == nil {
		environ := string(data)
		if strings.Contains(environ, aspnet) || strings.Contains(environ, dotnet) {
			return common.DotNetProgrammingLanguage, true
		}
	}

	return "", false
}
