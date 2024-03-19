package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var ctx = context.Background()

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}
}

var rootCmd = &cobra.Command{
	Use:   "signature-apis",
	Short: "Root command of signature-apis",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
