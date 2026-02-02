package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
    // --- 1. Validate CLI arguments ---
    if len(os.Args) != 3 {
        fmt.Fprintln(os.Stderr, "Usage: go run . <inputFile> <outputFile>")
        os.Exit(1)
    }

    inputPath := os.Args[1]
    outputPath := os.Args[2]

    // --- 2. Open input file ---
    in, err := os.Open(inputPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer in.Close()

    // --- 3. Create output file ---
    out, err := os.Create(outputPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer out.Close()

    scanner := bufio.NewScanner(in)

    // --- 4. Process line-by-line ---
    for scanner.Scan() {
        line := scanner.Text()

        // The FSM engine (defined elsewhere)
        result := ProcessLineFSM(line)

        _, err := fmt.Fprintln(out, result)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
    }

    // --- 5. Handle scanning errors ---
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
