package http

import (
	"net/http"

	"gitverse-internship-zg/internal"
	"gitverse-internship-zg/services/user-service/usecase"
)

func SetupRoutes(api *usecase.Usecase, mux *http.ServeMux) {
	// Создаем цепочку миддлварей и передаем API через замыкание
	// нужно создать композитный хендлер-миддлварь в котором я буду использовать свитч по методу у ручки чтобы
	// выбирать обработку по нужной миддлеварине
	// т.е нужна миддлеварь распределитель
	Usersget := internal.ChainMiddleware(
		UsersHandler(api),
		// internal.GetMethodMiddleware,
	)

	mux.Handle("/api/buy", Usersget) // get
	mux.Handle("/api/auth")          // post
	mux.Handle("/api/sendCoin")      // post
	mux.Handle("/api/info")          // get
}
