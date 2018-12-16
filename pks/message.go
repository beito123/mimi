package pks

// Client -> Server
type StartProgram struct {
	ProgramName string `json:"program"`
}

// Client -> Server
type StopProgram struct {
	ProgramName string `json:"program"`
}

// Client -> Server
type RestartProgram struct {
	ProgramName string `json:"program"`
}

// Server -> Client
type EndProgram struct {
	ProgramName string `json:"program"`
}

// Server -> Client
type RequestConsoleList struct {
}

// Server -> Client
type ResponseConsoleList struct {
}

// Server -> Client
type JoinConsole struct {
}

// Server -> Client
type QuitConsole struct {
}

// Server -> Client
type ConsoleMessage struct {
}

// Server -> Client
type SendCommand struct {
}
