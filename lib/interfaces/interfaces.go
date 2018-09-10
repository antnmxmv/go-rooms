package interfaces

type Game interface {
	// must return game name
	GetName() string
	// must return 1 or 2
	GetWinner() int

	Initialize()
	// should return bool for different specific goals and error != nil if parameters are unacceptable
	Turn(...int) (bool, error)

	GetGrid() [][]int
}