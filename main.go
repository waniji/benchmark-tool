package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var helpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [options]

VERSION:
   {{.Version}}{{if or .Author .Email}}

AUTHOR:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`

func main() {
	app := cli.NewApp()
	app.Name = "benchmark-tool"
	app.Usage = "Web server benchmarking tool"
	app.HideHelp = true
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url, u",
			Usage: "アクセスするURL",
		},
		cli.IntFlag{
			Name:  "count, c",
			Value: 1,
			Usage: "URLにアクセスする回数",
		},
		cli.IntFlag{
			Name:  "worker, w",
			Value: 1,
			Usage: "同時アクセス数",
		},
		cli.StringFlag{
			Name:  "basic-auth-user",
			Value: "",
			Usage: "BASIC認証に使用するユーザー",
		},
		cli.StringFlag{
			Name:  "basic-auth-pass",
			Value: "",
			Usage: "BASIC認証に使用するパスワード",
		},
	}
	app.Action = bench
	cli.AppHelpTemplate = helpTemplate

	app.Run(os.Args)
}
