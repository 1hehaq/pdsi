package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	blue    = "\033[34m"
	red     = "\033[31m"
	green   = "\033[32m"
	greenbg = "\x1b[42m"
	reset   = "\033[0m"
)

var sensitiveKeywords = []string{
	// personal & confidential
	"confidential", "private", "restricted", "internal",
	"not for distribution", "do not share", "proprietary",
	"trade secret", "classified", "sensitive",

	// financial & legal
	"bank statement", "invoice", "salary", "contract",
	"agreement", "nda", "non disclosure",

	// personal data
	"passport", "social security", "ssn", "date of birth",
	"credit card", "tax", "identity", "id number",

	// business critical
	"revenue", "profit", "loss", "financial report",
	"quarterly report", "annual report", "audit",
	"board meeting", "shareholders",

	// company internal
	"internal use", "draft", "preliminary",
	"company confidential", "staff only",
	"management only", "executive", "internal",
}

func highlightKeywords(input string) string {
	for _, keyword := range sensitiveKeywords {
		input = strings.ReplaceAll(input, keyword, fmt.Sprintf("\033[33m%s%s", keyword, reset))
	}
	return input
}

func analyzePDF(url string) {
	fmt.Printf("%s[-]%s %sreading:%s %s\n", blue, reset, blue, reset, url)

	cmd := exec.Command("wget", "--no-check-certificate", "-qO-", url)
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s[!]%s %snot fetching:%s %s\n", red, reset, red, reset, url)
		return
	}

	pdftotextCmd := exec.Command("pdftotext", "-", "-")
	pdftotextCmd.Stdin = bytes.NewReader(output)
	pdfOutput, err := pdftotextCmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s[!]%s %snot processing:%s %s\n", red, reset, red, reset, url)
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(pdfOutput))
	foundSensitive := false
	for scanner.Scan() {
		line := scanner.Text()
		for _, keyword := range sensitiveKeywords {
			if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
				highlighted := highlightKeywords(line)
				fmt.Printf("%s[+]%s %ssensitive:%s %s\n", green, reset, greenbg, reset, highlighted)
				foundSensitive = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s[!]%s %sfailed to scan pdf:%s %s\n", red, reset, red, reset, url)
	}

	if !foundSensitive {
		fmt.Printf("%s[*]%s %snothing found%s\n", blue, reset, blue, reset)
	}
}

func main() {
	pdfFlag := flag.String("pdf", "", "path to a file containing a list of PDF URLs")
	flag.Parse()

	var reader io.Reader

	if *pdfFlag != "" {
		file, err := os.Open(*pdfFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s[!]%s %sfailed to open:%s %v\n", red, reset, red, reset, err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		info, _ := os.Stdin.Stat()
		if info.Mode()&os.ModeCharDevice != 0 {
			fmt.Fprintf(os.Stderr, "%s[!]%s %sno input provided%s\n", red, reset, red, reset)
			os.Exit(1)
		}
		reader = os.Stdin
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		url := scanner.Text()
		if strings.TrimSpace(url) != "" {
			analyzePDF(url)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s[!]%s %serror reading input:%s %v\n", red, reset, red, reset, err)
		os.Exit(1)
	}
}
