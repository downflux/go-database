package database

import (
	"github.com/downflux/go-database/database/cache"
)

var (
	_ RO = &DB{}
	_ RO = &cache.DB{}
)
