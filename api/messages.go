package api

type TxInfo struct {
	Content interface{} `json:"content"`
	Id      string      `json:"id"`
	Type    string      `json:"type"`
}

func NewTxInfo(content interface{}, id, _type string) *TxInfo {
	return &TxInfo{
		Content: content,
		Id:      id,
		Type:    _type,
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
