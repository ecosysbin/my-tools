package repository

const (
	formatTimeLayout = "2006-01-02 15:04:05"
)

const (
	// kubeconfig 的过期时间，单位为秒，这里设置为 10 年
	kubeconfigExpirationSeconds = int64(10 * 365 * 24 * 60 * 60)

	quotaHard = "hard"
	quotaUsed = "used"
)
