//go:build !windows

package ssh

func adaptConsole(_ int) error {
	return nil
}
