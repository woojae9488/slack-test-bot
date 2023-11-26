package slacktestbot

type Server interface {
	SetupMiddlewares()
	SetupRoutes()
	Listen()
}

func StartServer(s Server) {
	s.SetupMiddlewares()
	s.SetupRoutes()
	s.Listen()
}
