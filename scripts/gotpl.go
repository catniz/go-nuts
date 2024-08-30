package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Usage: go run gotpl.go -t template.yaml.tpl -v values.yaml > template.yaml

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return fmt.Sprintf("%v", *s)
}
func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var (
	templateFile = flag.String("t", "", "template.yaml.tpl file")
	valueFiles   stringSliceFlag
)

func main() {
	// 커맨드라인 플래그 파싱
	flag.Var(&valueFiles, "v", "values.yaml files")
	flag.Parse()

	// 플래그 검증
	if _, err := os.Stat(*templateFile); os.IsNotExist(err) {
		fatalf("-t: %s 파일이 존재하지 않습니다.", *templateFile)
	}
	for _, valueFile := range valueFiles {
		if _, err := os.Stat(valueFile); os.IsNotExist(err) {
			fatalf("-v: %s 파일이 존재하지 않습니다.", valueFile)
		}
	}

	// values.yaml 파일들을 순서대로 읽어서 합칩니다.
	values := make(map[string]interface{})
	for _, valueFile := range valueFiles {
		// values.yaml 파일 읽기
		yamlFile, err := os.ReadFile(valueFile)
		if err != nil {
			fatalf("Error reading %s: %v", valueFile, err)
		}

		// values 맵에 언마샬링
		err = yaml.Unmarshal(yamlFile, &values)
		if err != nil {
			fatalf("Error unmarshaling %s: %v", valueFile, err)
		}
	}

	// 템플릿 함수 등록
	tmpl := template.New(filepath.Base(*templateFile)).Funcs(funcMap())

	// 템플릿 파일 파싱
	tmpl, err := tmpl.ParseFiles(*templateFile)
	if err != nil {
		fatalf("Error parsing file %s: %v", *templateFile, err)
	}

	// 템플릿 실행해서 stdout 으로 출력
	err = tmpl.Execute(os.Stdout, values)
	if err != nil {
		fatalf("Error executing template: %v", err)
	}
}

func fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}

func funcMap() template.FuncMap {
	// helm 이 사용하는 sprig go template function library 똑같이 사용합시다
	f := sprig.TxtFuncMap()

	// helm 에서 따로 추가한 함수들 가능한 추가
	f["toYAML"] = toYAML
	f["fromYAML"] = fromYAML
	f["fromYAMLArray"] = fromYAMLArray
	f["toJSON"] = toJSON
	f["fromJSON"] = fromJSON
	f["fromJSONArray"] = fromJSONArray

	// readFile 함수는 템플릿 파일의 경로를 기준으로 상대경로로 파일을 읽습니다.
	f["readFile"] = readFile
	return f
}

func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}

func fromYAML(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

func fromYAMLArray(str string) []interface{} {
	var a []interface{}
	if err := yaml.Unmarshal([]byte(str), &a); err != nil {
		a = []interface{}{err.Error()}
	}
	return a
}

func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

func fromJSON(str string) map[string]interface{} {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

func fromJSONArray(str string) []interface{} {
	var a []interface{}
	if err := json.Unmarshal([]byte(str), &a); err != nil {
		a = []interface{}{err.Error()}
	}
	return a
}

func readFile(path string) string {
	path = filepath.Join(filepath.Dir(*templateFile), path)
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(f)
}
