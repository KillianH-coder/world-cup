package results

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	scoreboardFile         = "scores/scoreboard.txt"
	defaultLeaderboardFile = "scores/default_leaderboard.txt"
)

func LeaderboardResults(
	last32, round16, quarterFinal, semiFinal, final map[string]int,
	teamsLast32, teamsRound16, teamsQuarterFinal, teamsSemiFinal, teamsFinal, winnerTeam []string,
) string {
	data, err := os.ReadFile(defaultLeaderboardFile)
	if err != nil {
		fmt.Println("Error at reading file *", err)
		os.Exit(1)
	}
	finalResult := string(data)

	// Last 32 results
	j := 0
	for i := 111; i <= 142; i++ {
		finalResult = strings.ReplaceAll(finalResult, strconv.Itoa(i), teamsLast32[j])
		j++
	}
	j = 0
	for i := 'À'; i <= 'ß'; i++ {
		finalResult = strings.ReplaceAll(finalResult, string(i), strconv.Itoa(last32[teamsLast32[j]]))
		j++
	}
	// Round 16 results
	j = 0
	for i := 211; i <= 226; i++ {
		finalResult = strings.ReplaceAll(finalResult, strconv.Itoa(i), teamsRound16[j])
		j++
	}
	j = 0
	for i := 'à'; i <= 'ï'; i++ {
		finalResult = strings.ReplaceAll(finalResult, string(i), strconv.Itoa(round16[teamsRound16[j]]))
		j++
	}
	// Quarter final results
	j = 0
	for i := 311; i <= 318; i++ {
		finalResult = strings.ReplaceAll(finalResult, strconv.Itoa(i), teamsQuarterFinal[j])
		j++
	}
	j = 0
	for i := 'ð'; i <= '÷'; i++ {
		finalResult = strings.ReplaceAll(finalResult, string(i), strconv.Itoa(quarterFinal[teamsQuarterFinal[j]]))
		j++
	}
	// Semi final results
	j = 0
	for i := 411; i <= 414; i++ {
		finalResult = strings.ReplaceAll(finalResult, strconv.Itoa(i), teamsSemiFinal[j])
		j++
	}
	j = 0
	for i := '¡'; i <= '¤'; i++ {
		finalResult = strings.ReplaceAll(finalResult, string(i), strconv.Itoa(semiFinal[teamsSemiFinal[j]]))
		j++
	}
	// Final results
	j = 0
	for i := 511; i <= 512; i++ {
		finalResult = strings.ReplaceAll(finalResult, strconv.Itoa(i), teamsFinal[j])
		j++
	}
	j = 0
	for i := '*'; i <= '+'; i++ {
		finalResult = strings.ReplaceAll(finalResult, string(i), strconv.Itoa(final[teamsFinal[j]]))
		j++
	}
	// Winner
	finalResult = strings.ReplaceAll(finalResult, "611", winnerTeam[0])
	return finalResult
}

func ScoreboardResults(scoreboard string) {
	file, err := os.OpenFile(scoreboardFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error at oppening file", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.WriteString(scoreboard)
	if err != nil {
		fmt.Println("Error at writing to file", err)
		os.Exit(1)
	}
}
