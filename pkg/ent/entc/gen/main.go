package main

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"flag"
	"log"
	"strings"
)

func main() {
	var target, feature, templateDir string
	flag.StringVar(&target, "target", "", "Target path to generate")
	flag.StringVar(&feature, "feature", "", "Feature names")
	flag.StringVar(&templateDir, "templatedir", "", "Directory that template files are located")
	flag.Parse()

	opts := []entc.Option{
		entc.TemplateDir(templateDir),
		entc.FeatureNames(strings.Split(feature, ",")...),
	}
	err := entc.Generate(target, &gen.Config{}, opts...)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
