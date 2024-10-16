package resources

import (
	"embed"
)

//go:embed sql.create/*
var StaticFiles embed.FS
