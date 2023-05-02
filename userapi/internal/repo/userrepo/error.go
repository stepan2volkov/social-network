package userrepo

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

var (
	errNotUnique = errors.New("unique violation")
)

func mapError(err error) error {
	merr, ok := err.(*mysql.MySQLError)
	if !ok {
		return err
	}

	switch merr.Number {
	case 1062:
		return errNotUnique
	}

	return err
}
