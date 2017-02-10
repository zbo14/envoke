package api

type ActionInfo struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

func NewActionInfo(id, _type string) *ActionInfo {
	return &ActionInfo{
		Id:   id,
		Type: _type,
	}
}

type QueryResult struct {
	Log string `json:"log"`
	Ok  bool   `json:"ok"`
}

func NewQueryResult(log string, ok bool) *QueryResult {
	return &QueryResult{
		Log: log,
		Ok:  ok,
	}
}

type UserInfo struct {
	Id      string `json:"user_id"`
	PrivKey string `json:"private_key"`
	PubKey  string `json:"public_key"`
}

func NewUserInfo(id, priv, pub string) *UserInfo {
	return &UserInfo{
		Id:      id,
		PrivKey: priv,
		PubKey:  pub,
	}
}
