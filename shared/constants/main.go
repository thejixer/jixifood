package constants

import "github.com/lib/pq"

var (
	Bearer = "ut"

	UserStatusComplete   = "complete"
	UserStatusIncomplete = "incomplete"

	RoleManager  = "manager"
	RoleOperator = "operator"
	RoleDelivery = "delivery"
	RoleCustomer = "customer"

	PermissionManageUser      = "manage_user"
	PermissionViewUser        = "view_user"
	PermissionManageMenu      = "manage_menu"
	PermissionViewOrder       = "view_order"
	PermissionManagerOrder    = "manage_order"
	PermissionAssignOrder     = "assign_orders"
	PermissionMarkAsDelivered = "mark_as_delivered"

	PGForeignKeyViolationCode = pq.ErrorCode("23503")
	PGDuplicateKeyErrorCode   = pq.ErrorCode("23505")
)
