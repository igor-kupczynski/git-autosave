package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/igor-kupczynski/git-autosave/repo"
	"github.com/spf13/cobra"
)

var cmdSave = &cobra.Command{
	Use:   "save",
	Short: "Commit any staged and unstaged changes to autosave/<branch>.",
	Run: func(cmd *cobra.Command, args []string) {

		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error autocommiting: %v", err)
		}

		_ = dir

		if err := CommitAllChanges(nil); err != nil {
			log.Fatalf("Error autocommiting: %v", err)
		}
	},
}

func CommitAllChanges(repo repo.Repository) error {
	headName, err := repo.GetCurrentBranch()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(headName, "autosave/") {
		if err := repo.CheckoutSpinOffBranch(fmt.Sprintf("autosave/%s", headName)); err != nil {
			return err
		}
	}

	if err := repo.CommitAllChanged(fmt.Sprintf("git-autosave at %s", time.Now())); err != nil {
		return err
	}

	return nil
}
