package ops

import (
	"fmt"

	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/provider"
	"github.com/nanovms/ops/types"
)

// ProviderByType returns provider identified by given type.
func ProviderByType(typeName string, config *types.Config) (lepton.Provider, error) {
	if config == nil {
		config = &types.Config{}
	}

	if typeName == "" {
		typeName = "onprem"
	}

	provider, err := provider.CloudProvider(typeName, &(config.CloudConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}
	return provider, nil
}
