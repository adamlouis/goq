package apiserver

func NewAPIHandler() APIHandler {
	return &hdl{}
}

type hdl struct{}
