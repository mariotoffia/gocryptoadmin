package common

type ConsoleFormatter interface {
	ConsoleHeader() string
	ConsoleString() string
}
