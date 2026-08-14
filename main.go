package main

import (
	"world-cup/groups"
	"world-cup/matches"
	"world-cup/results"
	"world-cup/teams"

	"fmt"
	"os"
)

const (
	teamsFile      = "groups/teams.txt"
	scoreboardFile = "scores/scoreboard.txt"
)

func main() {
	// teams slice for each round (teams in order)
	var teamsGroups, teamsLast32, teamsRound16, teamsQuarterFinal, teamsSemiFinal, teamsFinal, winnerTeam []string

	// teams in maps with there goals in the last match
	teamsScore := make(map[string]int)
	group := make(map[string]int)
	last32 := make(map[string]int)
	round16 := make(map[string]int)
	quarterFinal := make(map[string]int)
	semiFinal := make(map[string]int)
	final := make(map[string]int)
	winner := make(map[string]int)

	// reset scoreboard
	err := os.WriteFile(scoreboardFile, nil, 0644)
	if err != nil {
		fmt.Println("Error at writing in file |", err)
		os.Exit(1)
	}

	// read teams in the .txt file
	allGroups := groups.GroupsTeams(teamsFile)

	// set teams scores map, teams group phase map and group slice
	teamsScore, group, teamsGroups = teams.BuildTeams(allGroups, teamsScore, group)

	// do the group matches and keep teams qualified in last 32
	teamsLast32 = groups.GroupsPhase(teamsScore, group, teamsLast32, teamsGroups)

	// build the last 32 map
	for _, team := range teamsLast32 {
		last32[team] = 0
	}

	fmt.Println("last 32 scores:", last32)
	fmt.Println("last 32 teams:", teamsLast32)
	fmt.Println("number of teams:", len(last32), "/", len(teamsLast32))

	// do the last 32 matches
	fmt.Println("----- Last 32 Matches -----")
	results.ScoreboardResults("----------------------------------------------- Last32 -----------------------------------------------\n")
	last32, round16 = matches.MatchResults(teamsScore, last32, round16, teamsLast32)

	// qualified teams slice for round 16
	for _, team := range teamsLast32 {
		if _, ok := round16[team]; ok {
			teamsRound16 = append(teamsRound16, team)
		}
	}
	fmt.Println("round 16 scores:", round16)
	fmt.Println("round 16 teams:", teamsRound16)
	fmt.Println("number of teams:", len(round16), "/", len(teamsRound16))

	// do the round16 matches
	fmt.Println("----- Round 16 Matches -----")
	results.ScoreboardResults("---------------------------------------------- Round 16 ----------------------------------------------\n")
	round16, quarterFinal = matches.MatchResults(teamsScore, round16, quarterFinal, teamsRound16)

	// qualified teams slice for quarter final
	for _, team := range teamsRound16 {
		if _, ok := quarterFinal[team]; ok {
			teamsQuarterFinal = append(teamsQuarterFinal, team)
		}
	}
	fmt.Println("quarter final scores:", quarterFinal)
	fmt.Println("quarter final teams:", teamsQuarterFinal)
	fmt.Println("number of teams:", len(quarterFinal), "/", len(teamsQuarterFinal))

	// do the quarter final matches
	fmt.Println("----- Quarter Final Matches -----")
	results.ScoreboardResults("-------------------------------------------- QuarterFinal --------------------------------------------\n")
	quarterFinal, semiFinal = matches.MatchResults(teamsScore, quarterFinal, semiFinal, teamsQuarterFinal)

	// qualified teams slice for semi final
	for _, team := range teamsQuarterFinal {
		if _, ok := semiFinal[team]; ok {
			teamsSemiFinal = append(teamsSemiFinal, team)
		}
	}
	fmt.Println("semi final scores:", semiFinal)
	fmt.Println("semi final teams:", teamsSemiFinal)
	fmt.Println("number of teams:", len(semiFinal), "/", len(teamsSemiFinal))

	// do the semi final matches
	fmt.Println("----- Semi Final Matches -----")
	results.ScoreboardResults("--------------------------------------------- Semi Final ---------------------------------------------\n")
	semiFinal, final = matches.MatchResults(teamsScore, semiFinal, final, teamsSemiFinal)

	// qualified teams slice for final
	for _, team := range teamsSemiFinal {
		if _, ok := final[team]; ok {
			teamsFinal = append(teamsFinal, team)
		}
	}
	fmt.Println("final scores:", final)
	fmt.Println("final teams:", teamsFinal)
	fmt.Println("number of teams:", len(final), "/", len(teamsFinal))

	// do the final matches
	fmt.Println("----- Final Matches -----")
	results.ScoreboardResults("--------------------------------------------- GrandFinal ---------------------------------------------\n")
	final, winner = matches.MatchResults(teamsScore, final, winner, teamsFinal)

	// winner team slice
	for _, team := range teamsFinal {
		if _, ok := winner[team]; ok {
			winnerTeam = append(winnerTeam, team)
		}
	}
	fmt.Println("winner team:", winner)
	fmt.Println("number of teams:", len(winner), "/", len(winnerTeam))

	finalResult := results.LeaderboardResults(
		last32, round16, quarterFinal, semiFinal, final,
		teamsLast32, teamsRound16, teamsQuarterFinal, teamsSemiFinal, teamsFinal, winnerTeam)

	// new leaderboard
	err = os.WriteFile(scoreboardFile, []byte(finalResult), 0644)
	if err != nil {
		fmt.Println("Error at writing to file", err)
		os.Exit(1)
	}
}
