package main

import (
	"main/db"
	"os"

	"github.com/gofrs/uuid"
)

var storeFile = "uniqId"
var dbInstance = &db.DB{}

func getId() string {
	uid, err := dbInstance.Get([]string{storeFile})
	if os.IsNotExist(err) {
		id, err := uuid.NewV4()
		if err != nil {
			return ""
		}

		dbInstance.Append([]string{storeFile}, []byte(id.String()))
		return id.String()
	}

	if err != nil {
		return ""
	}

	return string(uid)
}
