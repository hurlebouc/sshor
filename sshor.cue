package sshor

_machine1: {
	plop: {
		host: "1.2.3.4"
		jump: plip
	}
	plip: {
		ip: "2.3.4.5"
	}
}

_machine2: {
	plap: {
		ip: "8.8.8.8"
	}
	kiwi: {
		host: "192.168.1.2"
		user: "partage"
		shellHook: {
			cmd: "su -"
		}
	}
	testjump: {
		host: "127.0.0.1"
		user: "partage"
		jump: kiwi
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
	}
}
