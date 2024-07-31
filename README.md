Multi-plateform SSH tools

# Features

- Read Keepass for ssh password
- Copy files from/to ssh remote host
- Open tunnels from local or remote host
- Open socks or http proxy through remote host

# Install

## From binary

Download appropriate version from [Github release](https://github.com/hurlebouc/sshor/releases/latest).

## From sources

If go 1.22 is installed on your machine, you can do

```sh
go build .
```

# Configure

Sshor currently use its own configuration (It does not read `~/.ssh/config`). This configuration is located at `~/.config/sshor` on Unix host, or at `AppData\Roaming\sshor` for windows.

This location is considered by Sshor as a [Cue package](https://cuelang.org/docs/concept/modules-packages-instances/#packages). This package must be named `sshor`.

Expected structure is an `hosts` map where keys are host names and values are `host` objects with following properties:

* `host`: address of the remote host (if missing uses the last jump host, or local host if none)
* `port`: port of the remote host (22 if missing)
* `user`: SSH user we are connecting to (if missing uses the last jump user, or local user if none)
* `keepass`: keepass location of SSH password (may be missing)
* `jump`: jump host between local and remote host (may be missing)

Field `jump` follows the same structure as `host` objects.

If `keepass` field is present, its value must be an object contaning requested following properties:

* `path`: location of the keepass database
* `id`: location of the password entry in database

Very simple example of configuration is given by the following snippet:

```cuelang
package sshor

hosts: {
    host1: {
        host: "example.com"
        port: 22
        user: "bob"
    }
    host2: {
        host: "my.ssh.host.com"
        port: 22
        user: "alice"
    }
}
```

More complicated example:

```cuelang
package sshor

_machine1: {
	plop: {
		host: "1.2.3.4"
	}
	plip: {
		host: "2.3.4.5"
	}
}

_machine2: {
	plap: {
		host: "8.8.8.8"
	}
	plup: {
		host: "192.168.1.2"
		user: "user"
	}
	testjump: {
		host: "127.0.0.1"
		user: "user"
		jump: plup
	}
}

// hosts whose credentials are stored in keepass database with the same name as the key of the map
hosts: {
	for k, v in _machine1 {
		(k): v
		(k): {
            keepass: {
                path: "/path/to/keepass.kdbx"
                id: "/id/in/keepass/\(k)"
            }
        }
	}
}

hosts: {
	for k, v in _machine2 {
		(k): v
	}
}
```
