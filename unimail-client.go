package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	unimail "github.com/unimails/unimail-go-sdk"
)

const VERSION = "0.0.1"
const defaultRepo = "unimails/unimail-client"

type cliOptions struct {
	key       string
	from      string
	receivers []string
	cc        string
	bcc       string
	subject   string
	txt       string
	html      string
	files     []string
}

type stringList []string

func (l *stringList) String() string {
	if l == nil {
		return ""
	}
	return strings.Join(*l, ",")
}

func (l *stringList) Set(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return errors.New("value cannot be empty")
	}
	*l = append(*l, v)
	return nil
}

type releaseInfo struct {
	TagName string `json:"tag_name"`
}

func main() {
	if len(os.Args) == 1 {
		usage()
		return
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			printVersion()
			return
		case "update", "upgrade":
			if err := updateToLatest(defaultRepo); err != nil {
				fmt.Fprintf(os.Stderr, "update failed: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	opts, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	req, err := buildRequest(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid arguments: %v\n", err)
		os.Exit(2)
	}

	client := unimail.New(opts.key)
	result := client.SendEmail(req)
	if !result.IsSucess() {
		fmt.Fprintf(os.Stderr, "send failed: code=%d msg=%s\n", result.Code, result.Msg)
		os.Exit(1)
	}

	fmt.Printf("send success: code=%d msg=%s\n", result.Code, result.Msg)
}

func parseArgs(args []string) (cliOptions, error) {
	var opts cliOptions

	for _, raw := range args {
		if raw == "-h" || raw == "--help" {
			usage()
			os.Exit(0)
		}
	}

	fs := flag.NewFlagSet("unimail-client", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = usage

	var receivers stringList
	var files stringList

	defaultKey := strings.TrimSpace(os.Getenv("UNIMAIL_KEY"))

	fs.StringVar(&opts.key, "k", defaultKey, "Unimail API key (default from env UNIMAIL_KEY)")
	fs.StringVar(&opts.key, "key", defaultKey, "Unimail API key (default from env UNIMAIL_KEY)")

	fs.StringVar(&opts.from, "f", "", "mail from display name")
	fs.StringVar(&opts.from, "form", "", "mail from display name")
	fs.StringVar(&opts.from, "from", "", "mail from display name")

	fs.Var(&receivers, "r", "receiver email, repeat for multiple receivers")
	fs.Var(&receivers, "receiver", "receiver email, repeat for multiple receivers")

	fs.StringVar(&opts.cc, "c", "", "cc emails (comma-separated)")
	fs.StringVar(&opts.cc, "cc", "", "cc emails (comma-separated)")

	fs.StringVar(&opts.bcc, "b", "", "bcc emails (comma-separated)")
	fs.StringVar(&opts.bcc, "bb", "", "bcc emails (comma-separated)")
	fs.StringVar(&opts.bcc, "bcc", "", "bcc emails (comma-separated)")

	fs.StringVar(&opts.subject, "s", "", "mail subject")
	fs.StringVar(&opts.subject, "subject", "", "mail subject")

	fs.StringVar(&opts.txt, "t", "", "plain text content")
	fs.StringVar(&opts.txt, "txt", "", "plain text content")

	fs.StringVar(&opts.html, "H", "", "html content (short flag -H)")
	fs.StringVar(&opts.html, "html", "", "html content")

	fs.Var(&files, "file", "attachment pair list: name;path or name;path;name2;path2")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			usage()
			os.Exit(0)
		}
		return opts, err
	}

	opts.receivers = receivers
	opts.files = files

	if strings.TrimSpace(opts.key) == "" {
		return opts, errors.New("missing key: set UNIMAIL_KEY or provide -k/--key")
	}

	return opts, nil
}

func buildRequest(opts cliOptions) (unimail.UnimailReq, error) {
	req := unimail.UnimailReq{
		From:        strings.TrimSpace(opts.from),
		Receivers:   opts.receivers,
		Cc:          strings.TrimSpace(opts.cc),
		Bcc:         strings.TrimSpace(opts.bcc),
		Subject:     strings.TrimSpace(opts.subject),
		TxtContent:  opts.txt,
		HtmlContent: opts.html,
	}

	if len(req.Receivers) == 0 {
		return req, errors.New("at least one -r/--receiver is required")
	}
	if req.Subject == "" {
		return req, errors.New("-s/--subject is required")
	}
	if strings.TrimSpace(req.TxtContent) == "" && strings.TrimSpace(req.HtmlContent) == "" {
		return req, errors.New("either -t/--txt or --html/-H must be provided")
	}

	for _, item := range opts.files {
		parts := splitAndTrim(item, ";")
		if len(parts)%2 != 0 {
			return req, fmt.Errorf("--file expects even number of ';' separated values, got: %q", item)
		}
		for i := 0; i < len(parts); i += 2 {
			name := parts[i]
			path := parts[i+1]
			if err := req.AppendFile(name, path); err != nil {
				return req, fmt.Errorf("append file failed (name=%s path=%s): %w", name, path, err)
			}
		}
	}

	return req, nil
}

func usage() {
	fmt.Printf(`unimail-client v%s

Usage:
  unimail-client [options]
  unimail-client version
  unimail-client update|upgrade

Options:
  -k, --key         API key (default from env UNIMAIL_KEY)
  -f, --form        set UnimailReq.From (also supports --from)
  -r, --receiver    set UnimailReq.Receivers, repeatable
  -c, --cc          set UnimailReq.Cc
  -b, --bb          set UnimailReq.Bcc (also supports --bcc)
  -s, --subject     set UnimailReq.Subject
  -t, --txt         set UnimailReq.TxtContent
  --html, -H        set UnimailReq.HtmlContent
  --file            attachment pairs separated by ';'
                    e.g. --file "report.pdf;./report.pdf"
                    e.g. --file "a.txt;./a.txt;b.txt;./b.txt"
  -h, --help        show help

Commands:
  version           print current version
  update/upgrade    check latest release of %s and install latest if outdated

Examples:
  unimail-client -k xxx -f "Sender" -r a@x.com -r b@x.com -s "Hello" -t "Plain"
  unimail-client --form "Sender" --receiver a@x.com --subject "Hi" --html "<b>Hi</b>"
  unimail-client --file "log.txt;./log.txt" -r a@x.com -s "with file" -t "ok"
`, VERSION, defaultRepo)
}

func printVersion() {
	fmt.Printf("unimail-client version %s\n", VERSION)
}

func updateToLatest(repo string) error {
	latest, err := fetchLatestReleaseTag(repo)
	if err != nil {
		return err
	}

	if compareSemver(VERSION, latest) >= 0 {
		fmt.Printf("already up to date: current=%s latest=%s\n", VERSION, latest)
		return nil
	}

	module := "github.com/" + repo
	fmt.Printf("updating %s from %s to %s ...\n", module, VERSION, latest)
	cmd := exec.Command("go", "install", module+"@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go install failed: %w", err)
	}

	fmt.Println("update completed")
	return nil
}

func fetchLatestReleaseTag(repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "unimail-client")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github releases api returned status %d", resp.StatusCode)
	}

	var info releaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", err
	}
	if strings.TrimSpace(info.TagName) == "" {
		return "", errors.New("latest release tag is empty")
	}

	return strings.TrimPrefix(strings.TrimSpace(info.TagName), "v"), nil
}

func splitAndTrim(s string, sep string) []string {
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}

func compareSemver(a string, b string) int {
	aParts := parseVersion(a)
	bParts := parseVersion(b)
	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}
	for i := 0; i < maxLen; i++ {
		av := 0
		bv := 0
		if i < len(aParts) {
			av = aParts[i]
		}
		if i < len(bParts) {
			bv = bParts[i]
		}
		if av > bv {
			return 1
		}
		if av < bv {
			return -1
		}
	}
	return 0
}

func parseVersion(v string) []int {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	if v == "" {
		return []int{0}
	}
	parts := strings.Split(v, ".")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			out = append(out, 0)
			continue
		}
		out = append(out, n)
	}
	return out
}
