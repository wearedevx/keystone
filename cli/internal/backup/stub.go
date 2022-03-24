package backup

type stubBackupService struct {
}

func (b *stubBackupService) Setup() error {
	return nil
}

func (b *stubBackupService) IsSetup() bool {
	return true
}

func (b *stubBackupService) Backup() (string, string, error) {
	return "", "", nil
}

func (b *stubBackupService) Restore() error {
	return nil
}
