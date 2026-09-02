package telemetry

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

const (
	treeIconOK   = "✓"
	treeIconFail = "✗"
	treeIconRun  = "·"
)

// TreeNode represents a single step in an execution hierarchy.
// It is concurrency-safe: child nodes may be added from goroutines.
type TreeNode struct {
	mu       sync.Mutex
	Name     string
	failed   bool
	err      error
	start    time.Time
	end      time.Time
	finished bool
	Children []*TreeNode
}

// NewRootNode creates the root of an execution tree with the current time as start.
func NewRootNode(name string) *TreeNode {
	return &TreeNode{Name: name, start: time.Now()}
}

// AddChild creates and appends a new child node. Safe to call concurrently.
func (n *TreeNode) AddChild(name string) *TreeNode {
	child := &TreeNode{Name: name, start: time.Now()}
	n.mu.Lock()
	n.Children = append(n.Children, child)
	n.mu.Unlock()
	return child
}

// Finish marks the node as done. Pass non-nil err to mark it as failed.
func (n *TreeNode) Finish(err error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.finished {
		return
	}
	n.finished = true
	n.end = time.Now()
	if err != nil {
		n.failed = true
		n.err = err
	}
}

// FinishIfRunning marks the node as done only if it hasn't finished yet.
func (n *TreeNode) FinishIfRunning(err error) {
	n.mu.Lock()
	if !n.finished {
		n.mu.Unlock()
		n.Finish(err)
	} else {
		n.mu.Unlock()
	}
}

// Print renders the complete tree to w once.
func (n *TreeNode) Print(w io.Writer) {
	n.mu.Lock()
	name := n.Name
	dur := n.duration()
	n.mu.Unlock()

	fmt.Fprintf(w, "\n%s\n", name)
	n.mu.Lock()
	children := make([]*TreeNode, len(n.Children))
	copy(children, n.Children)
	n.mu.Unlock()

	for i, child := range children {
		isLast := i == len(children)-1
		printTreeNode(w, child, "", isLast)
	}
	fmt.Fprintf(w, "  total: %s\n\n", formatDuration(dur))
}

func printTreeNode(w io.Writer, n *TreeNode, prefix string, isLast bool) {
	n.mu.Lock()
	name := n.Name
	failed := n.failed
	finished := n.finished
	err := n.err
	dur := n.duration()
	children := make([]*TreeNode, len(n.Children))
	copy(children, n.Children)
	n.mu.Unlock()

	connector := "├──"
	childPrefix := prefix + "│   "
	if isLast {
		connector = "└──"
		childPrefix = prefix + "    "
	}

	icon := treeIconOK
	if failed {
		icon = treeIconFail
	} else if !finished {
		icon = treeIconRun
	}

	// Right-align duration in 70-char wide terminal
	label := fmt.Sprintf("%s %s %s", connector, icon, name)
	durStr := formatDuration(dur)
	padding := 60 - len(prefix) - len(label)
	if padding < 1 {
		padding = 1
	}
	fmt.Fprintf(w, "%s%s%s%s\n", prefix, label, strings.Repeat(" ", padding), durStr)

	if failed && err != nil {
		errPrefix := childPrefix
		if len(children) == 0 {
			fmt.Fprintf(w, "%s    error: %s\n", errPrefix, err.Error())
		}
	}

	for i, child := range children {
		childIsLast := i == len(children)-1
		printTreeNode(w, child, childPrefix, childIsLast)
	}

	if failed && err != nil && len(children) > 0 {
		fmt.Fprintf(w, "%s    error: %s\n", childPrefix, err.Error())
	}
}

func (n *TreeNode) duration() time.Duration {
	if n.finished && !n.end.IsZero() {
		return n.end.Sub(n.start)
	}
	if !n.start.IsZero() {
		return time.Since(n.start)
	}
	return 0
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "—"
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}
