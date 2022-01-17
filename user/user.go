package user

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	uuid "github.com/nu7hatch/gouuid"
)

type User struct {
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	NickName string `json:"nick_name"`
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

func (s Service) Create(user User) (User, error) {
	user.Id = CreateId()
	return user, s.table.Put(user).Run()
}

func CreateId() (id string) {
	v4, _ := uuid.NewV4()
	return v4.String()
}
