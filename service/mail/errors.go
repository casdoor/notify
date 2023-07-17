package mail

// noop at the moment
func asNotifyError(err error) error {
	if err == nil {
		return nil
	}

	return err
}
