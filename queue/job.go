package queue

import (
	"github.com/RichardKnop/machinery/v2/tasks"
)

type Job struct {
	Name string
	Args []Arg
	// should be a function
	Run       interface{}
	Retries   int
	OnSuccess *Job
	OnError   *Job
}

type Arg struct {
	Name  string
	Type  string
	Value interface{}
}

func (arg Arg) ToTaskArg() tasks.Arg {
	return tasks.Arg{
		Name:  arg.Name,
		Type:  arg.Type,
		Value: arg.Value,
	}
}

func NewSendEmailJob(to string, content string, textContent string) Job {
	return Job{
		Name: "send_email",
		Args: []Arg{
			Arg{
				Name:  "email",
				Type:  "string",
				Value: to,
			},
			Arg{
				Name:  "content",
				Type:  "string",
				Value: content,
			},
			Arg{
				Name:  "textContent",
				Type:  "string",
				Value: textContent,
			},
		},
		Retries: 3,
	}
}
