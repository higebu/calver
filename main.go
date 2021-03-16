package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"
)

var (
	format   = flag.String("format", "YYYY.0M.0D", "The format string for version. See https://calver.org for the details.")
	major    = flag.String("major", "", "specific major version.")
	minor    = flag.String("minor", "", "specific minor version.")
	micro    = flag.String("micro", "", "specific micro version.")
	modifier = flag.String("modifier", "", "Modifier")

	formatMap = map[string]string{
		"YYYY":     "2006",
		"YY":       "6",
		"0Y":       "06",
		"MM":       "1",
		"0M":       "01",
		"WW":       "{{ .ShortWeek }}",
		"0W":       "{{ .ZeroPaddedWeek }}",
		"DD":       "2",
		"0D":       "02",
		"MAJOR":    "{{ .Major }}",
		"MINOR":    "{{ .Minor }}",
		"MICRO":    "{{ .Micro }}",
		"MODIFIER": "{{ .Modifier }}",
	}
)

func generateTimeFormat(format string) (string, error) {
	s := strings.Split(format, ".")
	if len(s) < 3 {
		return "", fmt.Errorf("invalid format: %s", format)
	}
	tf := make([]string, len(s))
	for i, ss := range s {
		tf[i] = formatMap[ss]
	}
	return strings.Join(tf, "."), nil
}

type Params struct {
	ShortWeek      string
	ZeroPaddedWeek string
	Major          string
	Minor          string
	Micro          string
	Modifier       string
}

func main() {
	flag.Parse()

	tf, err := generateTimeFormat(*format)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	t := time.Now().UTC()
	_, w := t.ISOWeek()
	v := t.Format(tf)
	params := Params{
		ShortWeek:      fmt.Sprintf("%d", w),
		ZeroPaddedWeek: fmt.Sprintf("%02d", w),
		Major:          *major,
		Minor:          *minor,
		Micro:          *micro,
		Modifier:       *modifier,
	}
	tmpl, err := template.New("calver").Parse(v)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = tmpl.Execute(os.Stdout, params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
