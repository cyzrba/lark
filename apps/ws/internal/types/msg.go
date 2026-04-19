package types

type ReadMsg struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

type WriteMsg struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}