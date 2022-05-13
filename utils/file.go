package utils

import (
	"io/ioutil"
)

func ReadFile(file string) []byte {

	data, err := ioutil.ReadFile(file)

	if err != nil {
		panic(err)
	}

	return data

}
