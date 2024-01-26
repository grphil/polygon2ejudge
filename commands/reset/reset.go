package reset

import (
	"github.com/hellflame/argparse"
	"polygon2ejudge/lib/config"
)

func AddResetConfigCommand(parser *argparse.Parser) {
	reset := parser.AddCommand("reset", "Reset all configuration files", nil)
	globalConfigs := reset.Flag("g", "global", &argparse.Option{
		Help: "Reset global options files",
	})

	reset.InvokeAction = func(invoked bool) {
		if !invoked {
			return
		}

		if *globalConfigs {
			config.ResetGlobalConfig()
		} else {
			config.ResetUserConfig()
		}
	}
}
