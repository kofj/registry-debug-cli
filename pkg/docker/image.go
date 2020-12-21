package docker

import (
	"encoding/json"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	dockerimage "github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	digest "github.com/opencontainers/go-digest"
)

// BuildConfigBytes ...
func BuildConfigBytes(digests ...digest.Digest) ([]byte, error) {
	var diffids []layer.DiffID
	for _, d := range digests {
		diffids = append(diffids, layer.DiffID(d))
	}
	var cfg = dockerimage.Image{
		V1Image: dockerimage.V1Image{
			Architecture: "amd64",
			OS:           "linux",
			Created:      time.Now().UTC(),
		},
		RootFS: &dockerimage.RootFS{
			Type:    dockerimage.TypeLayers,
			DiffIDs: diffids,
		},
	}

	return json.MarshalIndent(cfg, "  ", "  ")
}

// BuildManifest ...
func BuildManifest(configSize int64, configDigest digest.Digest, blobsDescriptors ...distribution.Descriptor) (distribution.Manifest, error) {
	var manifest = schema2.Manifest{
		Versioned: schema2.SchemaVersion,
		Config: distribution.Descriptor{
			MediaType: schema2.MediaTypeImageConfig,
			Size:      configSize,
			Digest:    configDigest,
		},
		Layers: blobsDescriptors,
	}
	return schema2.FromStruct(manifest)
}
