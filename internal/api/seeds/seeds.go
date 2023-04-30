package seeds

type Seed interface {
	Run()
}

type Seeds []Seed

func (s Seeds) Run() {
	for _, seed := range s {
		seed.Run()
	}
}

func NewSeeds(
	superAdminSeed SuperAdminSeed,
) Seeds {
	return Seeds{
		superAdminSeed,
	}
}
