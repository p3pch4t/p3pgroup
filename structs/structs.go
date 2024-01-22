package structs

import (
	"git.mrcyjanek.net/p3pch4t/p3pgo/lib/core"
	"gorm.io/gorm"
	"log"
)

// P3PGROUP_UserInfo This one is for bot.
type P3PGROUP_UserInfo struct {
	gorm.Model
	KeyID        string
	IsIntroduced bool
}

// P3PGROUP_MemberUserInfo This one is for groups
type P3PGROUP_MemberUserInfo struct {
	gorm.Model
	KeyID          string
	IsUserBanned   bool
	IsUserWelcomed bool
	IsUserAdmin    bool
	BanReason      string
	MessagesSent   int
}
type P3PGROUP_GroupInfo struct {
	gorm.Model
	GroupName    string
	GroupEmail   string
	EndpointPath string
	UniqueKey    string
	OwnerKeyID   string
	PI           *core.PrivateInfoS `gorm:"-"`
}

func DbAutoMigrateBot(pi *core.PrivateInfoS) {
	log.Println("dbAutoMigrateBot P3PGROUP_UserInfo", pi.DB.AutoMigrate(&P3PGROUP_UserInfo{}))
	log.Println("dbAutoMigrateBot P3PGROUP_GroupInfo", pi.DB.AutoMigrate(&P3PGROUP_GroupInfo{}))
}

func DbAutoMigrateGroup(pi *core.PrivateInfoS) {
	log.Println("dbAutoMigrateGroup P3PGROUP_MemberUserInfo", pi.DB.AutoMigrate(&P3PGROUP_MemberUserInfo{}))
}
