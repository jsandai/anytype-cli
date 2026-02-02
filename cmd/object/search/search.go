package search

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewSearchCmd() *cobra.Command {
	var (
		spaceID string
		query   string
		types   []string
		limit   int32
		jsonOut bool
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for objects",
		Long: `Search for objects in a space by query or type.

Examples:
  # Full-text search
  anytype object search --space <space-id> --query "meeting notes"

  # Filter by type
  anytype object search --space <space-id> --type ot-note

  # Combine query and type filter
  anytype object search --space <space-id> --query "project" --type ot-task`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if spaceID == "" {
				return fmt.Errorf("--space is required")
			}

			req := core.SearchObjectsRequest{
				SpaceID: spaceID,
				Query:   query,
				Types:   types,
				Limit:   limit,
			}

			results, err := core.SearchObjects(req)
			if err != nil {
				return output.Error("Search failed: %w", err)
			}

			if len(results) == 0 {
				output.Info("No objects found")
				return nil
			}

			if jsonOut {
				out, _ := json.MarshalIndent(results, "", "  ")
				output.Print(string(out))
			} else {
				output.Info("Found %d object(s):\n", len(results))
				for _, obj := range results {
					name := obj.Name
					if name == "" {
						name = "(untitled)"
					}
					output.Print("  %s", name)
					output.Print("    ID:   %s", obj.ID)
					output.Print("    Type: %s", obj.Type)
					if obj.Snippet != "" {
						// Truncate snippet
						snippet := obj.Snippet
						if len(snippet) > 100 {
							snippet = snippet[:100] + "..."
						}
						snippet = strings.ReplaceAll(snippet, "\n", " ")
						output.Print("    Preview: %s", snippet)
					}
					output.Print("")
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&spaceID, "space", "", "Space ID (required)")
	cmd.Flags().StringVar(&query, "query", "", "Full-text search query")
	cmd.Flags().StringArrayVar(&types, "type", nil, "Filter by object type (can be repeated)")
	cmd.Flags().Int32VarP(&limit, "limit", "n", 20, "Maximum number of results")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")

	return cmd
}
