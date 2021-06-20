package ops

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/onprem"
	"github.com/nanovms/ops/types"
)

func CurrentVersion() (string, error) {
	var err error

	local, remote := lepton.LocalReleaseVersion, lepton.LatestReleaseVersion
	if local == "0.0" {
		err = lepton.DownloadReleaseImages(remote)
		if err != nil {
			return "", err
		}
		return remote, nil
	}

	return local, nil
}

func ReadConfigFromFile(file string) (c *types.Config, err error) {
	c = &types.Config{}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config: %v\n", err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		err = fmt.Errorf("error config: %v", err)
		return
	}

	c.VolumesDir = lepton.LocalVolumeDir
	if c.Mounts != nil {
		err = onprem.AddMountsFromConfig(c)
	}

	return
}
