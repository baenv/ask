package command

// ICommand defines the interface for command operations
type ICommand interface {
	RegisterReg()
	AddHandler()
	RegisterLs()
	RegisterAi()
	RegisterStart()
	RegisterSum()
}
