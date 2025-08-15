package ui

import "fmt"

func WifiPassword(network, password string) string {
	return fmt.Sprintf("Сеть %s, пароль: %s", network, password)
}
