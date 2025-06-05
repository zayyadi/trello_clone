package main

import (
	"log"

	"github.com/zayyadi/trello/config"
	"github.com/zayyadi/trello/db"
	"github.com/zayyadi/trello/handlers"
	middleware "github.com/zayyadi/trello/middlewares"
	"github.com/zayyadi/trello/repositories"
	"github.com/zayyadi/trello/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	dbInstance, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Repositories
	userRepo := repositories.NewUserRepository(dbInstance)
	boardRepo := repositories.NewBoardRepository(dbInstance)
	listRepo := repositories.NewListRepository(dbInstance)
	cardRepo := repositories.NewCardRepository(dbInstance)
	boardMemberRepo := repositories.NewBoardMemberRepository(dbInstance)
	commentRepo := repositories.NewCommentRepository(dbInstance) // Initialize CommentRepository

	// Initialize Services
	authService := services.NewAuthService(userRepo, cfg.JWTSecretKey)
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	listService := services.NewListService(listRepo, boardRepo, boardMemberRepo)
	cardService := services.NewCardService(cardRepo, listRepo, boardRepo, boardMemberRepo, userRepo)          // Added userRepo
	commentService := services.NewCommentService(commentRepo, cardRepo, listRepo, boardRepo, boardMemberRepo) // Initialize CommentService

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	boardHandler := handlers.NewBoardHandler(boardService)
	listHandler := handlers.NewListHandler(listService)
	cardHandler := handlers.NewCardHandler(cardService)
	commentHandler := handlers.NewCommentHandler(commentService) // Initialize CommentHandler

	// Setup Gin router
	// gin.SetMode(gin.ReleaseMode) // Uncomment for production
	router := gin.Default()

	// CORS Middleware (Allow All for development)
	// For production, configure origins properly
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecretKey))
	{
		// Board routes
		api.POST("/boards", boardHandler.CreateBoard)
		api.GET("/boards", boardHandler.GetBoardsForUser)
		api.GET("/boards/:boardID", boardHandler.GetBoardByID)
		api.PUT("/boards/:boardID", boardHandler.UpdateBoard)
		api.DELETE("/boards/:boardID", boardHandler.DeleteBoard)

		// Board Member routes
		api.POST("/boards/:boardID/members", boardHandler.AddMemberToBoard)
		api.GET("/boards/:boardID/members", boardHandler.GetBoardMembers)
		api.DELETE("/boards/:boardID/members/:memberUserID", boardHandler.RemoveMemberFromBoard)

		// List routes
		api.POST("/boards/:boardID/lists", listHandler.CreateList)
		api.GET("/boards/:boardID/lists", listHandler.GetListsByBoardID)
		api.PUT("/lists/:listID", listHandler.UpdateList)
		api.DELETE("/lists/:listID", listHandler.DeleteList)
		// Consider a route for reordering lists: PATCH /api/boards/:boardID/lists/reorder

		// Card routes
		api.POST("/lists/:listID/cards", cardHandler.CreateCard)
		api.GET("/lists/:listID/cards", cardHandler.GetCardsByListID)
		api.GET("/cards/:cardID", cardHandler.GetCardByID)
		api.PUT("/cards/:cardID", cardHandler.UpdateCard)
		api.DELETE("/cards/:cardID", cardHandler.DeleteCard)
		api.PATCH("/cards/:cardID/move", cardHandler.MoveCard)
		// Consider a route for reordering cards within a list: PATCH /api/lists/:listID/cards/reorder

		// Comment routes
		api.POST("/cards/:cardID/comments", commentHandler.CreateComment)
		api.GET("/cards/:cardID/comments", commentHandler.GetCommentsByCardID)

		// Card Collaborator routes
		api.POST("/cards/:cardID/collaborators", cardHandler.AddCollaborator)
		api.GET("/cards/:cardID/collaborators", cardHandler.GetCollaborators)
		api.DELETE("/cards/:cardID/collaborators/:userID", cardHandler.RemoveCollaborator)
	}

	// Start server
	port := cfg.ServerPort
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
