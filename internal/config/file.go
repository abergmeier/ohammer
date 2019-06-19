package config

import (
	"os"

	cliconfig "github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
)

func DefaultFile() *configfile.ConfigFile {
	return cliconfig.LoadDefaultConfigFile(os.Stderr)
}
