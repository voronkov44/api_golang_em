package person

type Person struct {
	ID          int64 `gorm:"primary_key"`
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}

func (Person) TableName() string {
	return "persons"
}
