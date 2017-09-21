package message

type GetRequest struct {
	Request
	Targets []string
}

type GetResponse struct {
	Response
}
