package shell

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

func NewShellCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Start the Anytype interactive shell",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Starting Anytype interactive shell. Type 'exit' to quit.")
			return runShell(rootCmd)
		},
	}
}

func runShell(rootCmd *cobra.Command) error {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>> ",
		HistoryLimit:    1000,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		AutoComplete:    buildCompleter(rootCmd),
	})
	if err != nil {
		return fmt.Errorf("failed to initialize readline: %w", err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				fmt.Println("Use 'exit' or 'quit' to leave the shell")
				continue
			}
		} else if err == io.EOF {
			return nil
		} else if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		line = strings.TrimSpace(line)

		if line == "exit" || line == "quit" {
			return nil
		}

		if line == "" {
			continue
		}

		args := strings.Split(line, " ")
		rootCmd.SetArgs(args)

		if err := rootCmd.Execute(); err != nil {
			fmt.Println("Command error:", err)
		}
	}
}

func buildCompleter(rootCmd *cobra.Command) *readline.PrefixCompleter {
	var items []readline.PrefixCompleterInterface

	for _, cmd := range rootCmd.Commands() {
		if cmd.Hidden {
			continue
		}

		var subItems []readline.PrefixCompleterInterface
		for _, subCmd := range cmd.Commands() {
			if !subCmd.Hidden {
				subItems = append(subItems, readline.PcItem(subCmd.Name()))
			}
		}

		if len(subItems) > 0 {
			items = append(items, readline.PcItem(cmd.Name(), subItems...))
		} else {
			items = append(items, readline.PcItem(cmd.Name()))
		}
	}

	items = append(items,
		readline.PcItem("exit"),
		readline.PcItem("quit"),
		readline.PcItem("help"),
	)

	return readline.NewPrefixCompleter(items...)
}
