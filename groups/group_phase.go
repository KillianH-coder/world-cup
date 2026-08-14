package groups

import (
	"world-cup/matches"

	"fmt"
	"os"
	"sort"
	"strings"
)

type Team struct {
	Name            string
	Points          int
	GoalsDifference int
}

func GroupsTeams(file string) [][]string {
	var groups [][]string
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error at reading teams.txt |", err)
		os.Exit(1)
	}
	datas := strings.Split(string(data), "\n")
	for _, group := range datas {
		groups = append(groups, strings.Split(group, "-"))
	}
	return groups
}

func buildTeams(teams []string) []Team {
	result := make([]Team, len(teams))
	for index, name := range teams {
		result[index] = Team{Name: name}
	}
	return result
}

func GroupsPhase(teamsScore, teamsGroup map[string]int, teamsLast32, teams []string) []string {
	allTeams := buildTeams(teams)

	for index := 0; index < len(teams); index += 4 {
		teamA := &allTeams[index]
		teamB := &allTeams[index+1]
		teamC := &allTeams[index+2]
		teamD := &allTeams[index+3]

		teamsGroup = PoolMatches(teamsScore, teamsGroup, teamA, teamB)
		teamsGroup = PoolMatches(teamsScore, teamsGroup, teamA, teamC)
		teamsGroup = PoolMatches(teamsScore, teamsGroup, teamA, teamD)
		teamsGroup = PoolMatches(teamsScore, teamsGroup, teamB, teamC)
		teamsGroup = PoolMatches(teamsScore, teamsGroup, teamB, teamD)
		teamsGroup = PoolMatches(teamsScore, teamsGroup, teamC, teamD)
	}
	sort.SliceStable(allTeams, func(i, j int) bool {
		if allTeams[i].Points != allTeams[j].Points {
			return allTeams[i].Points > allTeams[j].Points
		}
		return allTeams[i].GoalsDifference > allTeams[j].GoalsDifference
	})
	for i := 0; i < 32; i += 2 {
		teamsLast32 = append(teamsLast32, allTeams[i].Name)
		teamsLast32 = append(teamsLast32, allTeams[31-i].Name)
	}
	return teamsLast32
}

func PoolMatches(teamsScore, teamsGroup map[string]int, a, b *Team) map[string]int {
	teamsGroup = matches.Matches(teamsScore, teamsGroup, a.Name, b.Name)

	scoreA := teamsGroup[a.Name]
	scoreB := teamsGroup[b.Name]
	a.GoalsDifference += scoreA - scoreB
	b.GoalsDifference += scoreB - scoreA

	switch {
	case scoreA > scoreB:
		a.Points += 3
	case scoreB > scoreA:
		b.Points += 3
	default:
		a.Points += 1
		b.Points += 1
	}
	return teamsGroup
}
