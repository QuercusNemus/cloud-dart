package match

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	uuid "github.com/nu7hatch/gouuid"
	"time"
)

type Match struct {
	MatchId      string `dynamo:"match_id"`
	SortKey      string `dynamo:"sort_key"`
	Time         int64  `dynamo:"time"`
	NumberOfSets int    `dynamo:"number_of_sets"`
	NumberOfLegs int    `dynamo:"number_of_legs"`
	StartScore   int    `dynamo:"start_score"`
	Winner       string `dynamo:"winner"`
	Sets         []Set  `dynamo:"sets"`
}

type Set struct {
	Time   int64  `dynamo:"time"`
	Winner string `dynamo:"winner"`
	Legs   []Leg  `dynamo:"legs"`
}

type Leg struct {
	Time    int64       `dynamo:"time"`
	Players []PlayerLeg `dynamo:"players"`
	Winner  string      `dynamo:"winner"`
	Number  int         `dynamo:"number"`
	Throws  []Throw     `dynamo:"throws"`
}

type PlayerLeg struct {
	PlayerId string `dynamo:"player_id"`
	Score    int    `dynamo:"score"`
}

type Throw struct {
	Time     int64  `dynamo:"time"`
	Number   int    `dynamo:"number"`
	PlayerId string `dynamo:"player_id"`
	Score    int    `dynamo:"score"`
}

type MatchPlayer struct {
	MatchId  string `dynamo:"match_id"`
	SortKey  string `dynamo:"sort_key"`
	PlayerId string `dynamo:"player_id"`
}

type Service struct {
	table dynamo.Table
}

func NewService(tableName, region string) *Service {
	sess, _ := session.NewSession()
	db := dynamo.New(sess, &aws.Config{Region: aws.String(region)})
	table := db.Table(tableName)

	return &Service{table: table}
}

//Save takes a match and a slices of players and save to the match dynamo table.
//If the game hase an id allready it will only save the match, other wise it will
//save each player from the player slices to the table as a reference.
func (s Service) Save(match Match, players []string) (Match, error) {
	if match.MatchId == "" {
		match.MatchId = CreateId()
		match.SortKey = "MATCH"
		match.Time = time.Now().Unix()

		for _, player := range players {
			matchPlayer := MatchPlayer{
				MatchId:  match.MatchId,
				SortKey:  "PLAYER#" + player,
				PlayerId: player,
			}
			err := s.table.Put(matchPlayer).Run()
			if err != nil {
				return Match{}, err
			}
		}
	}

	err := s.table.Put(match).Run()
	if err != nil {
		return Match{}, err
	}

	return match, nil
}

//GetById takes a Match ID and returns all entries in the table corresponding to that ID.
func (s Service) GetById(matchId string) (match []Match, err error) {
	err = s.table.Get("match_id", matchId).All(&match)

	if err != nil {
		return []Match{}, err
	}
	return
}

//GetByPlayerId takes a Player ID and return entries corrsponding to that index.
func (s Service) GetByPlayerId(playerId string) (matches []MatchPlayer, err error) {
	err = s.table.Get("player_id", playerId).Index("player_id").All(&matches)
	if err != nil {
		return nil, err
	}
	return
}

//CreateId returs a string of UUID V4
func CreateId() (id string) {
	v4, _ := uuid.NewV4()
	return v4.String()
}
