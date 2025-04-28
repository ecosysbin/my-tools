package v1

type VCluster struct {
	Name    string `json:"name"   binding:"required"`
	Version string `json:"version"   binding:"required"`
}
