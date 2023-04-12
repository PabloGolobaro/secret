package cobra

import (
	"fmt"
	"github.com/pablogolobaro/secret"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		vault := secret.File(encodingKey, secretsPath())

		key := args[0]

		value, err := vault.Get(key)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s=%s\n", key, value)

	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
