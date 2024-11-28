package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ssss-tantalum/todos-api/cmd/todos/migrations"
	"github.com/ssss-tantalum/todos-api/internal/config"
	"github.com/ssss-tantalum/todos-api/internal/database"
	"github.com/ssss-tantalum/todos-api/internal/routes"
	"github.com/ssss-tantalum/todos-api/internal/todos"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "todos",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "dev",
				Usage: "environment",
			},
		},
		Commands: []*cli.Command{
			apiCommand,
			newDBCommand(migrations.Migrations),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var apiCommand = &cli.Command{
	Name:  "api",
	Usage: "start API server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Value: ":8000",
			Usage: "serve address",
		},
	},
	Action: func(c *cli.Context) error {
		cfg, err := config.Load(c.Command.Name, c.String("env"))
		if err != nil {
			return err
		}

		db, err := database.Connect(cfg.DB.DSN, cfg.Debug)
		if err != nil {
			return err
		}

		ctx, app := todos.Start(c.Context, cfg, db)
		srv := &http.Server{
			Addr:    c.String("addr"),
			Handler: app.Router(),
		}

		routes.InitRoutes(app)

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("ListenAndServe failed: %s", err)
			}
		}()

		fmt.Printf("listening on %s\n", srv.Addr)
		fmt.Println(app.WaitExitSignal())

		return srv.Shutdown(ctx)
	},
}

func newDBCommand(migrations *migrate.Migrations) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					cfg, err := config.Load(c.Command.Name, c.String("env"))
					if err != nil {
						return err
					}
					db, err := database.Connect(cfg.DB.DSN, cfg.Debug)
					if err != nil {
						return err
					}
					ctx, app := todos.Start(c.Context, cfg, db)

					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Init(ctx)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					cfg, err := config.Load(c.Command.Name, c.String("env"))
					if err != nil {
						return err
					}
					db, err := database.Connect(cfg.DB.DSN, cfg.Debug)
					if err != nil {
						return err
					}
					ctx, app := todos.Start(c.Context, cfg, db)

					migrator := migrate.NewMigrator(app.DB(), migrations)

					group, err := migrator.Migrate(ctx)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no new migrations to run\n")
						return nil
					}

					fmt.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					cfg, err := config.Load(c.Command.Name, c.String("env"))
					if err != nil {
						return err
					}
					db, err := database.Connect(cfg.DB.DSN, cfg.Debug)
					if err != nil {
						return err
					}
					ctx, app := todos.Start(c.Context, cfg, db)

					migrator := migrate.NewMigrator(app.DB(), migrations)

					group, err := migrator.Rollback(ctx)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no groups to roll back\n")
						return nil
					}

					fmt.Printf("rolled back %s\n", group)
					return nil
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					cfg, err := config.Load(c.Command.Name, c.String("env"))
					if err != nil {
						return err
					}
					db, err := database.Connect(cfg.DB.DSN, cfg.Debug)
					if err != nil {
						return err
					}
					ctx, app := todos.Start(c.Context, cfg, db)

					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Lock(ctx)
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					cfg, err := config.Load(c.Command.Name, c.String("env"))
					if err != nil {
						return err
					}
					db, err := database.Connect(cfg.DB.DSN, cfg.Debug)
					if err != nil {
						return err
					}
					ctx, app := todos.Start(c.Context, cfg, db)

					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Unlock(ctx)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					cfg, err := config.Load(c.Command.Name, c.String("env"))
					if err != nil {
						return err
					}
					db, err := database.Connect(cfg.DB.DSN, cfg.Debug)
					if err != nil {
						return err
					}
					ctx, app := todos.Start(c.Context, cfg, db)

					migrator := migrate.NewMigrator(app.DB(), migrations)

					name := strings.Join(c.Args().Slice(), "_")
					mf, err := migrator.CreateGoMigration(ctx, name)
					if err != nil {
						return err
					}
					fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)

					return nil
				},
			},
			{
				Name:  "create_sql",
				Usage: "create up and down SQL migrations",
				Action: func(c *cli.Context) error {
					cfg, err := config.Load(c.Command.Name, c.String("env"))
					if err != nil {
						return err
					}
					db, err := database.Connect(cfg.DB.DSN, cfg.Debug)
					if err != nil {
						return err
					}
					ctx, app := todos.Start(c.Context, cfg, db)

					migrator := migrate.NewMigrator(app.DB(), migrations)

					name := strings.Join(c.Args().Slice(), "_")
					files, err := migrator.CreateSQLMigrations(ctx, name)
					if err != nil {
						return err
					}

					for _, mf := range files {
						fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					cfg, err := config.Load(c.Command.Name, c.String("env"))
					if err != nil {
						return err
					}
					db, err := database.Connect(cfg.DB.DSN, cfg.Debug)
					if err != nil {
						return err
					}
					ctx, app := todos.Start(c.Context, cfg, db)

					migrator := migrate.NewMigrator(app.DB(), migrations)

					ms, err := migrator.MigrationsWithStatus(ctx)
					if err != nil {
						return err
					}
					fmt.Printf("migrations: %s\n", ms)
					fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
					fmt.Printf("last migration group: %s\n", ms.LastGroup())

					return nil
				},
			},
			{
				Name:  "mark_applied",
				Usage: "mark migrations as applied without actually running them",
				Action: func(c *cli.Context) error {
					cfg, err := config.Load(c.Command.Name, c.String("env"))
					if err != nil {
						return err
					}
					db, err := database.Connect(cfg.DB.DSN, cfg.Debug)
					if err != nil {
						return err
					}
					ctx, app := todos.Start(c.Context, cfg, db)

					migrator := migrate.NewMigrator(app.DB(), migrations)

					group, err := migrator.Migrate(ctx, migrate.WithNopMigration())
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no new migrations to mark as applied\n")
						return nil
					}

					fmt.Printf("marked as applied %s\n", group)
					return nil
				},
			},
		},
	}
}
