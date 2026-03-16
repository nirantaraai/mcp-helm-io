package domain

// ReleaseStatus represents the status of a Helm release
type ReleaseStatus string

const (
	// StatusPending indicates the release is pending deployment
	StatusPending ReleaseStatus = "pending"

	// StatusDeployed indicates the release is successfully deployed
	StatusDeployed ReleaseStatus = "deployed"

	// StatusUninstalling indicates the release is being uninstalled
	StatusUninstalling ReleaseStatus = "uninstalling"

	// StatusUninstalled indicates the release has been uninstalled
	StatusUninstalled ReleaseStatus = "uninstalled"

	// StatusSuperseded indicates the release has been superseded by another
	StatusSuperseded ReleaseStatus = "superseded"

	// StatusFailed indicates the release deployment failed
	StatusFailed ReleaseStatus = "failed"

	// StatusPendingInstall indicates the release is pending installation
	StatusPendingInstall ReleaseStatus = "pending-install"

	// StatusPendingUpgrade indicates the release is pending upgrade
	StatusPendingUpgrade ReleaseStatus = "pending-upgrade"

	// StatusPendingRollback indicates the release is pending rollback
	StatusPendingRollback ReleaseStatus = "pending-rollback"

	// StatusUnknown indicates the release status is unknown
	StatusUnknown ReleaseStatus = "unknown"
)

// String returns the string representation of the status
func (s ReleaseStatus) String() string {
	return string(s)
}

// IsValid checks if the status is valid
func (s ReleaseStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusDeployed, StatusUninstalling, StatusUninstalled,
		StatusSuperseded, StatusFailed, StatusPendingInstall, StatusPendingUpgrade,
		StatusPendingRollback, StatusUnknown:
		return true
	default:
		return false
	}
}

// IsTerminal checks if the status is terminal (no further changes expected)
func (s ReleaseStatus) IsTerminal() bool {
	return s == StatusDeployed || s == StatusUninstalled || s == StatusFailed
}

// Made with Bob
