package config

type Host struct {
	User       string `json:"user"`
	Keepass    string `json:"keepass"`
	KeepassPwd string `json:"keepassPwd"`
	KeepassId  string `json:"keepassId"`
}

type Config struct {
	Hosts map[string]Host `json:"hosts"`
}
