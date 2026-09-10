package telnet

import "fmt"

func failureText(resultCode string) string {
	return fmt.Sprintf("Failure: %v\n", resultCode)
}
