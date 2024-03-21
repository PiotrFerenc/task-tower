package controllers

type Controller interface {
	Run(address, port string) error
}
