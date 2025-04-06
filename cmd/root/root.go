package root

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "my-app",
	Short: "Мое cli приложение",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Что-то создает",
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Что-то удаляет",
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Создает пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		usernamesStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("user %s created", usernamesStr)
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Удаляет пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		usernamesStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("user %s deleted", usernamesStr)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)

	createCmd.AddCommand(createUserCmd)
	deleteCmd.AddCommand(deleteUserCmd)

	createUserCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	err := createUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatal(err)
	}

	deleteUserCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	err = deleteUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatal(err)
	}

}
