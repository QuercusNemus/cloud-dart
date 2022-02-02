package player

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	uuid "github.com/nu7hatch/gouuid"
)

type Player struct {
	Age      int      `dynamo:"age"`
	Email    string   `dynamo:"email"`
	Id       string   `dynamo:"player_id"`
	Name     string   `dynamo:"name"`
	NickName string   `dynamo:"nick_name"`
	Matches  []string `dynamo:"matches"`
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

func (s Service) Create(player Player) (Player, error) {
	users, err := s.GetByEmail(player.Email)
	if err != nil {
		return Player{}, err
	}
	if len(users) > 0 {
		return player, errors.New("player with this email is already created")
	}

	player.Id = CreateId()
	return player, s.table.Put(player).Run()
}

func (s Service) AddMatch(user Player, matchId string) (Player, error) {
	user.Matches = append(user.Matches, matchId)
	return user, s.table.Update("id", user.Id).Set("matches", user.Matches).Value(&user)
}

func (s Service) GetAll() (users []Player, err error) {
	err = s.table.Scan().All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s Service) GetByEmail(email string) (user []Player, err error) {
	err = s.table.Get("email", email).Index("email-index").All(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func CreateId() (id string) {
	v4, _ := uuid.NewV4()
	return v4.String()
}
