/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/meanii/teagit/database"
	"github.com/meanii/teagit/models"
	util "github.com/meanii/teagit/utils"
	"github.com/spf13/cobra"
)

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "manager git users",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if cmd.Flag("add").Changed {
			user := models.Users{}

			user.Name = util.Ask("Enter your name: ")
			user.Email = util.Ask("Enter your email: ")

			// save user
			database.NewDatabase().Db.Create(&user)

			fmt.Println("User added successfully")
			fmt.Printf("Name: %s\nEmail: %s\n", user.Name, user.Email)
			return
		}

		// list users
		var users []models.Users
		database.NewDatabase().Db.Find(&users)

		if len(users) == 0 {
			fmt.Println("No users found")
			return
		}

		fmt.Println("Users:")
		for _, user := range users {
			fmt.Printf("Name: %s\nEmail: %s\n\n", user.Name, user.Email)
		}
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)

	// add new user
	usersCmd.Flags().BoolP("add", "a", false, "add new user")
}
