package hsbas

import (
	"io/ioutil"
	"math/rand"
	"strings"
)

const (
	dbusPath    = "/var/lib/dbus/machine-id"
	dbusPathEtc = "/etc/machine-id"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func readFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func trim(s string) string {
	return strings.TrimSpace(strings.Trim(s, "\n"))
}

func machineID() (string, error) {
	id, err := readFile(dbusPath)
	if err != nil {
		// try fallback path
		id, err = readFile(dbusPathEtc)
	}
	if err != nil {
		return "", err
	}
	return trim(string(id)), nil
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
