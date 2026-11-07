package network

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var dirHeader = regexp.MustCompile(`(?i)^\s*Directory of (.+?)\s*$`)
var dirRow = regexp.MustCompile(`(?i)^\s*(\d{2,4}[-/.]\d{2}[-/.]\d{2,4})\s+(\d{1,2}:\d{2}(?::\d{2})?(?:\s*[AP]M)?)\s+(<[^>]+>|[0-9][0-9,.]*)\s+(.+?)\s*$`)
var dirEmpty = regexp.MustCompile(`(?i)\b0\s+File\(s\)`)

func dirCommand(path string) (string, error) {
	// CMD expands percent/bang variables even inside quotes. Reject those and
	// literal quotes/control characters rather than interpreting a path as code.
	if strings.ContainsAny(path, "\"%!\r\n\x00*?") || strings.IndexFunc(path, func(r rune) bool { return r < 32 || r == 127 }) >= 0 {
		return "", fmt.Errorf("dir browsing paths cannot contain quotes, percent, exclamation, wildcards, or control characters; use native listing for those names")
	}
	if path == "" {
		path = "."
	}
	return `dir /a /-c /4 /t:w /n /-s /-p /-w /-d /-b /-q /-r /-x "` + path + `" 2>&1`, nil
}

func windowsParent(path string) string {
	path = strings.ReplaceAll(path, "/", `\`)
	trimmed := strings.TrimRight(path, `\`)
	if len(trimmed) == 2 && trimmed[1] == ':' {
		return trimmed + `\`
	}
	if strings.HasPrefix(trimmed, `\\`) && len(strings.Split(strings.TrimPrefix(trimmed, `\\`), `\`)) <= 2 {
		return path
	}
	at := strings.LastIndex(trimmed, `\`)
	if at < 0 {
		return path
	}
	parent := trimmed[:at]
	if len(parent) == 2 && parent[1] == ':' {
		parent += `\`
	}
	return parent
}

func parseDirOutput(output string) (*DirectoryListing, error) {
	listing := &DirectoryListing{Entries: []DirectoryEntry{}}
	recognized := 0
	for _, line := range strings.Split(strings.TrimRight(output, "\x00"), "\n") {
		if match := dirHeader.FindStringSubmatch(strings.TrimSpace(line)); match != nil {
			listing.Path = strings.TrimSpace(match[1])
			listing.Parent = windowsParent(listing.Path)
			continue
		}
		match := dirRow.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		recognized++
		if listing.Path == "" {
			return nil, fmt.Errorf("dir output is missing its directory header")
		}
		name := match[4]
		if name == "." || name == ".." {
			continue
		}
		if len(listing.Entries) == 2000 {
			listing.Truncated = true
			break
		}
		kind := strings.ToUpper(match[3])
		isLink := kind == "<JUNCTION>" || strings.HasPrefix(kind, "<SYMLINK")
		if isLink {
			if at := strings.LastIndex(name, " ["); at >= 0 && strings.HasSuffix(name, "]") {
				name = name[:at]
			}
		}
		entry := DirectoryEntry{Name: name, Path: strings.TrimRight(listing.Path, `\`) + `\` + name, IsDirectory: kind == "<DIR>" || kind == "<JUNCTION>" || kind == "<SYMLINKD>", IsLink: isLink, ModifiedText: strings.Join(strings.Fields(match[1]+" "+match[2]), " ")}
		if !strings.HasPrefix(kind, "<") {
			n, err := strconv.ParseUint(strings.NewReplacer(",", "", ".", "").Replace(kind), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid file size in dir response")
			}
			entry.Size = &n
		}
		for _, layout := range []string{"01/02/2006 03:04 PM", "01/02/2006 03:04:05 PM", "01/02/2006 03:04:05PM", "01/02/2006 03:04PM", "01/02/2006 15:04", "2006-01-02 15:04", "2006-01-02 15:04:05", "02.01.2006 15:04", "02.01.2006 15:04:05"} {
			if value, err := time.Parse(layout, strings.ToUpper(entry.ModifiedText)); err == nil {
				entry.ModifiedSort = value.Format("2006-01-02T15:04:05.000Z")
				break
			}
		}
		listing.Entries = append(listing.Entries, entry)
	}
	if listing.Path == "" || (recognized == 0 && !dirEmpty.MatchString(output)) {
		return nil, fmt.Errorf("could not parse Windows dir output (English directory header required); the path may be unavailable or the output format unsupported")
	}
	return listing, nil
}
