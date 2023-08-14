package db

var ConfigBountyRoles []BountyRoles = []BountyRoles{
	{
		Name: "ADD BOUNTY",
	},
	{
		Name: "UPDATE BOUNTY",
	},
	{
		Name: "DELETE BOUNTY",
	},
	{
		Name: "PAY BOUNTY",
	},
	{
		Name: "ADD USER",
	},
	{
		Name: "UPDATE USER",
	},
	{
		Name: "DELETE USER",
	},
	{
		Name: "ADD BUDGET",
	},
	{
		Name: "WITHDRAW BUDGET",
	},
	{
		Name: "VIEW REPORT",
	},
}

var Updatables = []string{
	"name", "description", "tags", "img",
	"owner_alias", "price_to_join", "price_per_message",
	"escrow_amount", "escrow_millis",
	"unlisted", "private", "deleted",
	"app_url", "bots", "feed_url", "feed_type",
	"owner_route_hint", "updated", "pin",
	"profile_filters",
}
var Botupdatables = []string{
	"name", "description", "tags", "img",
	"owner_alias", "price_per_use",
	"unlisted", "deleted",
	"owner_route_hint", "updated",
}
var Peopleupdatables = []string{
	"description", "tags", "img",
	"owner_alias",
	"unlisted", "deleted",
	"owner_route_hint",
	"price_to_meet", "updated",
	"extras",
}
var Channelupdatables = []string{
	"name", "deleted"}
