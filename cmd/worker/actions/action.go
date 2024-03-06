package actions

type Action interface {
	Execute(parameters ActionContext) string
}
