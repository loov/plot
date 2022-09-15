package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"time"

	"github.com/loov/plot"
	"github.com/loov/plot/plotsvg"
)

type Commit struct {
	Hash string `json:"hash"`
	// Subject    string `json:"subject"`
	Email      string    `json:"email"`
	AuthorTime time.Time `json:"author-time"`
	CommitTime time.Time `json:"commit-time"`
}

func (c *Commit) Lifetime() time.Duration {
	return c.CommitTime.Sub(c.AuthorTime)
}

type Group struct {
	Time    time.Time
	Commits []Commit
}

func (group *Group) Lifetimes() []time.Duration {
	var xs []time.Duration
	for _, c := range group.Commits {
		xs = append(xs, c.Lifetime())
	}
	return xs
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	out, err := exec.Command("git",
		"log", "--date=iso-strict",
		`--pretty=format:{"commit": "%H", "email":"%aE", "author-time":"%ad", "commit-time": "%cd"}`,
	).CombinedOutput()
	if err != nil {
		fail(err)
	}

	maxLifetime := 6 * 30 * 24 * time.Hour

	var commits []Commit

	for _, row := range bytes.Split(out, []byte("\n")) {
		var commit Commit
		err := json.Unmarshal(row, &commit)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to parse", string(row))
			fail(err)
		}
		commits = append(commits, commit)
	}

	groupByMonth := func(c Commit) time.Time {
		t := c.CommitTime.UTC()
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	grouping := groupByMonth

	groups := map[time.Time]*Group{}
	for _, c := range commits {
		x := grouping(c)
		g, ok := groups[x]
		if !ok {
			g = &Group{Time: x}
			groups[x] = g
		}
		g.Commits = append(g.Commits, c)
	}

	ordered := []*Group{}
	for _, g := range groups {
		ordered = append(ordered, g)
	}
	sort.Slice(ordered, func(i, k int) bool {
		return ordered[k].Time.Before(ordered[i].Time)
	})

	defaultMargin := plot.R(5, 5, 5, 5)

	p := plot.New()
	stack := plot.NewVStack()
	stack.Margin = defaultMargin
	p.Add(stack)

	p.X.Min = 0
	p.Y.Max = 1
	p.Y.Min = -0.3

	p.X.Transform = plot.NewLog1pTransform(2)
	p.X.Ticks = plot.ManualTicks{
		{Value: 0, Label: "0"},
		{Value: 1, Label: "1h"},
		{Value: 2, Label: "2h"},
		{Value: 6, Label: "6h"},
		{Value: 12, Label: "12h"},
		{Value: 1 * 24, Label: "1d"},
		{Value: 2 * 24, Label: "2d"},
		{Value: 7 * 24, Label: "1w"},
		{Value: 2 * 7 * 24, Label: "2w"},
		{Value: 1 * 30 * 24, Label: "1m"},
		{Value: 2 * 30 * 24, Label: "2m"},
		{Value: 6 * 30 * 24, Label: "6m"},
	}

	for _, g := range ordered {
		flex := plot.NewHFlex()

		text := plot.NewTextbox(g.Time.Format("2006-01"))
		text.Style.Origin = plot.Point{X: -1, Y: 0}
		flex.Add(80, text)

		lifetimes := g.Lifetimes()
		for i, v := range lifetimes {
			if v > maxLifetime {
				lifetimes[i] = maxLifetime
			}
		}

		labels := plot.NewTickLabelsX()
		labels.Style.Origin = plot.Point{X: 0, Y: -1}

		density := plot.NewDensity("", plot.DurationTo(lifetimes, time.Hour))
		density.Kernel = 0.05

		line := plot.NewLine("", plot.Ps(
			0, 0.1,
			1, 0.2,
			2, 0.3,
			6, 0.4,
			12, 0.5,
			1*24, 0.6,
			2*24, 0.7,
			7*24, 0.8,
			2*7*24, 0.9,
			1*30*24, 1.0,
			2*30*24, 1.1,
			6*30*24, 1.2,
		))
		_ = line

		flex.AddGroup(0,
			// plot.NewGrid(),
			plot.NewGizmo(),
			density,
			// line,
			labels,
		)
		stack.Add(flex)
	}

	svg := plotsvg.New(400, float64(50*len(ordered)))
	p.Draw(svg)

	os.Stdout.Write(svg.Bytes())
}
