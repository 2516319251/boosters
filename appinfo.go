package boosters

type AppInfo interface {
	ID() string
	Name() string
	Version() string
	Endpoint() []string
}

func (boosters *Boosters) ID() string {
	return boosters.opts.id
}

func (boosters *Boosters) Name() string {
	return boosters.opts.name
}

func (boosters *Boosters) Version() string {
	return boosters.opts.version
}

func (boosters *Boosters) Endpoint() []string {
	if boosters.instance != nil {
		return boosters.instance.Endpoints
	}
	return nil
}
