package memd

import (
	"fmt"

	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
)

// update players, seasons, and teams in memory structs slices
func (s *InMemStore) Rebuild(db pgdb.DB) error {
	var err error
	s.Players, err = UpdatePlayers(db)
	if err != nil {
		return fmt.Errorf("failed updating players structs: %v", err)
	}
	s.Seasons, err = UpdateSeasons(db)
	if err != nil {
		return fmt.Errorf("failed updating season structs: %v", err)
	}
	s.Teams, err = UpdateTeams(db)
	if err != nil {
		return fmt.Errorf("failed updating team structs: %v", err)
	}

	s.TeamRecs, err = UpdateTeamRecords(db, &s.CurrentSzns)
	if err != nil {
		return fmt.Errorf("failed updating team records: %v", err)
	}
	s.TopLgPlayers, err = QueryTopLgPlayers(db, &s.CurrentSzns, "100")
	if err != nil {
		return fmt.Errorf("failed updating league top players struct: %v", err)
	}
	return nil
}
