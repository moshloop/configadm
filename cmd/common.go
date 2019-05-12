package cmd

import (
	_ "github.com/moshloop/configadm/pkg"
	"github.com/moshloop/configadm/pkg/os"
	"github.com/moshloop/configadm/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func GetConfig(cmd *cobra.Command) *types.Config {

	configs, err := cmd.Flags().GetStringSlice("config")
	if err != nil {
		log.Fatalf("%s", err)
	}
	vars, err := cmd.Flags().GetStringSlice("var")
	if err != nil {
		log.Fatalf("%s", err)
	}

	flags := []types.Flag{}

	if ok, _ := cmd.Flags().GetBool("detect-tags"); ok {
		for _, _os := range os.SupportedOperatingSystems {
			if _os.DetectAtRuntime() {
				log.Infof("Detected %s\n", _os.GetTag())
				flags = append(flags, types.FLAG_MAP[_os.GetTag()])
			}
		}
	}

	flagNames, err := cmd.Flags().GetStringSlice("tag")
	for _, name := range flagNames {

		if flag, ok := types.FLAG_MAP[name]; ok {
			flags = append(flags, flag)
		} else {
			log.Fatalf("Unknown flag %s", name)
		}

	}
	if err != nil {
		log.Fatalf("%s", err)
	}

	cfg, err := types.NewConfig(configs...).
		WithVars(vars...).
		WithFlags(flags...).
		Build()

	if err != nil {
		panic(nil)
	}
	return cfg

}
