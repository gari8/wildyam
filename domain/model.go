package domain

const (
	Run = 1<<iota + 1
	Help
	Guide
)

type Recepter struct {
	SubCmd uint
}
