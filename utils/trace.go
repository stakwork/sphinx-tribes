package utils

import (
	"bytes"
	"log"
	"net/http"
	"regexp"
	"runtime/trace"
	"strings"
)

func TraceWithLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		// Start trace capturing
		if err := trace.Start(&buf); err != nil {
			http.Error(w, "Failed to start trace", http.StatusInternalServerError)
			log.Fatalf("Failed to start trace: %v", err)
			return
		}
		defer func() {
			trace.Stop()

			// Extract and process trace output
			traceOutput := buf.String()
			//log.Println("Raw Trace Output (for debugging):", traceOutput)
			// Split the trace output into lines
			traceLines := strings.Split(traceOutput, "\n")
			startedFormatting := false

			// Log each trace line with a line number and format
			for _, line := range traceLines {
				if !startedFormatting && strings.Contains(line, "start trace") {
					startedFormatting = true
				} else {
					continue
				}

				if line != "" {
				  traceFormatted := formatTrace(line)
					for j, formattedLine := range traceFormatted {
								log.Printf("%03d %s", j+1, formattedLine)
				}
				}
			}
		}()

		next(w, r)
	}
}

func splitTrace(traceString string) []string {
    // Define the regex pattern for detecting file paths ending with .go
    re := regexp.MustCompile(`/([^/\s]+\.go)`)

    // Replace each match with the file path followed by a newline
    formattedTrace := re.ReplaceAllString(traceString, "$0\n")

		cleanRegex := regexp.MustCompile(`[^a-zA-Z0-9\s\./:_\-@()]`)
    cleanedTrace := cleanRegex.ReplaceAllString(formattedTrace, "")


    // Split the formatted trace into lines and return as a string array
    return strings.Split(strings.TrimSpace(cleanedTrace), "\n")
}
func formatTrace(input string) []string {
    // Find the position where "start trace" appears
    startIdx := strings.Index(input, "start trace")
    if startIdx == -1 {
        return nil // "start trace" not found, return nil slice
    }

    // Slice input string from "start trace" onwards
    traceContent := input[startIdx:]

    // Use splitTrace to split and format the trace content
    return splitTrace(traceContent)
}


