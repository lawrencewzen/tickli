package utils

import (
	"fmt"
	"github.com/sho0pi/tickli/internal/types"
	"github.com/sho0pi/tickli/internal/types/project"
)

func GetTaskDescription(task types.Task, projectColor project.Color) string {
	projectLine := projectColor.Sprint("----------------------")

	description := fmt.Sprintf(`
Task Details:

%s %s
%s
Desc: %s 
Content: %s
Priority: %s
Group: %s

Time: 
StartDate: %s
DueDate: %s
CompletedTime: %s

Tasks:`,
		task.Status,
		projectLine,
		task.Title,
		task.Desc,
		task.Content,
		task.Priority.String(),
		task.ProjectID,
		task.StartDate.Humanize(),
		task.DueDate,
		task.CompletedTime.String(),
	)

	return description
}
