package bs

type StreamerCheckLoginRsp struct {
	ID      int64  `json:"room_id"`
	Name    string `json:"name"`
	Account string `json:"account_name"`
}

type StreamerLoginReq struct {
	Account  string `json:"name"`
	Password string `json:"password"`
}
