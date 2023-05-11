package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"mangosteen/internal/database"
	"mangosteen/internal/email"
	"mangosteen/internal/jwt_helper"
	"mangosteen/internal/router"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run() {
	rootCmd := &cobra.Command{
		Use: "mangosteen",
	}
	srvCmd := &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			RunServer()
		},
	}
	dbCmd := &cobra.Command{
		Use: "db",
	}

	emailCmd := &cobra.Command{
		Use: "email",
		Run: func(cmd *cobra.Command, args []string) {
			email.Send()
		},
	}

	mgrtCmd := &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			database.Migrate()
		},
	}
	mgrtDownCmd := &cobra.Command{
		Use: "migrate:down",
		Run: func(cmd *cobra.Command, args []string) {
			database.MigrateDown()
		},
	}
	crudCmd := &cobra.Command{
		Use: "crud",
		Run: func(cmd *cobra.Command, args []string) {
			database.Crud()
		},
	}

	createMgrtCmd := &cobra.Command{
		Use: "create:migration",
		Run: func(cmd *cobra.Command, args []string) {
			database.CreateMigration(args[0])
		},
	}
	generateHMACKeyCmd := &cobra.Command{
		Use: "generateHMACKey",
		Run: func(cmd *cobra.Command, args []string) {
			bytes, _ := jwt_helper.GenerateHMACKey()
			keyPath := viper.GetString("jwt.hmac.key_path")
			dir := filepath.Dir(keyPath)
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				log.Fatalln(err)
			}
			if err := ioutil.WriteFile(keyPath, bytes, 0644); err != nil {
				log.Fatalln(err)
			}
			fmt.Println("HMAC key saved to " + keyPath)
		},
	}

	database.Connect()
	defer database.Close()

	rootCmd.AddCommand(srvCmd, dbCmd, emailCmd, generateHMACKeyCmd)
	dbCmd.AddCommand(mgrtCmd, crudCmd, mgrtDownCmd, createMgrtCmd)
	rootCmd.Execute()
}

func RunServer() {
	r := router.New()
	err := r.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("r.Run 的下一行")

}
