package match

import (
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
	Legs   []Leg  `dynamo:"legs"`
	Winner string `dynamo:"winner"`
}

type Leg struct {
	Throws []Throw `dynamo:"throws"`
	Winner string  `dynamo:"winner"`
}

type Throw struct {
	UserId string `dynamo:"user"`
	Score  int    `dynamo:"score"`
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
	return match, s.table.Put(match).Run()
}

func CreateId() (id string) {
	v4, _ := uuid.NewV4()
	return v4.String()
}
