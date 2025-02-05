package app

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CompleteTask     struct{}
	ChangeTaskStatus struct{}
	AssignTask       struct{}
	UnassignTask     struct{}
}

type Queries struct {
	AllTasks struct{}
}
