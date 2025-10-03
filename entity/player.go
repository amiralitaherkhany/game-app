package entity

type Player struct {
	ID      uint
	UserID  uint
	GameID  uint
	score   uint
	Answers []PlayerAnswer
}

type PlayerAnswer struct {
	ID         uint
	PlayerID   uint
	QuestionID uint
	Choice     PossibleAnswerChoice
}
