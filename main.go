package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/CodeMonk/GoSquaretoShopify/squarespace"
)

var (
	verbose = flag.Bool("verbose", false, "Show verbose logging to stderr")
	output  = flag.String("output", "-", "Choose output file (default: stdout)")
)

func init() {
	flag.Parse()
}

// getFileHandles uses our arguments to set up our input and output files
func getFileHandles() (in, out *os.File, err error) {
	// Input file should be our only command line argument
	if len(flag.Args()) != 1 {
		err = fmt.Errorf("Must have exactly one argument (input filename).  Received: %v",
			flag.Args())
		return
	}
	in, err = os.Open(flag.Arg(0))
	if err != nil {
		err = fmt.Errorf("Error: Unable to open input file %v: %v", flag.Arg(0),
			err)
		return
	}

	// Output file
	out = os.Stdout
	if *output != "-" {
		out, err = os.Create(*output)
		if err != nil {
			err = fmt.Errorf("Error: Unable to open output file %v: %v", *output,
				err)
			return
		}
	}

	return
}

func convertFiles(in, out *os.File) error {

	square, err := squarespace.New(in)
	if err != nil {
		return fmt.Errorf("Unable to read input file!: %v", err)
	}

	fmt.Fprintf(os.Stderr, "DEBUG:  Squarespace: %#v", square)

	return errors.New("Not implemented yet!")
}

func main() {

	inputFd, outputFd, err := getFileHandles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening files: %s\n", err)
		return
	}

	if *verbose {
		fmt.Fprintln(os.Stderr, "GoSquareToShopify:")
		fmt.Fprintf(os.Stderr, "             Verbose : %v\n", *verbose)
		fmt.Fprintf(os.Stderr, "     Input (File/Fd) : %v / %v \n", flag.Arg(0), inputFd)
		fmt.Fprintf(os.Stderr, "    Output (File/Fd) : %v / %v \n", *output, outputFd)
	}

	err = convertFiles(inputFd, outputFd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting files: %s\n", err)
		return
	}

	fmt.Fprintln(os.Stderr, "Done.")
}
