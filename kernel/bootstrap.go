package kernel

func (k *Kernel) Bootstrap() error {
	var err error
	for _, b := range k.bootstrap {
		err = b()
		if err != nil {
			return err
		}
	}
	return nil
}
