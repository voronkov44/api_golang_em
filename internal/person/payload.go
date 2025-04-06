package person

type PersonCreateRequest struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"` // поправил опечатку: "patronomic" → "patronymic"
}

type PersonUpdateRequest struct {
	Name       *string `json:"name,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Patronymic *string `json:"patronymic,omitempty"`
}

type PersonResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	//Age         int    `json:"age"`
	//Gender      string `json:"gender"`
	//Nationality string `json:"nationality"`
}
