package images

import "fmt"

// Repository maps component versions to their corresponding Docker image URLs
type Repository struct {
	mapping map[string]string
}

const (
	opLabsToolsRegistry = "us-docker.pkg.dev/oplabs-tools-artifacts/images"
	paradigmRegistry    = "ghcr.io/paradigmxyz"
)

// NewRepository creates a new Repository instance with predefined mappings
func NewRepository() *Repository {
	return &Repository{
		mapping: map[string]string{
			// OP Labs images
			"mc-deployer":   opLabsToolsRegistry,
			"op-geth":       opLabsToolsRegistry,
			"mc-node":       opLabsToolsRegistry,
			"mc-batcher":    opLabsToolsRegistry,
			"mc-proposer":   opLabsToolsRegistry,
			"mc-challenger": opLabsToolsRegistry,
			// Paradigm images
			"op-reth": paradigmRegistry,
		},
	}
}

// GetImage returns the full Docker image URL for a given component and version
func (r *Repository) GetImage(component string, version string) string {
	if imageTemplate, ok := r.mapping[component]; ok {

		if version == "" {
			version = "latest"
		}
		return fmt.Sprintf("%s/%s:%s", imageTemplate, component, version)
	}

	// TODO: that's our way to convey that the "default" image should be used.
	// We should probably have a more explicit way to do this.
	return ""
}
