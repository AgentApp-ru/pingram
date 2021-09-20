package healthcheck

type (
	Checker interface {
		Check() bool
	}

	HealthChecker interface {
		Update()
		AddChecker(Checker)
		CheckAll()
		// GetFailed()
	}
)
