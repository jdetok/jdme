package api

import "strconv"

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
func HandleSeasonId(sId uint64, p *Player) uint64 {
	if strconv.FormatUint(sId, 10)[1:] == "9999" { // agg seasons
		return sId
	} else if sId >= 80000 && sId < 90000 {
		return p.SeasonIdMax // return most recent season
	} else if sId >= 70000 && sId < 80000 {
		return p.PSeasonIdMax // return most recent season
	} else if sId >= 40000 && sId < 50000 {
		if p.PSeasonIdMax < 40000 { // player has no playeroff, return max reg season
			return p.SeasonIdMax // return reg season if player has no playoffs
		}
		if sId == 49999 { // playoff career
			return sId
		}

		if sId > p.PSeasonIdMax {
			return p.PSeasonIdMax
		}
		if sId < p.PSeasonIdMin {
			return p.PSeasonIdMin
		}
	} else if sId >= 20000 && sId < 30000 {
		if sId > p.SeasonIdMax {
			if sId == 29999 { // reg season career
				return sId
			}
			return p.SeasonIdMax
		}
		if sId < p.SeasonIdMin {
			return p.SeasonIdMin
		}
	}
	return sId
}

/*
accept the slice of all players and a seasonId, return a slice with just the
active players from the passed season id
*/
// func SlicePlayersSzn(players []Player, seasonId uint64) ([]Player, error) {
func SlicePlayersSzn(players []Player, seasonId uint64) ([]Player, error) {
	var plslice []Player

	sl := LgSznsByMonth()

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
		// append players to the random slice if the passed season id between player min and max season
		if (seasonId >= 20000 && seasonId < 30000) &&
			(seasonId <= p.SeasonIdMax && seasonId >= p.SeasonIdMin) ||
			(seasonId >= 40000 && seasonId < 50000) &&
				(seasonId <= p.PSeasonIdMax && seasonId >= p.PSeasonIdMin) {
			plslice = append(plslice, p)
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
