package config

type Host struct {
	Host      string `json:"host"`
	User      string `json:"user"`
	Keepass   string `json:"keepass"`
	KeepassId string `json:"keepassId"`
	Port      uint16 `json:"port"`
	Jump      *Host  `json:"jump"`
}

type Config struct {
	Hosts map[string]Host `json:"hosts"`
}

func (c *Config) GetHost(host string) *Host {
	if c == nil {
		return nil
	}
	for k, v := range c.Hosts {
		if k == host {
			return &v
		}
	}
	return nil
}

func (h *Host) GetKeepass() *string {
	if h == nil {
		return nil
	}
	if h.Keepass == "" {
		return nil
	}
	return &h.Keepass
}

func (h *Host) GetKeepassId() *string {
	if h == nil {
		return nil
	}
	if h.KeepassId == "" {
		return nil
	}
	return &h.KeepassId
}

func (h *Host) GetUser() *string {
	if h == nil {
		return nil
	}
	if h.User == "" {
		return nil
	}
	return &h.User
}

func (h *Host) GetHost() *string {
	if h == nil {
		return nil
	}
	if h.Host == "" {
		return nil
	}
	return &h.Host
}

func (h *Host) GetPort() *uint16 {
	if h == nil {
		return nil
	}
	if h.Port == 0 {
		return nil
	}
	return &h.Port
}

func (h Host) GetPortOrDefault(v uint16) uint16 {
	port := h.GetPort()
	if port == nil {
		return v
	} else {
		return *port
	}
}

func (h *Host) GetJump() *Host {
	if h == nil {
		return nil
	}
	return h.Jump
}
