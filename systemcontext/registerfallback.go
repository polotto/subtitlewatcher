// +build linux darwin

package systemcontext

import (
	"fmt"
)

func addRegister(keyName string, executablePath string, rightClickText string) error {
	return fmt.Errorf("you are not in Windows")
}
