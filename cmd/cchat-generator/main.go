package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
	"github.com/diamondburned/cchat/repository"
)

const OutputDir = "."

func init() {
	log.SetFlags(0)
}

func main() {
	for pkgPath, pkg := range repository.Main {
		g := generate(pkgPath, pkg)

		var destDir = filepath.FromSlash(trimPrefix(repository.RootPath, pkgPath))
		var destFle = filepath.Base(pkgPath)

		// Guarantee that the directory exists.
		if destDir != "" {
			if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
				log.Fatalln("Failed to mkdir -p:", err)
			}
		}

		f, err := os.Create(filepath.Join(destDir, destFle+".go"))
		if err != nil {
			log.Fatalln("Failed to create output file:", err)
		}

		if err := g.Render(f); err != nil {
			log.Fatalln("Failed to render output:", err)
		}

		f.Close()
	}
}

func trimPrefix(rootPrefix, path string) string {
	return strings.Trim(strings.TrimPrefix(path, rootPrefix), "/")
}

// recvName is used to get the receiver variable name. It returns the first
// letter lower-cased. It does NOT do length checking. It only works with ASCII.
func recvName(name string) string {
	return string(unicode.ToLower(rune(name[0])))
}

func generate(pkgPath string, repo repository.Package) *jen.File {
	gen := jen.NewFilePath(pkgPath)
	gen.HeaderComment("DO NOT EDIT: THIS FILE IS GENERATED!")
	gen.PackageComment(repo.Comment.GoString(1))
	gen.Add(generateTypeAlises(repo.TypeAliases))
	gen.Line()
	gen.Add(generateEnums(repo.Enums))
	gen.Line()
	gen.Add(generateStructs(repo.Structs))
	gen.Line()
	gen.Add(generateErrorStructs(repo.ErrorStructs))
	gen.Line()
	gen.Add(generateInterfaces(repo.Interfaces))
	gen.Line()

	return gen
}
