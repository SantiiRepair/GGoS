package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		target, _ := cmd.Flags().GetString("target")
		duration, _ := cmd.Flags().GetInt("duration")
		rps, _ := cmd.Flags().GetInt("rps")
		plist, _ := cmd.Flags().GetString("plist")
		executors, _ := cmd.Flags().GetInt("executors")
		run(target, plist, duration, rps, executors)
	},
}

func init() {
	rootCmd.Flags().StringP("target", "t", "", "")
	rootCmd.Flags().StringP("plist", "p", "", "")
	rootCmd.Flags().IntP("duration", "d", 0, "")
	rootCmd.Flags().IntP("rps", "r", 0, "")
	rootCmd.Flags().IntP("executors", "e", 0, "")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
