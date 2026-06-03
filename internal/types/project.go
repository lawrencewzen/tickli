package types

import "github.com/sho0pi/tickli/internal/types/project"

// InboxProject the Inbox project representation (cause is not returned by the api)
var InboxProject = Project{
	ID:        "inbox",
	Name:      "📥Inbox",
	Color:     project.DefaultColor,
	SortOrder: 0,
	Closed:    false,
	Kind:      project.KindTask,
	ViewMode:  project.ViewModeList,
}

type Project struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Color      project.Color    `json:"color"`
	SortOrder  int64            `json:"sortOrder"`
	Closed     bool             `json:"closed"`
	GroupID    string           `json:"groupId"`
	ViewMode   project.ViewMode `json:"viewMode"`
	Permission string           `json:"permission"`
	Kind       project.Kind     `json:"kind"`
}
