package cmd

import (
	"fmt"
	"gugit/internal"
	"strings"
)

type commitNode struct {
	commit internal.Commit
	refs   []string
	column int
}

func K() {
	commits := make(map[string]*commitNode)
	refs := make(map[string]string)

	for _, ref := range genRefs() {
		resolved := ResolveName(ref)
		if resolved == "" {
			continue // Skip if resolution fails
		}
		refs[ref] = resolved

		// Get actual commit
		commit := getCommit(resolved)
		if _, ok := commits[resolved]; !ok {
			commits[resolved] = &commitNode{
				commit: commit,
				refs:   []string{ref},
			}
		} else {
			commits[resolved].refs = append(commits[resolved].refs, ref)
		}
	}

	// If no commits found, print message and return
	if len(commits) == 0 {
		fmt.Println("No commits yet")
		return
	}

	// Get HEAD for special marking
	headOid := ""
	if _, head, err := getRef("HEAD", true); err == nil && head.value != "" {
		headOid = head.value
	}

	var history []*commitNode
	visited := make(map[string]bool)
	queue := []string{headOid}

	for len(queue) > 0 {
		oid := queue[0]
		queue = queue[1:]

		if visited[oid] {
			continue
		}
		visited[oid] = true

		node := commits[oid]
		if node == nil {
			commit := getCommit(oid)
			node = &commitNode{
				commit: commit,
			}
			commits[oid] = node
		}
		history = append(history, node)

		queue = append(queue, node.commit.Parent...)
	}

	maxCol := 0
	for i, node := range history {
		if node.column == 0 && i > 0 {
			node.column = maxCol
			maxCol++
		}
		for _, parent := range node.commit.Parent {
			if pNode := commits[parent]; pNode != nil && pNode.column == 0 {
				pNode.column = node.column
			}
		}
	}

	for i, node := range history {
		line := make([]string, maxCol*2+1)
		for j := range line {
			line[j] = " "
		}

		col := node.column
		line[col*2] = "*"

		for _, parent := range node.commit.Parent {
			if pNode := commits[parent]; pNode != nil {
				parentCol := pNode.column
				start := min(col, parentCol) * 2
				end := max(col, parentCol) * 2
				for j := start; j <= end; j++ {
					if line[j] == " " {
						line[j] = "-"
					}
				}
			}
		}

		// Print commit info
		graph := strings.Join(line, "")
		shortOid := node.commit.Oid[:8]
		msg := node.commit.Msg
		refs := ""
		if len(node.refs) > 0 {
			prettyRefs := make([]string, 0)
			for _, ref := range node.refs {
				ref = strings.TrimPrefix(ref, "refs/heads/")
				ref = strings.TrimPrefix(ref, "refs/tags/")
				prettyRefs = append(prettyRefs, ref)
			}
			refs = fmt.Sprintf(" (%s)", strings.Join(prettyRefs, ", "))
		}

		fmt.Printf("%s %s%s %s\n", graph, shortOid, refs, msg)

		if i < len(history)-1 {
			for j := range line {
				line[j] = " "
			}
			for _, next := range history[i+1:] {
				line[next.column*2] = "|"
			}
			fmt.Println(strings.Join(line, ""))
		}
	}
}
