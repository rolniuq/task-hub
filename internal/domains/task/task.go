package task

import "taskhub/pkg/base/entity"

type Task struct {
	entity.BaseEntity
	Title   string
	Content string
}
