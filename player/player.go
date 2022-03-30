package player

type Player struct {
	Age      int      `dynamo:"age"`
	Email    string   `dynamo:"email"`
	Id       string   `dynamo:"player_id"`
	Name     string   `dynamo:"name"`
	NickName string   `dynamo:"nick_name"`
	Matches  []string `dynamo:"matches"`
}
