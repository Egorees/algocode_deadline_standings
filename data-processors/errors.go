package data_processors

import "fmt"

type DeadlineResultsError struct {
	Reason string
}

type DataError struct {
	Reason string
}

func (e *DeadlineResultsError) Error() string {
	return fmt.Sprintf("Deadline result error: %s", e.Reason)
}

func (e *DataError) Error() string {
	return fmt.Sprintf("Data error: %s", e.Reason)
}
