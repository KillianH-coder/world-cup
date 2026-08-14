package matches

import (
	"world-cup/results"
	"world-cup/teams"

	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
)

const defaultScoreboardFile = "scores/default_scoreboard.txt"

func Matches(scores, currentRound map[string]int, team1, team2 string) map[string]int {
	difference := scores[team1] - scores[team2]
	// 0 = 5%/35%, 1 = 10%/50%, 2 = 35%/5%, 3 = 25%/5%, 4 = 16%/4%, 5 = 9%/1%
	// (0: -15%/+15%, 1: -20%/+20%, 2: +15%/-15%, 3: +10%/-10%, 4: +6%/-6%, 5: +4%/-4%)
	gap40 := [10]int{91, 99, 75, 95, 50, 90, 15, 85, 5, 35}
	// 0 = 10%/30%, 1 = 17%/43%, 2 = 30%/10%, 3 = 23%/7%, 4 = 13%/7%, 5 = 7%/3%
	// (0: -10%/+10%, 1: -13%/+13%, 2: +10%/-10%, 3: +8%/-8%,   4: +3%/-3%, 5: +2%/-2%)
	gap20 := [10]int{93, 97, 80, 90, 57, 83, 27, 73, 10, 30}
	// 0 = 15%/25%, 1 = 23%/37%, 2 = 25%/15%, 3 = 19%/11%, 4 = 12%/8%, 5 = 6%/4%
	// (0: -5%/+5%,   1: -7%/+7%,   2: +5%/-5%,   3: +4%/-4%,   4: +2%/-2%, 5: +1%/-1%)
	gap10 := [10]int{94, 96, 82, 88, 63, 77, 38, 62, 15, 25}
	// 0 = 20%/20%, 1 = 30%/30%, 2 = 20%/20%, 3 = 15%/15%, 4 = 10%/10%, 5 = 5%/5%
	defaultGap := [10]int{95, 95, 85, 85, 70, 70, 50, 50, 20, 20}
	switch {
	case difference >= 40:
		currentRound = Goals(currentRound, team1, team2, difference, gap40)
	case difference >= 20:
		currentRound = Goals(currentRound, team1, team2, difference, gap20)
	case difference >= 10:
		currentRound = Goals(currentRound, team1, team2, difference, gap10)
	default:
		currentRound = Goals(currentRound, team1, team2, difference, defaultGap)
	}
	return currentRound
}

func Goals(currentRound map[string]int, team1, team2 string, difference int, gap [10]int) map[string]int {
	currentRound[team1] = rand.IntN(99)
	currentRound[team2] = rand.IntN(99)
	switch {
	case currentRound[team1] >= gap[0]:
		currentRound[team1] = 5
	case currentRound[team1] >= gap[2]:
		currentRound[team1] = 4
	case currentRound[team1] >= gap[4]:
		currentRound[team1] = 3
	case currentRound[team1] >= gap[6]:
		currentRound[team1] = 2
	case currentRound[team1] >= gap[8]:
		currentRound[team1] = 1
	default: // 20/100
		currentRound[team1] = 0
	}
	switch {
	case currentRound[team2] >= gap[1]:
		currentRound[team2] = 5
	case currentRound[team2] >= gap[3]:
		currentRound[team2] = 4
	case currentRound[team2] >= gap[5]:
		currentRound[team2] = 3
	case currentRound[team2] >= gap[7]:
		currentRound[team2] = 2
	case currentRound[team2] >= gap[9]:
		currentRound[team2] = 1
	default: // 20/100
		currentRound[team2] = 0
	}
	return currentRound
}

func Overtime(currentRound map[string]int, team1, team2 string) map[string]int {
	currentRound[team1] += rand.IntN(2)
	currentRound[team2] += rand.IntN(2)
	return currentRound
}

func Penalty(scoreboard string) (int, int, string) {
	var shoot, score1, score2 int
	for range 5 {
		shoot = rand.IntN(3)
		if shoot > 0 {
			score1 += 1
		}
		shoot = rand.IntN(3)
		if shoot > 0 {
			score2 += 1
		}
	}

	// replace penalty goals in scoreboard
	scoreboard = strings.ReplaceAll(scoreboard, "[", strconv.Itoa(score1))
	scoreboard = strings.ReplaceAll(scoreboard, "]", strconv.Itoa(score2))
	if score1 != score2 {
		scoreboard = strings.ReplaceAll(scoreboard, "{", "X")
		scoreboard = strings.ReplaceAll(scoreboard, "}", "X")
	}

	fmt.Println("penalty score: team1", score1, ":", score2, "team2")

	for score1 == score2 {
		shoot = rand.IntN(3)
		if shoot > 0 {
			score1 += 1
		}
		shoot = rand.IntN(3)
		if shoot > 0 {
			score2 += 1
		}
	}
	return score1, score2, scoreboard
}

func MatchResults(
	scores, currentRound, nextRound map[string]int, allTeams []string,
) (map[string]int, map[string]int) {
	// match
	for index := 0; index < len(allTeams); index += 2 {
		var team1 = allTeams[index]
		var team2 = allTeams[index+1]
		var isEqual bool

		// initialize scoreboard
		data, err := os.ReadFile(defaultScoreboardFile)
		if err != nil {
			fmt.Println("Error at reading file", err)
			os.Exit(1)
		}
		scoreboard := string(data) + "\n"
		// var scoreboard string
		// scoreboard += ".__________._______________.__________.    ._______________.    ._______________.    ._______________.\n"
		// scoreboard += "|-<Scores>-|----<Match>----|-<Scores>-| -> |--<Over Time>--| -> |---<Penalty>---| -> |-<Tie Breaker>-|\n"
		// scoreboard += "|==========:===============:==========| -> |===============| -> |===============| -> |===============|\n"
		// scoreboard += "| s1 <- S1 | tm1 # : @ tm2 | S2 -> s2 | -> | tm1 ( : ) tm2 | -> | tm1 [ : ] tm2 | -> | tm1 { : } tm2 |\n"
		// scoreboard += "|__________|_______________|__________| -> |_______________| -> |_______________| -> |_______________|\n\n"

		// replace teams in scoreboard
		scoreboard = strings.ReplaceAll(scoreboard, "tm1", team1)
		scoreboard = strings.ReplaceAll(scoreboard, "tm2", team2)

		// replace scores in scoreboard
		scoreboard = strings.ReplaceAll(scoreboard, "S1", strconv.Itoa(scores[team1]))
		scoreboard = strings.ReplaceAll(scoreboard, "S2", strconv.Itoa(scores[team2]))

		// match
		currentRound = Matches(scores, currentRound, team1, team2)
		nextRound, isEqual = teams.QualifiedTeam(currentRound, nextRound, team1, team2)

		// replace match goals in scoreboard
		scoreboard = strings.ReplaceAll(scoreboard, "#", strconv.Itoa(currentRound[team1]))
		scoreboard = strings.ReplaceAll(scoreboard, "@", strconv.Itoa(currentRound[team2]))
		if isEqual == false {
			scoreboard = strings.ReplaceAll(scoreboard, "(", "X")
			scoreboard = strings.ReplaceAll(scoreboard, ")", "X")
			scoreboard = strings.ReplaceAll(scoreboard, "[", "X")
			scoreboard = strings.ReplaceAll(scoreboard, "]", "X")
			scoreboard = strings.ReplaceAll(scoreboard, "{", "X")
			scoreboard = strings.ReplaceAll(scoreboard, "}", "X")
		}

		fmt.Println("Match result:", team1, currentRound[team1], ":", currentRound[team2], team2)

		// overtime
		if isEqual == true {
			currentRound = Overtime(currentRound, team1, team2)
			nextRound, isEqual = teams.QualifiedTeam(currentRound, nextRound, team1, team2)

			// replace overtime goals in scoreboard
			scoreboard = strings.ReplaceAll(scoreboard, "(", strconv.Itoa(currentRound[team1]))
			scoreboard = strings.ReplaceAll(scoreboard, ")", strconv.Itoa(currentRound[team2]))
			if isEqual == false {
				scoreboard = strings.ReplaceAll(scoreboard, "[", "X")
				scoreboard = strings.ReplaceAll(scoreboard, "]", "X")
				scoreboard = strings.ReplaceAll(scoreboard, "{", "X")
				scoreboard = strings.ReplaceAll(scoreboard, "}", "X")
			}

			fmt.Println("Overtime result:", team1, currentRound[team1], ":", currentRound[team2], team2)

			// penalty
			if isEqual == true {
				score1, score2, newScoreboard := Penalty(scoreboard)
				if score1 > score2 {
					nextRound[team1] = currentRound[team1]
				} else {
					nextRound[team2] = currentRound[team2]
				}

				// replace final penalty goals in scoreboard
				scoreboard = strings.ReplaceAll(newScoreboard, "{", strconv.Itoa(score1))
				scoreboard = strings.ReplaceAll(scoreboard, "}", strconv.Itoa(score2))

				fmt.Println("final penalty score:", team1, score1, ":", score2, team2)
			}
		}

		// new scores
		if _, ok := nextRound[team1]; ok {
			scores = teams.NewScore(scores, currentRound, team1, team2)
		} else {
			scores = teams.NewScore(scores, currentRound, team2, team1)
		}

		// replace new scores in scoreboard
		scoreboard = strings.ReplaceAll(scoreboard, "s1", strconv.Itoa(scores[team1]))
		scoreboard = strings.ReplaceAll(scoreboard, "s2", strconv.Itoa(scores[team2]))

		// update scoreboard.txt
		results.ScoreboardResults(scoreboard)
		fmt.Println(scoreboard)
	}
	return currentRound, nextRound
}
