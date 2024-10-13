package resources

import (
	"embed"
	_ "embed"
)

//go:embed sql.create/*
var StaticFiles embed.FS
