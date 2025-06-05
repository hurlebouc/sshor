package ssh

import (
	"testing"
)

func TestKeepassXC(t *testing.T) {
	password := ReadKeepass("test/keepassxc.kdbx", "protected", "/plop", "toto")
	if password != "M6mwwAvNgwvKo3d7Twz7b4GdnzZKh5XV" {
		panic("Expected password 'M6mwwAvNgwvKo3d7Twz7b4GdnzZKh5XV', got: " + password)
	}

	password = ReadKeepass("test/keepassxc.kdbx", "protected", "/Other/plap", "tata")
	if password != "bAFwEpxtnxU7tv1nm2Jp412tXVRysgLE" {
		panic("Expected password 'bAFwEpxtnxU7tv1nm2Jp412tXVRysgLE', got: " + password)
	}
}
