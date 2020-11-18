//+build windows

package systemcontext

import (
	"golang.org/x/sys/windows/registry"
	"log"
	"strings"
)

func addRegister(keyName string, executablePath string, rightClickText string) error {
	executablePathConverted := strings.Replace(executablePath, "/", "\\", -1)

	k, err := registry.OpenKey(registry.CURRENT_USER, `Computer\HKEY_CLASSES_ROOT\*\shell\`+keyName, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = k.SetStringValue("Icon", `"`+executablePathConverted+`"`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = k.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = k.SetStringValue("Default", rightClickText)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = k.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}

	k, err = registry.OpenKey(registry.CURRENT_USER, `Computer\HKEY_CLASSES_ROOT\*\shell\`+keyName+`\command`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = k.SetStringValue("Default", `"`+executablePathConverted+`" "%1"`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = k.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
