package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/but80/go-smaf/v2/chunk"
	"github.com/but80/go-smaf/v2/log"
	"github.com/but80/go-smaf/v2/voice"
	"github.com/golang/protobuf/proto"
	"github.com/urfave/cli"
)

var version string

func init() {
	if version == "" {
		version = "unknown"
	}
}

var dumpCmd = cli.Command{
	Name:      "dump",
	Aliases:   []string{"d"},
	Usage:     "Dump SMAF format files (.mmf|.spf|.vma|.vm3|.vm5)",
	ArgsUsage: "<filename>",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "json, j",
			Usage: `Dump in JSON format`,
		},
		cli.BoolFlag{
			Name:  "protobuf, p",
			Usage: `Dump in protobuf format`,
		},
		cli.BoolFlag{
			Name:  "voice, v",
			Usage: `Dump voice data only`,
		},
		cli.BoolFlag{
			Name:  "exclusive, x",
			Usage: `Dump exclusives only`,
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: `Show debug messages`,
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: `Suppress information messages`,
		},
		cli.BoolFlag{
			Name:  "silent, Q",
			Usage: `Do not output any messages`,
		},
	},
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() < 1 {
			cli.ShowCommandHelp(ctx, "dump")
			os.Exit(1)
		}
		if ctx.Bool("debug") {
			log.SetLevel(log.LevelDebug)
		} else if ctx.Bool("silent") {
			log.SetLevel(log.LevelNone)
		} else if ctx.Bool("quiet") {
			log.SetLevel(log.LevelWarn)
		}
		file := ctx.Args()[0]
		ext := ""
		i := len(file) - 4
		if 0 <= i {
			ext = strings.ToLower(file[i:])
		}
		var data fmt.Stringer
		var err error
		switch ext {
		case ".mmf", ".spf":
			fc, err := chunk.NewFileChunk(file)
			data = fc
			if err == nil && (ctx.Bool("voice") || ctx.Bool("exclusive")) {
				exclusives := fc.CollectExclusives()
				data = exclusives
				if ctx.Bool("voice") {
					data = exclusives.Voices()
				}
			}
		case ".vma":
			data, err = voice.NewVMAVoiceLib(file)
		case ".vm3":
			data, err = voice.NewVM3VoiceLib(file)
		case ".vm5":
			data, err = voice.NewVM5VoiceLib(file)
		default:
			return cli.NewExitError(fmt.Errorf("Unknown file extension"), 1)
		}
		if err != nil {
			switch data.(type) {
			case nil:
				return cli.NewExitError(err, 1)
			default:
				log.Warnf(err.Error())
			}
		}
		if ctx.Bool("json") {
			j, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			fmt.Println(string(j))
		} else if ctx.Bool("protobuf") {
			switch d := data.(type) {
			case *voice.VM5VoiceLib:
				b, err := proto.Marshal(d.ToPB())
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				fmt.Print(string(b))
			default:
				return cli.NewExitError(fmt.Errorf("Protobuf conversion for %s is not supported", ext), 1)
			}
		} else {
			fmt.Println(data.String())
		}
		return nil
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "go-smaf"
	app.Version = version
	app.Usage = "Convert SMAF format files"
	app.Authors = []cli.Author{
		{
			Name:  "but80",
			Email: "mersenne.sister@gmail.com",
		},
	}
	app.HelpName = "go-smaf"

	app.Commands = []cli.Command{
		dumpCmd,
	}

	app.Action = func(ctx *cli.Context) error {
		cli.ShowAppHelp(ctx)
		return nil
	}

	app.Run(os.Args)
}