package cobra

import (
	"fmt"
	"github.com/pablogolobaro/secret"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		vault := secret.File(encodingKey, secretsPath())

		key, value := args[0], args[1]

		err := vault.Set(key, value)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Value set succesfully\n")
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
