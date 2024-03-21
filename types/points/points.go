package points

type Points struct {
	Id      string
	UserId  string
	Balance int32
}

type GetPointsRequestBody struct{}

type GetPointsByUserRequestBody struct {
	UserId string
}

type UpdatePointsRequestBody struct {
	UserId     string
	NewBalance int32
}

type UpdatePointsResponseBody struct {
	Status bool
}
