package player

type Store interface {
	Load(name string) (*Record, bool, error)
	Create(r *Record) (*Record, error)
	Save(r *Record) error
}
