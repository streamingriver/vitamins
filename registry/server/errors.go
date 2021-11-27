package registry

type RegistryError struct {
	err      string
	notFound bool
	expired  bool
}

func (re *RegistryError) Error() string {
	return re.err
}

func (re *RegistryError) NotFound() bool {
	return re.notFound
}

func (re *RegistryError) Expired() bool {
	return re.expired
}
