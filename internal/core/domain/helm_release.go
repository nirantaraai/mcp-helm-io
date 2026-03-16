package domain

import "time"

// HelmRelease represents a Helm release in the domain
type HelmRelease struct {
	Name      string
	Namespace string
	Chart     string
	Version   string
	Values    map[string]interface{}
	Status    ReleaseStatus
	Revision  int
	UpdatedAt time.Time
	CreatedAt time.Time
}

// NewHelmRelease creates a new HelmRelease instance
func NewHelmRelease(name, namespace, chart, version string, values map[string]interface{}) *HelmRelease {
	now := time.Now()
	return &HelmRelease{
		Name:      name,
		Namespace: namespace,
		Chart:     chart,
		Version:   version,
		Values:    values,
		Status:    StatusPending,
		Revision:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateStatus updates the release status
func (r *HelmRelease) UpdateStatus(status ReleaseStatus) {
	r.Status = status
	r.UpdatedAt = time.Now()
}

// IncrementRevision increments the revision number
func (r *HelmRelease) IncrementRevision() {
	r.Revision++
	r.UpdatedAt = time.Now()
}

// IsDeployed checks if the release is deployed
func (r *HelmRelease) IsDeployed() bool {
	return r.Status == StatusDeployed
}

// IsFailed checks if the release has failed
func (r *HelmRelease) IsFailed() bool {
	return r.Status == StatusFailed
}

// Made with Bob
