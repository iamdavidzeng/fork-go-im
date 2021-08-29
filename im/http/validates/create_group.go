package validates

type CreateGroupParams struct {
	UserId    map[string]string `valid:"user_id"`
	GroupName string            `valid:"group_name"`
}
