package common

type ApiReturn struct {
	Ret int `json:"ret"`
	Msg string `json:"msg,omitempty"`
	Data *ApiData `json:"data,omitempty"`
}

type ApiData struct {
	Ext interface{} `json:"ent,omitempty"`
	Ent interface{} `json:"ext,omitempty"`
}
