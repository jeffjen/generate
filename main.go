package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"encoding/json"
	"errors"
	"fmt"
	"os"
	"text/template"
)

func main() {
	app := cli.NewApp()
	app.Name = "generate"
	app.Usage = "gobble up context and template, out comes a generated file"
	app.Authors = []cli.Author{
		cli.Author{"Yi-Hung Jen", "yihungjen@gmail.com"},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input",
			Usage: "input file (json encoded) to read into context",
		},
		cli.StringFlag{
			Name:  "json",
			Usage: "json string to read into context",
		},
		cli.StringFlag{
			Name:  "output",
			Value: "-",
			Usage: "output location for the generated content, defaults to stdout",
		},
	}
	app.Action = ParseEndGenerate
	app.Run(os.Args)
}

func check(c *cli.Context) (ctx interface{}, tmpl *template.Template, err error) {
	if tmplPath := c.Args(); len(tmplPath) == 0 {
		err = errors.New("you need to specify a path for template")
		return
	} else if len(tmplPath) > 1 {
		err = errors.New("you should specify one template")
		return
	} else if _, err = os.Stat(tmplPath[0]); err != nil {
		err = fmt.Errorf("%s: unable to access specified tempalte", tmplPath[0])
		return
	} else if tmpl, err = template.ParseFiles(tmplPath[0]); err != nil {
		return
	}
	inn, jstr := c.String("input"), c.String("json")
	if inn != "" && jstr != "" {
		err = errors.New("input and json flag are exclusive")
		return
	} else if inn != "" {
		if f, frr := os.Open(inn); frr != nil {
			err = fmt.Errorf("%s: unable to process context", inn)
			return
		} else {
			defer f.Close()
			err = json.NewDecoder(f).Decode(&ctx)
		}
	} else if jstr != "" {
		err = json.Unmarshal([]byte(jstr), &ctx)
	} else {
		err = errors.New("no context; cowardly refuse")
	}
	return
}

func ParseEndGenerate(c *cli.Context) {
	ctx, tmpl, err := check(c)
	if err != nil {
		log.Fatal(err)
	}

	var out *os.File
	if output := c.String("output"); output == "-" || output == "" {
		out = os.Stdout
	} else if out, err = os.Create(output); err != nil {
		log.Fatal(err)
	} else {
		defer out.Close()
	}

	if err = tmpl.Execute(out, ctx); err != nil {
		log.Fatal(err)
	}
}
