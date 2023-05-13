package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"mangosteen/internal/database"
	"mangosteen/internal/email"
	"mangosteen/internal/jwt_helper"
	"mangosteen/internal/router"
	"net/http"
	"os"
	"os/exec"
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

	coverCmd := &cobra.Command{
		Use: "cover",
		Run: func(cmd *cobra.Command, args []string) {
			os.MkdirAll("coverage", os.ModePerm)
			if err := exec.Command(
				"go", "test", "-coverprofile=coverage/cover.out", "./...",
			).Run(); err != nil {
				log.Fatalln(err)
			}
			if err := exec.Command(
				"go", "tool", "cover", "-html=coverage/cover.out", "-o", "coverage/index.html",
			).Run(); err != nil {
				log.Fatalln(err)
			}

			path := "coverage/index.html"
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				log.Fatalln(err)
			}
			defer file.Close()
			css, err := ioutil.ReadFile("config/coverage.css")
			if err != nil {
				log.Fatalln(err)
			}
			cssString := string(css)
			if _, err = file.Write([]byte("\n<style>" + cssString + "</style>")); err != nil {
				log.Fatalln(err)
			}
			var port string
			if len(args) > 0 {
				port = args[0]
			} else {
				port = "8888"
			}
			fmt.Println("http://localhost:" + port + "/coverage/index.html")
			if err := http.ListenAndServe(":"+port, http.FileServer(http.Dir("."))); err != nil {
				log.Fatalln(err)
			}
		},
	}

	database.Connect()
	defer database.Close()

	rootCmd.AddCommand(srvCmd, dbCmd, emailCmd, generateHMACKeyCmd, coverCmd)
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
