// Package notification содержит типы уведомлений из экспорта чата Telegram
// с ботом @MajesticRolePlayBot.
package notification

import "time"

// NotificationType — строковый тип уведомления.
type NotificationType string

const (
	TypeItemSold             NotificationType = "item_sold"
	TypeItemInWarehouse      NotificationType = "item_in_warehouse"
	TypePunishment           NotificationType = "punishment"
	TypeSubscriptionExpiring NotificationType = "subscription_expiring"
	TypeSubscriptionWarning  NotificationType = "subscription_warning"
	TypeSubscriptionFrozen   NotificationType = "subscription_frozen"
	TypeOrgAttack            NotificationType = "org_attack"    // наша орга напала
	TypeOrgDefend            NotificationType = "org_defend"    // на нашу напали
	TypePropertySold         NotificationType = "property_sold"
	TypeVehicleRented        NotificationType = "vehicle_rented"
)

// BaseNotification — общие поля для всех уведомлений.
type BaseNotification struct {
	Type          NotificationType `json:"type"`
	Server        string           `json:"server"`
	CharID        string           `json:"char_id,omitempty"`
	Timestamp     time.Time        `json:"timestamp"`
	HTMLMessageID int              `json:"html_message_id"`
}

// ItemSoldNotification — продажа предмета на маркете.
type ItemSoldNotification struct {
	BaseNotification
	Character string  `json:"character"`
	ItemName  string  `json:"item_name"`
	Quantity  int     `json:"quantity"`
	SalePrice float64 `json:"sale_price"`
	Buyer     string  `json:"buyer"`
}

// ItemInWarehouseNotification — предмет доставлен на склад.
type ItemInWarehouseNotification struct {
	BaseNotification
	Character           string `json:"character"`
	ItemName            string `json:"item_name"`
	Quantity            int    `json:"quantity"`
	PickupDeadlineHours int    `json:"pickup_deadline_hours"`
}

// PunishmentNotification — наказание персонажа.
type PunishmentNotification struct {
	BaseNotification
	Character      string  `json:"character"`
	PunishmentType string  `json:"punishment_type"`
	Admin          string  `json:"admin"`
	Reason         string  `json:"reason"`
	DurationHours  float64 `json:"duration_hours"`
}

// SubscriptionExpiringNotification — подписка заканчивается через N дней.
type SubscriptionExpiringNotification struct {
	BaseNotification
	Login    string `json:"login"`
	DaysLeft int    `json:"days_left"`
}

// SubscriptionWarningNotification — подписка скоро заканчивается.
type SubscriptionWarningNotification struct {
	BaseNotification
}

// SubscriptionFrozenNotification — подписка заморожена.
type SubscriptionFrozenNotification struct {
	BaseNotification
}

// OrgAttackNotification — ваша организация напала на чужую.
type OrgAttackNotification struct {
	BaseNotification
	Character     string    `json:"character"`
	OrgName       string    `json:"org_name"`
	EnemyName     string    `json:"enemy_name"`
	AttackStart   time.Time `json:"attack_start"`
	SquareName    string    `json:"square_name"`
	SquareNumber  string    `json:"square_number"`
	AttackerCount int       `json:"attacker_count"`
	WeaponCaliber string    `json:"weapon_caliber"`
}

// OrgDefendNotification — на вашу организацию напали.
type OrgDefendNotification struct {
	BaseNotification
	Character     string    `json:"character"`
	OrgName       string    `json:"org_name"`
	AttackerName  string    `json:"attacker_name"`
	AttackStart   time.Time `json:"attack_start"`
	SquareName    string    `json:"square_name"`
	SquareNumber  string    `json:"square_number"`
	AttackerCount int       `json:"attacker_count"`
	WeaponCaliber string    `json:"weapon_caliber"`
}

// VehicleRentedNotification — транспорт сдан в аренду.
type VehicleRentedNotification struct {
	BaseNotification
	Character     string  `json:"character"`
	VehicleName   string  `json:"vehicle_name"`
	VehicleNumber string  `json:"vehicle_number"`
	Price         float64 `json:"price"`
	DurationHours float64 `json:"duration_hours"`
	Renter        string  `json:"renter"`
}

// PropertySoldNotification — продажа имущества.
type PropertySoldNotification struct {
	BaseNotification
	Character string  `json:"character"`
	Name      string  `json:"name"`
	SalePrice float64 `json:"sale_price"`
	Buyer     string  `json:"buyer"`
}
