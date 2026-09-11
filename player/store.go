package player

type Store interface {
	Load(name string) (*Record, bool, error)
	Save(r *Record) error
}
