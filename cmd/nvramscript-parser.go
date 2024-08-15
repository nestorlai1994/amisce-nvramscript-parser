package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type Item struct {
	Question        []byte
	Value           [][]byte
	Content         [][]byte
	ValueStartIndex int
	Index           int
}

func (item Item) String() {
	fmt.Printf("Count: %5d, Question: %s, Value: %s \n", item.Index+1, item.Question, item.Value)
}

func (item Item) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Question string   `json:"question"`
		Value    []string `json:"value"`
		Content  []string `json:"content"`
	}{
		Question: string(item.Question),
		Value:    convertByteArrayToString(item.Value),
		Content:  convertByteArrayToString(item.Content),
	})
}

func convertByteArrayToString(b [][]byte) []string {
	strs := make([]string, len(b))
	for i, v := range b {
		strs[i] = string(v)
	}
	return strs
}

func help() {
	fmt.Println("Usage: nvramscript-parser <file_path>")
	fmt.Println("Example: nvramscript-parser ./example/nvram_script_clean_up_S5B_3A11.Q10_bm2_uefi")
	flag.PrintDefaults()
}

func main() {
	fmt.Println("Welcome to the NVRAM Script Parser!")
	fmt.Println("This program parses a file containing NVRAM script data and extracts the setup questions and their corresponding values.")

	exportToPtr := flag.String("export-to", "json", "creaet a file with supported format")
	exportPath := flag.String("export-path", "", "path to export file")
	flag.Usage = help
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	now := time.Now()
	defaultExportPath := fmt.Sprintf("/tmp/%s-parsed-%d-%d-%d-%d-%d-%d", filepath.Base(flag.Arg(0)), now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	if *exportPath == "" {
		exportPath = &defaultExportPath
	}

	fmt.Println("export-to:", *exportToPtr)
	fmt.Println("export-path:", *exportPath)

	filePath := os.Args[1]
	dat, _ := os.ReadFile(filePath)

	confs := make([]Item, 0)

	lines := make([][]byte, 0)
	line := make([]byte, 0)
	reComment := regexp.MustCompile(`(\s*)(//.*)?$`)
	reQuestion := regexp.MustCompile(`^Setup\s+Question\s*=\s*(.*)`)
	reOptions := regexp.MustCompile(`^Options\s*=\s*(.*)|^Value\s*=\s*(.*)`)
	reValue := regexp.MustCompile(`^Options\s*=\s*\*+(.*)|^Value\s*=\s*(.*)|\s*\*+(.*)`)

	for _, b := range dat {
		if b == '\n' {
			line = ridOfComment(line, reComment)
			if len(line) == 0 || (!reQuestion.Match(line) && len(confs) == 0) {
				// skip empty line and only comment line
				// will skip any contents before not appear Setup Questions
				line = make([]byte, 0)
				continue
			}

			if reQuestion.Match(line) {
				q := reQuestion.FindAllSubmatch(line, -1)
				confs = append(confs, Item{
					Question: bytes.TrimSpace(q[0][1]),
					Index:    len(confs),
					Value:    make([][]byte, 0),
				})
			}

			item := &confs[len(confs)-1]

			if reOptions.Match(line) {
				item.ValueStartIndex = len(lines)
			}

			if item.ValueStartIndex > 0 && reValue.Match(line) {
				val := reValue.FindAllSubmatch(line, -1)

				valWithOption := bytes.TrimSpace(val[0][1])
				valWithValue := bytes.TrimSpace(val[0][2])
				valOnly := bytes.TrimSpace(val[0][3])

				if len(valWithOption) > 0 {
					item.Value = append(item.Value, valWithOption)
				} else if len(valWithValue) > 0 {
					item.Value = append(item.Value, valWithValue)
				} else if len(valOnly) > 0 {
					item.Value = append(item.Value, valOnly)
				}
			}

			if reQuestion.Match(line) && len(lines) > 0 {
				item = &confs[len(confs)-2]
				item.Content = lines
				lines = make([][]byte, 0)
				//item.String()
			}

			lines = append(lines, line)
			line = make([]byte, 0)
			continue
		}
		line = append(line, b)
	}

	exportToJson(confs, *exportPath)

}

func ridOfComment(line []byte, re *regexp.Regexp) []byte {
	if re.Match(line) {
		line = re.ReplaceAll(line, []byte("$1"))
	}
	return bytes.TrimSpace(line)
}

func exportToJson(confs []Item, path string) {
	// jsonMarshals := make([][]byte, len(confs))
	// for i, conf := range confs {
	// 	dat, _ := conf.MarshalJSON()
	// 	jsonMarshals[i] = dat
	// }
	dat, _ := json.MarshalIndent(confs, "", "  ")
	os.WriteFile(path, dat, 0644)
}
