package constants

var (
	Bearer = "ut"

	UserStatusComplete   = "complete"
	UserStatusIncomplete = "incomplete"

	RoleManager  = "manager"
	RoleOperator = "operator"
	RoleDelivery = "delivery"
	RoleCustomer = "customer"

	PermissionManageUser      = "manage_user"
	PermissionManageMenu      = "manage_menu"
	PermissionViewOrder       = "view_order"
	PermissionManagerOrder    = "manage_order"
	PermissionAssignOrder     = "assign_orders"
	PermissionMarkAsDelivered = "mark_as_delivered"
)
