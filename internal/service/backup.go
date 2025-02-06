package service

type backup struct {
	date string
	size string
	name string
}

func NewBackup(date string, size string, name string) *backup {
	return &backup{
		date: date,
		size: size,
		name: name,
	}
}
