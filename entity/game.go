package entity

type Game struct {
	ID          uint
	CategoryID  uint
	QuestionIDs []uint
	PlayerIDs   []uint
}
