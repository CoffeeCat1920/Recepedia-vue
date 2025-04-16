package database

import "errors"

var ErrItemAlreadyExists = errors.New("Item Already Exists")
var ErrItemNotFound = errors.New("Item Not Found")
var ErrItemMismatch = errors.New("Item Mismatch")
var ErrInternalDatabaseError = errors.New("Internal Database Error")
