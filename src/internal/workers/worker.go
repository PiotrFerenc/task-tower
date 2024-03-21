package workers

type Worker interface {
	Run(address, port string)
}
