package match

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	uuid "github.com/nu7hatch/gouuid"
	"strconv"
	"time"
)

type Match struct {
	MatchId      string `dynamo:"match_id"`
	SortKey      string `dynamo:"sort_key"`
	Time         int64  `dynamo:"time"`
	NumberOfSets int    `dynamo:"number_of_sets"`
	NumberOfLegs int    `dynamo:"number_of_legs"`
	StartScore   int    `dynamo:"start_score"`
	CurrentSet   int    `dynamo:"current_set"`
	CurrentLeg   int    `dynamo:"current_leg"`
	Winner       string `dynamo:"winner"`
}

type Set struct {
	MatchId string `dynamo:"match_id"`
	SortKey string `dynamo:"sort_key"`
	Time    int64  `dynamo:"time"`
	Winner  string `dynamo:"winner"`
	Number  int    `dynamo:"number"`
}

type Leg struct {
	MatchId string      `dynamo:"match_id"`
	SortKey string      `dynamo:"sort_key"`
	Time    int64       `dynamo:"time"`
	Players []PlayerLeg `dynamo:"players"`
	Winner  string      `dynamo:"winner"`
	Number  int         `dynamo:"number"`
}

type PlayerLeg struct {
	PlayerId string `dynamo:"player_id"`
	Score    int    `dynamo:"score"`
}

type Throw struct {
	MatchId  string `dynamo:"match_id"`
	SortKey  string `dynamo:"sort_key"`
	Time     int64  `dynamo:"time"`
	Number   int    `dynamo:"number"`
	PlayerId string `dynamo:"player_id"`
	Score    int    `dynamo:"score"`
}

type ThrowIdentity struct {
	MatchId string
	SetId   string
	LegId   string
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

func (s Service) Create(match Match, players []string) (Match, error) {
	match.MatchId = CreateId()
	match.SortKey = "INFO"
	match.CurrentSet = 1
	match.CurrentLeg = 1
	match.Time = time.Now().Unix()
	err := s.table.Put(match).Run()
	if err != nil {
		return Match{}, err
	}

	set := Set{
		MatchId: match.MatchId,
		SortKey: "SET" + strconv.Itoa(match.CurrentSet),
		Winner:  "",
		Number:  1,
		Time:    match.Time,
	}
	err = s.table.Put(set).Run()
	if err != nil {
		return Match{}, err
	}

	var playerSlice []PlayerLeg

	for _, player := range players {
		playerSlice = append(playerSlice, PlayerLeg{
			PlayerId: player,
			Score:    match.StartScore,
		})
	}

	leg := Leg{
		MatchId: match.MatchId,
		SortKey: set.SortKey + "#LEG" + strconv.Itoa(match.CurrentLeg),
		Winner:  "",
		Players: playerSlice,
		Number:  1,
		Time:    match.Time,
	}
	err = s.table.Put(leg).Run()
	if err != nil {
		return Match{}, err
	}

	return match, nil
}

func (s Service) AddThrow(match Match, throw Throw) (Throw, error) {
	throw.MatchId = match.MatchId
	throw.SortKey =
		"SET" + strconv.Itoa(match.CurrentSet) + "#" +
			"LEG" + strconv.Itoa(match.CurrentLeg) + "#" +
			"THROW" + strconv.Itoa(throw.Number)
	throw.Time = time.Now().Unix()

	return throw, s.table.Put(throw).
		If("attribute_not_exists(match_id) AND attribute_not_exists(sort_key)", throw.MatchId, throw.SortKey).
		Run()
}

func (s Service) GetInfoById(matchId string) (match Match, err error) {
	err = s.table.Get("match_id", matchId).
		Range("sort_key", dynamo.Equal, "INFO").
		One(&match)

	if err != nil {
		return Match{}, err
	}
	return
}

func (s Service) GetById(matchId string) (match []Match, err error) {
	err = s.table.Get("match_id", matchId).All(&match)

	if err != nil {
		return []Match{}, err
	}
	return
}

func CreateId() (id string) {
	v4, _ := uuid.NewV4()
	return v4.String()
}
