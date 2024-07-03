package sshor

_machine1: {
	plop: {
		ip: "1.2.3.4"
	}
	plip: {
		ip: "2.3.4.5"
	}
}

_machine2: {
	plap: {
		ip: "8.8.8.8"
	}
}

// _machine1: [Name=_]: "keepass-access": "/chemin/vers/\(Name)"

hosts: {
	for k, v in _machine1 {
		(k): v
		(k): {keepassAccess: "/chemin/vers/\(k)"}
	}
}
hosts: {
	for k, v in _machine2 {
		(k): v
		(k): {aaaa: v.ip}
	}
}
