package teams

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func BuildTeams(
	groups [][]string, teamsScore, group map[string]int,
) (map[string]int, map[string]int, []string) {
	var newGroups []string
	for _, value := range groups {
		for _, chars := range value {
			team, score, sep := strings.Cut(chars, ":")
			if !sep {
				fmt.Println("No ':' found in string")
				os.Exit(1)
			}
			// Building score map
			teamsScore[team], _ = strconv.Atoi(score)
			// Building group map
			group[team] = 0
			// Building groups slice
			newGroups = append(newGroups, team)
		}
	}
	return teamsScore, group, newGroups
}

func QualifiedTeam(currentRound, nextRound map[string]int, team1, team2 string) (map[string]int, bool) {
	switch {
	case currentRound[team1] > currentRound[team2]:
		nextRound[team1] = currentRound[team1]
	case currentRound[team2] > currentRound[team1]:
		nextRound[team2] = currentRound[team2]
	default:
		return nextRound, true
	}
	return nextRound, false
}

func NewScore(scores map[string]int, currentRound map[string]int, winner string, looser string) map[string]int {
	var k, newScore float64
	var goalsGap = float64(currentRound[winner] - currentRound[looser])

	fmt.Println("old score : team 1 =", scores[winner], ":", scores[looser], "= team2")

	if scores[winner] < scores[looser] {
		k = 0.5 * math.Sqrt(float64(scores[winner])/float64(scores[looser]))                                  // (team avec le moins de point gagne) k = 0.5 * racine(score[plus bas] / score[plus haut])
		newScore = float64(scores[looser]) - float64(scores[looser]-scores[winner])*math.Exp(-(1+goalsGap)*k) // nouveau score = score[plus haut] - (score[plus haut] - score[plus bas])*exp(-(1+écart buts)*k)

	} else if scores[winner] >= scores[looser] {
		k = 0.1 * math.Sqrt(float64(scores[looser])/float64(scores[winner])) // (team avec le plus de score gagne) k = 0.1 * racine(score[plus bas] / score[plus haut])
		newScore = 99 - float64(99-scores[winner])*math.Exp(-(1+goalsGap)*k) // nouveau score = 99 - (99 - score[plus bas])*exp(-(1+écart buts)*k)
	}

	// new scores
	oldWinner := scores[winner]
	newWinner := int(math.Round(newScore))

	delta := newWinner - oldWinner

	scores[winner] = newWinner
	scores[looser] -= delta

	fmt.Println("new score : team 1 =", scores[winner], ":", scores[looser], "= team2")

	return scores
}
