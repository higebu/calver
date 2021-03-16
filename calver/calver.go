package calver

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"
)

var (
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

type params struct {
	ShortWeek      string
	ZeroPaddedWeek string
	Major          string
	Minor          string
	Micro          string
	Modifier       string
}

func Generate(format, major, minor, micro, modifier string) (string, error) {
	tf, err := generateTimeFormat(format)
	if err != nil {
		return "", err
	}
	t := time.Now().UTC()
	_, w := t.ISOWeek()
	v := t.Format(tf)
	p := params{
		ShortWeek:      fmt.Sprintf("%d", w),
		ZeroPaddedWeek: fmt.Sprintf("%02d", w),
		Major:          major,
		Minor:          minor,
		Micro:          micro,
		Modifier:       modifier,
	}
	tmpl, err := template.New("calver").Parse(v)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, p)
	if err != nil {
		return "", err
	}
	return string(b.Bytes()), nil
}
