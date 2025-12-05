package resp

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jdetok/go-api-jdeko.me/pkg/clnd"
	"github.com/jdetok/go-api-jdeko.me/pkg/errd"
	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
)

// read player from query string and clean the value (remove accents, lowercase)
func PlayerFromQ(r *http.Request, m *memd.StMaps) (uint64, error) {
	pStr := r.URL.Query().Get("player")

	// check if integer (return if conversion successful)
	plrIdInt, err := strconv.ParseUint(pStr, 10, 64)
	if err == nil {
		return plrIdInt, nil
	}

	// clean string & remove accents on letters (all standard ascii)
	cleanPlrReq := CleanPlayerReq(pStr)

	fmt.Println("in player_req: ", pStr)
	// check if requested player name exists in maps, reutrn player id if so
	plrIdUint, err := m.GetPlrIdFromName(cleanPlrReq)
	if err != nil {
		return 0, &errd.ValidationError{Val: pStr}
	}

	return plrIdUint, nil
	// fmt.Printf("player request (raw): %s | cleaned: %s\n", pStr, p_cln)
	// return p_cln
}
func CleanPlayerReq(pname string) string {
	p_lwr := strings.ToLower(pname)
	return clnd.ConvToASCII(p_lwr)
}

// accept http request, get the "season" passed in the query string, return as int
func SeasonFromQ(r *http.Request, maxRSzn, maxPSzn int) (int, error) {
	s := r.URL.Query().Get("season")
	s_int, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("INVALID SEASON: could not convert %s to an int\n%s",
			s, err.Error())
	}

	firstDig := strconv.Itoa(s_int)[0]
	// lastDigs := strconv.Itoa(s_int)[1:]
	var maxSzn int
	switch firstDig {
	case '2':
		maxSzn = maxRSzn
	case '4':
		maxSzn = maxPSzn
	case '8':
		break
	default:
		return 0, fmt.Errorf("INVALID SEASON: %s | must begin with 2, 4, or 8  to an int\n%v",
			s, err)
	}

	if s_int > maxSzn && s_int != 88888 && s_int != 29999 && s_int != 49999 {
		return 0, fmt.Errorf("INVALID SEASON: %s | must be less than %d", s,
			maxSzn)
	}

	return s_int, nil
}

// returns team arg from query string
// as uint64 if passed a teamId or as string if passed abbr
// returns a 0 if error occurs or team is not included in query string
// if the returned value is a string, the caller must use the league,
// either from the league argument or from the player's league, to get the
// team id as a uint64
func TeamFromQ(r *http.Request, m *memd.StMaps) (any, error) {
	t := r.URL.Query().Get("team")
	if t != "" {
		// handle request for team abbr
		if _, err := strconv.Atoi(t); err != nil {
			if tmId, ok := m.TmAbbrId[t]; !ok {
				return t, fmt.Errorf("couldn't process request for team %v", t)
			} else {
				return tmId, nil // return string team abbr
			}
		}
		// handle request for team id
		teamId, err := m.GetTeamIDUintCC(t)
		if err != nil {
			return uint64(0), err
		}
		return teamId, nil
	}
	return uint64(0), nil
}

func LgFromQ(r *http.Request) (int, error) {
	lg := r.URL.Query().Get("league")
	lgId, err := strconv.Atoi(lg)
	if err != nil {
		switch lg {
		case "all", "":
			return 10, nil
		case "nba":
			return 0, nil
		case "wnba":
			return 1, nil
		default:
			return 99999, err
		}
	}
	return lgId, nil
}
