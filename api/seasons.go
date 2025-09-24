package api

import (
	"fmt"
	"time"
)

/*
accept a season id and a pointer to a Player struct, validate the player was active
in the passed season, return a valid season ID if not. if season id starts with an
8 the player's max regular season will be returned. if it starts with a 7, their
max playoff season will be returned. if it starts with a 4, it will first verify
player has played in a playoff game, and will return their max regular season if
they haven't. a season id starting with 2 will return a regular season. for both
regular season and playoffs, the function will verify the player played in said
season, and return either their max or min (whichever is closer) season  if they
did not
*/
func HandleSeasonId(sId uint64, p *Player, errStr *string) uint64 {
	if sId == 99999 || sId == 29999 { // agg seasons
		msg := fmt.Sprintf("aggregate season requested%d | %d\n", sId, sId)
		fmt.Println(msg)
		return sId
	} else if sId == 88888 {
		msg := fmt.Sprintf("returning latest regular season for player%d | %d\n",
			sId, p.SeasonIdMax)
		fmt.Println(msg)
		return p.SeasonIdMax // return most recent season
	} else if sId >= 40000 && sId < 50000 {
		if p.PSeasonIdMax < 40000 { // player has no playeroff, return max reg season
			msg := fmt.Sprintf(
				"%s has not played in the post-season | displaying latest regular season stats",
				p.Name)
			*errStr = msg
			fmt.Println(msg)
			return p.SeasonIdMax // return reg season if player has no playoffs
		}
		if sId == 49999 {
			msg := fmt.Sprintf(
				"requested career playoff stats %d | %d\n",
				sId, sId)
			// *errStr = msg
			fmt.Println(msg)
			return sId
		}
		if sId > p.PSeasonIdMax {
			msg := fmt.Sprintf(
				// "szn > playoff max, returning playoff max%d | %d\n",
				// sId, p.PSeasonIdMax)
				"%d was after %s's last playoff season | displaying the %d playoffs",
				sId, p.Name, p.PSeasonIdMax)
			*errStr = msg
			fmt.Println(msg)
			return p.PSeasonIdMax
		}
		if sId < p.PSeasonIdMin {
			msg := fmt.Sprintf(
				"the first playoffs for %s was the %d season",
				p.Name, p.PSeasonIdMin)
			*errStr = msg
			fmt.Println(msg)
			return p.PSeasonIdMin
		}
	} else if sId >= 20000 && sId < 30000 {
		if sId > p.SeasonIdMax {
			msg := fmt.Sprintf(
				"%s has not played games in the %d season | displaying %d stats instead\n",
				p.Name, sId, p.SeasonIdMax)
			*errStr = msg
			fmt.Println(msg)
			return p.SeasonIdMax
		}
		if sId < p.SeasonIdMin {
			msg := fmt.Sprintf(
				"%s was not in the league yet for the %d season | displaying their rookie season %d stats instead\n",
				p.Name, sId, p.SeasonIdMin)
			*errStr = msg
			fmt.Println(msg)
			return p.SeasonIdMin
		}
	}
	msg := fmt.Sprintf("validated: %d | %d\n", sId, sId)
	fmt.Println(msg)
	return sId
}

/*
accept the slice of all players and a seasonId, return a slice with just the
active players from the passed season id
*/
// func SlicePlayersSzn(players []Player, seasonId uint64) ([]Player, error) {
func SlicePlayersSzn(players []Player, seasonId uint64, lg string) ([]Player, error) {
	var plslice []Player

	// get struct with current seasons
	sl := LgSznsByMonth(time.Now())

	for _, p := range players { // EXPAND THIS IF TO CATCH PLAYOFF SEASONS AS WELL

		// handle random season id
		if seasonId == 88888 {
			switch p.League {
			case "nba":
				seasonId = sl.SznId
			case "wnba":
				seasonId = sl.WSznId
			}
		}

		if seasonId == 49999 {
			if p.PSeasonIdMin > 0 && p.SeasonIdMax >= (sl.WSznId-3) {
				if lg == "all" || lg == p.League {
					plslice = append(plslice, p)
				}
			}
		}

		if seasonId == 29999 {
			if p.SeasonIdMax >= (sl.WSznId - 3) {
				if lg == "all" || lg == p.League {
					plslice = append(plslice, p)
				}
			}
		}

		// append players to the random slice if the passed season id between player min and max season
		if seasonId >= 20000 && seasonId < 30000 {
			if seasonId <= p.SeasonIdMax && seasonId >= p.SeasonIdMin {
				if lg == "all" || lg == p.League {
					plslice = append(plslice, p)
				}
			}
		}

		if seasonId >= 40000 && seasonId < 50000 {
			if seasonId <= p.PSeasonIdMax && seasonId >= p.PSeasonIdMin {
				if lg == "all" || lg == p.League {
					plslice = append(plslice, p)
				}
			}
		}

	}
	return plslice, nil
}

// accept pointers of league and season, switch season/wseason on league
func (t *RespSeasonTmp) HndlSeason(league *string, season *string) {
	switch *league {
	case "NBA":
		*season = t.Season
	case "WNBA":
		*season = t.WSeason
	}
}
