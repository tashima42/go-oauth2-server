package main

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tashima42/go-oauth2-server/db"
	"github.com/tashima42/go-oauth2-server/helpers"
)

func main() {
	log.Println("Initializing application")
	rootCmd := &cobra.Command{
		Use:   "go-oauth2-server",
		Short: "OAuth2 server",
		Long:  `OAuth2 server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var conf db.Config
			conf.FromEnv()
			repo, err := db.Open(conf)
			if err != nil {
				return errors.Wrap(err, "failed to open database")
			}
			hashHelper := helpers.GetHashHelperInstance()
			Serve(repo, hashHelper)
			// TODO: get signal and close database
			return nil
		},
	}
	rootCmd.Execute()
}
