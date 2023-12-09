package main

import (
	"log"

	"git.mrcyjanek.net/p3pch4t/p3pgo/lib/core"
	"gorm.io/gorm"
)

// P3PGROUP_UserInfo This one is for bot.
type P3PGROUP_UserInfo struct {
	gorm.Model
	KeyID        string
	isIntroduced bool
}

// P3PGROUP_MemberUserInfo This one is for groups
type P3PGROUP_MemberUserInfo struct {
	gorm.Model
	KeyID        string
	IsUserBanned bool
	BanReason    string
	MessagesSent int
}

func dbAutoMigrateBot(pi *core.PrivateInfoS) {
	log.Println("dbAutoMigrateBot P3PGROUP_UserInfo", pi.DB.AutoMigrate(&P3PGROUP_UserInfo{}))
	log.Println("dbAutoMigrateBot P3PGROUP_GroupInfo", pi.DB.AutoMigrate(&P3PGROUP_GroupInfo{}))
}

func dbAutoMigrateGroup(pi *core.PrivateInfoS) {
	log.Println("dbAutoMigrateGroup P3PGROUP_MemberUserInfo", pi.DB.AutoMigrate(&P3PGROUP_MemberUserInfo{}))
}
