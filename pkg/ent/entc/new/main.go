package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"text/template"
	"unicode"
)

// go run ./pkg/ent/cmd/new --target ./pkg/ent/schema --template ./pkg/ent/template/entnew.tmpl tableNames schemaName
func main() {
	var tableNames, schemaName, target, templateFile string
	flag.StringVar(&tableNames, "tables", "", "Table names to create new schema")
	flag.StringVar(&schemaName, "schema", "", "Schema name to create new schema")
	flag.StringVar(&target, "target", "", "Target directory to create new schema")
	flag.StringVar(&templateFile, "template", "", "Template file to use")
	flag.Parse()

	if schemaName == "" {
		log.Fatalf("schemaName is required")
	}
	if tableNames == "" || target == "" || templateFile == "" {
		log.Fatalf("tableNames, target, and template are required")
	}

	tmpl, err := template.New(path.Base(templateFile)).Funcs(template.FuncMap{
		"snake": snake,
	}).ParseFiles(templateFile)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	tables := strings.Split(tableNames, ",")
	for _, table := range tables {
		generateFile(tmpl, table, schemaName)
	}
}

func generateFile(tmpl *template.Template, table, schema string) {
	log.Infof("Generating schema for table: %s, schema: %s", table, schema)

	// 출력 파일 생성
	filePath := "./pkg/ent/schema/" + strings.ToLower(table) + ".go"
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// 템플릿 실행 및 파일에 쓰기
	err = tmpl.Execute(file, struct {
		Table, Schema string
	}{
		Table:  table,
		Schema: schema,
	})
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

// cfg) https://github.com/ent/ent/blob/master/entc/gen/func.go
// snake converts the given struct or field name into a snake_case.
//
//	Username => username
//	FullName => full_name
//	HTTPCode => http_code
func snake(s string) string {
	var (
		j int
		b strings.Builder
	)
	for i := 0; i < len(s); i++ {
		r := rune(s[i])
		// Put '_' if it is not a start or end of a word, current letter is uppercase,
		// and previous is lowercase (cases like: "UserInfo"), or next letter is also
		// a lowercase and previous letter is not "_".
		if i > 0 && i < len(s)-1 && unicode.IsUpper(r) {
			if unicode.IsLower(rune(s[i-1])) ||
				j != i-1 && unicode.IsLower(rune(s[i+1])) && unicode.IsLetter(rune(s[i-1])) {
				j = i
				b.WriteString("_")
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}
