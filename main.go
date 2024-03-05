package main

import (
	"fmt"
	"github.com/spf13/cobra"
	cmd2 "gugit/cmd"
	"gugit/internal"
	"log"
)

var (
	rootCMD = &cobra.Command{
		Use:   "Use like a git",
		Short: "VCS implemented in GO",
		Long:  "Longer description i guess",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("See usage")
		},
	}

	initCMD = &cobra.Command{
		Use:   "init",
		Short: "Creates ugit directory.",
		Long:  "Creates ugit and ugit/objects directories. Must be used before trying other cmds.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd2.Init()
		},
	}

	fileCMD = &cobra.Command{Use: "file",
		Short: "Adds file to objects dir",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			args[0] = cmd2.ResolveName(args[0])
			cmd2.FileCMD(args[0], internal.BLOB)
		},
	}

	catCMD = &cobra.Command{Use: "cat",
		Short: "Prints contents of file to stdout.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			args[0] = cmd2.ResolveName(args[0])
			cmd2.CatCMD(args[0])
		}}

	wTreeCMD = &cobra.Command{Use: "writeTree",
		Short: "Write tree (directory) to objects subdir.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			args[0] = cmd2.ResolveName(args[0])
			cmd2.WriteTree(args[0])
		}}

	rTreeCMD = &cobra.Command{Use: "readTree",
		Short: "Write tree (directory) to objects subdir.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			args[0] = cmd2.ResolveName(args[0])
			cmd2.ReadTree(args[0])
		}}

	commitCMD = &cobra.Command{Use: "commit",
		Short: "Create new commit file containing saved tree, time, author and parent commit." +
			"Commit moves head to newest created commit.",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cmd2.Commit()
		}}

	logCMD = &cobra.Command{Use: "log",
		Short: "Log history of commits. If no arg is given it will start from HEAD.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				args[0] = cmd2.ResolveName(args[0])
				cmd2.Log(args[0])
			} else {
				cmd2.Log("")
			}
		}}

	checkoutCMD = &cobra.Command{Use: "checkout",
		Short: "Checks commit with oid given as only argument. Allows to move in past.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmd2.Checkout(args[0])
		}}

	tagCMD = &cobra.Command{Use: "tag",
		Short: "Creates new tag (alias) for a commit. Must be two arguments name and oid of commit." +
			"Will create directory .ugit/objects/tag if needed.",
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				cmd2.Tag(args[0], "")
			} else {
				args[1] = cmd2.ResolveName(args[1])
				cmd2.Tag(args[0], args[1])
			}
		}}

	kCMD = &cobra.Command{Use: "k",
		Short: "Visualize commit tree.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd2.K()
		}}

	branchCMD = &cobra.Command{Use: "branch",
		Short: "Create new branch. First argument is name and second OID.",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				log.Fatal("Incorrect number of arguments.")
			}
			if len(args) == 0 {
				cmd2.ListBranches()
			} else {
				args[1] = cmd2.ResolveName(args[1])
				cmd2.Branch(args[0], args[1])
			}
		}}

	statusCMD = &cobra.Command{Use: "status",
		Short: "Show your current branch.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd2.Status()
		}}

	resetCMD = &cobra.Command{Use: "reset",
		Short: "Moves HEAD and current branch to commit with oid given as argument.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmd2.Reset(args[0])
		}}

	diffCMD = &cobra.Command{Use: "diff",
		Short: "Show line by line difference between new commit and its parent.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmd2.Show(args[0])
		}}

	mergeCMD = &cobra.Command{Use: "merge",
		Short: "Merge other branch given as argument to branch pointed by HEAD.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmd2.Merge("refs/heads/" + args[0])
		}}
)

func main() {
	rootCMD.AddCommand(initCMD, fileCMD, catCMD, wTreeCMD, rTreeCMD, tagCMD,
		logCMD, commitCMD, checkoutCMD, kCMD, branchCMD, statusCMD, resetCMD, diffCMD, mergeCMD)
	err := rootCMD.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
