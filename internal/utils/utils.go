package utils

import (
	"fmt"
	"github.com/sho0pi/tickli/internal/types"
	"github.com/sho0pi/tickli/internal/types/project"
)

func GetProjectDescription(project types.Project) string {
	var projectStatus string
	if project.Closed {
		projectStatus = "Closed"
	} else {
		projectStatus = "Open"
	}

	projectLine := project.Color.Sprint("■■■■■■■■■■■■■■■■■■■■■■■■", project.Color)

	description := fmt.Sprintf(`
Project Details:

%s
Name: %s
ID: %s
Type: %s 
Status: %s
Group: %s

Tasks:`,
		projectLine,
		project.Name,
		project.ID,
		project.Kind,
		projectStatus,
		project.GroupID,
	)

	return description
}

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
