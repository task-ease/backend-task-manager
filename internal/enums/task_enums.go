package enums

type TaskTypes string

const (
	TaskTypeEpic    TaskTypes = "EPIC"
	TaskTypeTask    TaskTypes = "TASK"
	TaskTypeSubtask TaskTypes = "SUBTASK"
	TaskTypeBug     TaskTypes = "BUG"
	TaskTypeFeature TaskTypes = "FEATURE"
	TaskTypeChore   TaskTypes = "CHORE"
	TaskTypeSpike   TaskTypes = "SPIKE"
)

type TaskPriorities string

const (
	TaskPrioritiesVeryLow  TaskPriorities = "VERY_LOW"
	TaskPrioritiesLow      TaskPriorities = "LOW"
	TaskPrioritiesMid      TaskPriorities = "MID"
	TaskPrioritiesHigh     TaskPriorities = "HIGH"
	TaskPrioritiesVeryHigh TaskPriorities = "VERY_HIGH"
)
