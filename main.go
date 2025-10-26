package main

import (
	"flag"
	"fmt"
	"os"
	"sojebsikder/i2p/internal/converter"
	"sojebsikder/i2p/internal/util"
)

var version = "0.0.2"

func showUsage() {
	fmt.Println("Usage:")
	fmt.Println("  i2p convert [--input-file FILE] [--output-file FILE]")
	fmt.Println("  i2p help")
	fmt.Println("  i2p version")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --input-file FILE   Specify the input file (default: insomnia.yaml)")
	fmt.Println("  --output-file FILE  Specify the output file (default: postman_collection.json)")
}

func main() {
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "convert":
		inputFile := "insomnia.yaml"
		outputFile := "postman_collection.json"

		fs := flag.NewFlagSet("convert", flag.ExitOnError)
		fs.StringVar(&inputFile, "input-file", inputFile, "Specify the input file")
		fs.StringVar(&outputFile, "output-file", outputFile, "Specify the output file")
		fs.Parse(os.Args[2:])

		inputFileEx := util.GetFileExtensionFromURL(inputFile)

		if inputFileEx == ".yaml" || inputFileEx == ".yml" {
			converter.ConvertInsomniaToPostman(inputFile, outputFile)
		} else if inputFileEx == ".json" {
			converter.ConvertPostmanToInsomnia(inputFile, outputFile)
		} else {
			fmt.Println("Input file format is not supported: " + inputFileEx)
		}

	case "help":
		showUsage()
	case "version":
		fmt.Println("i2p version " + version)
	default:
		fmt.Println("Unknown command:", cmd)
		fmt.Println("Use 'i2p help' to see available commands.")
		os.Exit(1)
	}
}
