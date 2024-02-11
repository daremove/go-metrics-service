package healthcheck

type HealthCheck struct {
	storage Storage
}

type Storage interface {
	Ping() error
}

func New(storage Storage) *HealthCheck {
	return &HealthCheck{
		storage,
	}
}

func (hc *HealthCheck) CheckStorageConnection() error {
	return hc.storage.Ping()
}
