package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/ghcri/ghcri/kakashi"
	"github.com/ghcri/ghcri/stackbrew"
)

var Version string

var app = cli.App{
	Name:    "ghcri",
	Version: Version,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "registry",
			Usage: "The registry url, like docker.io or ghcr.io",
			EnvVars: []string{
				"GHCRI_REGISTRY",
			},
			Required: true,
			Value:    "ghcr.io",
		},
		&cli.StringFlag{
			Name:  "owner",
			Usage: "The registry owner",
			EnvVars: []string{
				"GHCRI_OWNER",
			},
			Required: true,
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "The registry username",
			EnvVars: []string{
				"GHCRI_USERNAME",
			},
			Required: true,
		},
		&cli.StringFlag{
			Name:  "password",
			Usage: "The registry password",
			EnvVars: []string{
				"GHCRI_PASSWORD",
			},
			Required: true,
		},
	},
	Before: func(c *cli.Context) error {
		if c.Args().Len() < 1 {
			return fmt.Errorf("please input source files from docker-libiray")
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		logger, _ := zap.NewDevelopment()
		pool, _ := ants.NewPool(8) // We will use 8 workers.
		wg := &sync.WaitGroup{}

		k, err := kakashi.New(
			c.String("registry"),
			c.String("username"),
			c.String("password"))
		if err != nil {
			return err
		}

		dentry, err := os.ReadDir(c.Args().First())
		if err != nil {
			logger.Error("read dir",
				zap.String("path", c.Args().First()),
				zap.Error(err))
			return err
		}

		for _, file := range dentry {
			path := filepath.Join(c.Args().First(), file.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				logger.Error("open file",
					zap.String("path", c.Args().First()),
					zap.Error(err))
				return err
			}

			sb := stackbrew.ParseBytes(content)

			logger.Info("Start handling", zap.String("file", file.Name()))

			imageName := file.Name()
			for _, stack := range sb.Stacks {
				var tags []string
				tags = append(tags, stack.Tags...)
				tags = append(tags, stack.SharedTags...)

				for _, tag := range tags {
					oldImage := fmt.Sprintf("docker://%s:%s", imageName, tag)
					newImage := fmt.Sprintf("docker://%s/%s/%s:%s", c.String("registry"), c.String("owner"), imageName, tag)

					wg.Add(1)
					err = pool.Submit(func() {
						defer wg.Done()
						_ = k.Copy(oldImage, newImage)
					})
					if err != nil {
						logger.Error("submit task", zap.Error(err))
						return err
					}
				}
			}
		}

		wg.Wait()
		return nil
	},
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
