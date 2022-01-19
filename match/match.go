package match

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	uuid "github.com/nu7hatch/gouuid"
)

type Match struct {
	Id           string   `dynamo:"id"`
	Players      []string `dynamo:"players"`
	NumberOfSets int      `dynamo:"number_of_sets"`
	Sets         []Set    `dynamo:"sets"`
	NumberOfLegs int      `dynamo:"number_of_legs"`
	Startscore   int      `dynamo:"startscore"`
	Winner       string   `dynamo:"winner"`
}

type Set struct {
	Id     string `dynamo:"id"`
	Legs   []Leg  `dynamo:"legs"`
	Winner string `dynamo:"winner"`
}

type Leg struct {
	Id     string  `dynamo:"id"`
	Throws []Throw `dynamo:"throws"`
	Winner string  `dynamo:"winner"`
}

type Throw struct {
	UserId string `dynamo:"user"`
	Score  int    `dynamo:"score"`
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

func (s Service) Create(match Match) (Match, error) {
	match.Id = CreateId()
	leg := Leg{
		Id: CreateId(),
	}
	set := Set{
		Id:   CreateId(),
		Legs: []Leg{leg},
	}
	match.Sets = append(match.Sets, set)
	return match, s.table.Put(match).Run()
}

func (s Service) GetById(matchId string) (match Match, err error) {
	err = s.table.Get("id", matchId).One(&match)
	if err != nil {
		return Match{}, err
	}
	return
}

func (s Service) AddThrow(throw Throw, identity ThrowIdentity) (Throw, error) {
	m, err := s.GetById(identity.MatchId)
	if err != nil {
		return Throw{}, errors.New("unable to find given match by id")
	}
	for _, set := range m.Sets {
		if set.Id == identity.SetId {
			for _, leg := range set.Legs {
				if leg.Id == identity.LegId {
					leg.Throws = append(leg.Throws, throw)
				} else {
					return Throw{}, errors.New("unable to find given leg in this match")
				}
			}
		} else {
			return Throw{}, errors.New("unable to find given set in this match")
		}
	}
	return throw, s.table.Update("id", m.Id).Set("match", m).Value(&m)
}

func CreateId() (id string) {
	v4, _ := uuid.NewV4()
	return v4.String()
}
