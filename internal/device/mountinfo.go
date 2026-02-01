package device

import (
	"fmt"

	"github.com/moby/sys/mountinfo"
	"github.com/rs/zerolog/log"
)

// MountEntry represents a single mount point with its metadata.
type MountEntry struct {
	// MountPoint is the path where the filesystem is mounted.
	MountPoint string
	// FSType is the filesystem type (e.g. "ext4", "tmpfs").
	FSType string
	// Source is the device or source of the mount (e.g. "/dev/sda1").
	Source string
	// Options is a comma-separated list of mount options.
	Options string
}

// MountInfoProvider abstracts the retrieval of system mount information.
type MountInfoProvider interface {
	// GetMounts returns all current mount entries.
	GetMounts() ([]MountEntry, error)
}

// SystemMountInfoProvider implements MountInfoProvider using moby/sys/mountinfo.
type SystemMountInfoProvider struct{}

// NewSystemMountInfoProvider creates a new SystemMountInfoProvider.
func NewSystemMountInfoProvider() *SystemMountInfoProvider {
	return &SystemMountInfoProvider{}
}

// GetMounts reads /proc/self/mountinfo and returns all mount entries.
func (p *SystemMountInfoProvider) GetMounts() ([]MountEntry, error) {
	log.Debug().Msg("reading mount info from /proc/self/mountinfo")

	mounts, err := mountinfo.GetMounts(nil)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get mount info")
		return nil, fmt.Errorf("failed to get mount info: %w", err)
	}

	entries := make([]MountEntry, 0, len(mounts))
	for _, m := range mounts {
		entries = append(entries, MountEntry{
			MountPoint: m.Mountpoint,
			FSType:     m.FSType,
			Source:     m.Source,
			Options:    m.Options,
		})
	}

	log.Debug().Int("count", len(entries)).Msg("mount entries loaded")
	return entries, nil
}
