package main

import (
	"fmt"
	"os"

	"github.com/CornerMonkey/media-mcp/internal/config"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "media-mcp configured: sonarr=%t radarr=%t mutations=%t\n",
		cfg.Sonarr.Configured(),
		cfg.Radarr.Configured(),
		cfg.AllowMutations,
	)
	fmt.Fprintln(os.Stderr, "MCP server runtime is not implemented yet")
}
