package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// viperOverrides allows giving precence to config file and env vars compared to flag default value.
// Viper and Cobra integration is not really smooth, so this helper function is needed to override
// flag default values with viper provided ones.
// Should be used as parameter for cmd.Flags().VisitAll() to set flag value based on viper value.
//
// NOTE: viper value is not updated with flag value, so flag value should be read.
// See https://github.com/spf13/viper/discussions/1061
// See https://github.com/spf13/viper/issues/671
// See https://github.com/spf13/viper/issues/375
func viperOverrides(c *cobra.Command) func(*pflag.Flag) {
	return func(f *pflag.Flag) {
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			err := c.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
