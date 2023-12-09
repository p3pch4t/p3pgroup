package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"strconv"

	"git.mrcyjanek.net/p3pch4t/p3pgo/lib/core"
	"github.com/joho/godotenv"
)

var botPi *core.PrivateInfoS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("I2P_HTTP_PROXY") == "" {
		log.Fatalln("I2P_HTTP_PROXY is \"\"")
	}
	if os.Getenv("PRIVATEINFO_ROOT_ENDPOINT") == "" {
		log.Fatalln("PRIVATEINFO_ROOT_ENDPOINT is \"\"")
	}

	core.I2P_HTTP_PROXY = os.Getenv("I2P_HTTP_PROXY")
	core.LOCAL_SERVER_PORT, err = strconv.Atoi(os.Getenv("LOCAL_SERVER_PORT"))
	if err != nil {
		log.Fatalln("LOCAL_SERVER_PORT", err)
	}
	botPi = core.OpenPrivateInfo(path.Join(os.Getenv("HOME"), ".config", ".p3pgroup"), "Group Host", "")
	botPi.Endpoint = core.Endpoint(os.Getenv("PRIVATEINFO_ROOT_ENDPOINT"))
	botPi.MessageCallback = append(botPi.MessageCallback, botMsgHandler)
	botPi.IntroduceCallback = append(botPi.IntroduceCallback, botIntroduceHandler)
	dbAutoMigrateBot(botPi)
	if !botPi.IsAccountReady() {
		botPi.Create("Group Host", "no@no.no", 4096)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	log.Println("p3p is running...")
	loadGroups()
	log.Println("Groups loaded")
	log.Println("Group Server is available at:", botPi.Endpoint)
	for sig := range c {
		log.Println("Closing [", sig, "] ...")
		return
	}
}
