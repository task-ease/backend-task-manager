package enums

type UpdateColumnTemplateValue string

const (
	ChangeColumnTemplateName  UpdateColumnTemplateValue = "CHANGE_NAME"
	ChangeColumnTemplateColor UpdateColumnTemplateValue = "CHANGE_COLOR"
)

type UpdateColumnTemplateStatus string

const (
	ChangeColumnTemplateRequired UpdateColumnTemplateStatus = "CHANGE_REQUIRE"
	ChangeColumnTemplateActive   UpdateColumnTemplateStatus = "CHANGE_ACTIVE"
	ChangeColumnTemplateDone     UpdateColumnTemplateStatus = "CHANGE_DONE"
)
