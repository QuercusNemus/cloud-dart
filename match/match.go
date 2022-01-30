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
	Id           string `json:"id"`
	SK           string `json:"sk"`
	NumberOfSets int    `json:"number_of_sets"`
	NumberOfLegs int    `json:"number_of_legs"`
	StartScore   int    `json:"start_score"`
	CurrentSet   int    `json:"current_set"`
	CurrentLeg   int    `json:"current_leg"`
	Winner       string `json:"winner"`
	Time         int64  `json:"time"`
}

type Set struct {
	Id     string `json:"id"`
	SK     string `json:"sk"`
	Winner string `json:"winner"`
	Number int    `json:"number"`
	Time   int64  `json:"time"`
}

type Leg struct {
	Id      string      `json:"id"`
	SK      string      `json:"sk"`
	Players []PlayerLeg `json:"players"`
	Winner  string      `json:"winner"`
	Number  int         `json:"number"`
	Time    int64       `json:"time"`
}

type PlayerLeg struct {
	PlayerId string `json:"player_id"`
	Score    int    `json:"score"`
}

type Throw struct {
	UserId string `json:"user"`
	Score  int    `json:"score"`
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
	match.Id = CreateId()
	match.SK = "INFO"
	match.CurrentSet = 1
	match.CurrentLeg = 1
	match.Time = time.Now().Unix()
	err := s.table.Put(match).Run()
	if err != nil {
		return Match{}, err
	}

	set := Set{
		Id:     match.Id,
		SK:     "SET" + strconv.Itoa(match.CurrentSet) + "#",
		Winner: "",
		Number: 1,
		Time:   match.Time,
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
		Id:      match.Id,
		SK:      set.SK + "LEG" + strconv.Itoa(match.CurrentLeg) + "#",
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

func (s Service) GetById(matchId string) (match Match, err error) {
	err = s.table.Get("id", matchId).One(&match)
	if err != nil {
		return Match{}, err
	}
	return
}

func CreateId() (id string) {
	v4, _ := uuid.NewV4()
	return v4.String()
}
