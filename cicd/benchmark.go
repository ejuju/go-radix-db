package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	fpath := "bench.txt"
	log.Printf("Running benchmarks and saving results to file: %q", fpath)
	f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	cmd := exec.Command("go", "test", "./...", "-bench=.", "-benchmem")
	cmd.Stdout = f
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
