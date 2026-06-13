package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
)

type flags struct {
	Input  string
	Output string
	Prefix string
	Suffix string
	Size   int
}

var config flags

func main() {
	flag.StringVar(&config.Input, "input", "", "")
	flag.StringVar(&config.Output, "output", "", "")
	flag.StringVar(&config.Prefix, "prefix", "", "")
	flag.StringVar(&config.Suffix, "suffix", "", "")
	flag.IntVar(&config.Size, "size", 512, "")

	flag.Parse()

	if config.Input == "" || config.Output == "" {
		slog.Error("input and output must be set")
		os.Exit(1)
	}

	err := os.MkdirAll(config.Output, 0755)
	if err != nil {
		slog.Error("make output directory", "error", err)
		os.Exit(1)
	}

	entries, err := os.ReadDir(config.Input)
	if err != nil {
		slog.Error("read input dir", "error", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		walk(config.Input, entry)
	}
}

func walk(dir string, entry os.DirEntry) {
	fullPath := path.Join(dir, entry.Name())

	if entry.IsDir() {
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			slog.Error("read input dir", "error", err)
			os.Exit(1)
		}

		for _, entry := range entries {
			walk(fullPath, entry)
		}

		return
	}

	hashName := config.Prefix + entry.Name() + config.Suffix
	hashPath := path.Join(dir, hashName)

	hashed := hash(hashPath)
	output := config.Output + "/" + hashed

	generate(fullPath, output, config.Size)
}

func hash(path string) string {
	sum := md5.Sum([]byte(path))
	str := fmt.Sprintf("%x", sum)

	return str
}

func generate(in, out string, width int) {
	cmd := exec.Command("magick", "convert", in, "-geometry", fmt.Sprintf("%dx", width), out)

	slog.Info("generate", "cmd", cmd.String())

	if err := cmd.Run(); err != nil {
		slog.Error("run command", "error", err)
		os.Exit(1)
	}
}
