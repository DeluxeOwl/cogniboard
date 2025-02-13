package app

import (
	"errors"
	"log/slog"

	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/DeluxeOwl/cogniboard/internal/project/app/commands"
	"github.com/DeluxeOwl/cogniboard/internal/project/app/operations"
	"github.com/DeluxeOwl/cogniboard/internal/project/app/queries"
)

type Application struct {
	Commands   Commands
	Queries    Queries
	Operations Operations
}

type Commands struct {
	CreateTask        commands.CreateTaskHandler
	ChangeTaskStatus  commands.ChangeTaskStatusHandler
	EditTask          commands.EditTaskHandler
	AttachFilesToTask commands.AttachFilesToTaskHandler
}

type Queries struct {
	AllTasks queries.AllTasksHandler
}

type Operations struct {
	ChatWithProject operations.ChatWithProjectHandler
}

// New creates a new Application instance with the provided dependencies
func New(
	repo project.TaskRepository,
	logger *slog.Logger,
	fileStorage project.FileStorage,
	chatService operations.ChatService,
	embeddings project.EmbeddingStorage,
) (*Application, error) {
	if repo == nil {
		return nil, errors.New("repo cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if fileStorage == nil {
		return nil, errors.New("file storage cannot be nil")
	}
	if chatService == nil {
		return nil, errors.New("chat service cannot be nil")
	}
	if embeddings == nil {
		return nil, errors.New("embeddings cannot be nil")
	}

	return &Application{
		Commands: Commands{
			CreateTask:       commands.NewCreateTaskHandler(repo, logger),
			ChangeTaskStatus: commands.NewChangeStatusHandler(repo, logger),
			EditTask:         commands.NewEditTaskHandler(repo, logger),
			AttachFilesToTask: commands.NewAttachFilesToTaskHandler(
				repo,
				logger,
				fileStorage,
				embeddings,
			),
		},
		Queries: Queries{
			AllTasks: queries.NewAllTasksHandler(repo, logger),
		},
		Operations: Operations{
			ChatWithProject: operations.NewChatWithProjectHandler(chatService, repo, logger),
		},
	}, nil
}
