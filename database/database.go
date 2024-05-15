package database

import "fmt"

type D struct{}

func (d D) Get(input []byte) (string, error) {
	return fmt.Sprintf(`{"result": "%s"}`, string(input)), nil
}
