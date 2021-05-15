package cmd

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/jakubd/ipd"
	"github.com/spf13/cobra"
)

var pipeCmd = &cobra.Command{
	Use:   "pipe",
	Short: "get info on an input list piped in via <stdin>",
	Long:  "get info on an input list piped in via <stdin>",
	Run: func(cmd *cobra.Command, args []string) {
		LookupPipe()
	},
}

func LookupPipe() {
	_, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	var runeInput []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		runeInput = append(runeInput, input)
	}

	var thisLine string
	for j := 0; j < len(runeInput); j++ {
		thisLine = thisLine + string(runeInput[j])
		if runeInput[j] == '\n' {
			if len(thisLine) > 1 {
				thisLine = strings.TrimSuffix(thisLine, "\n")
				ipd.OutputLookup(thisLine, false)
				thisLine = ""
			}
		}
	}
}
