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
	pi := core.OpenPrivateInfo(path.Join(os.Getenv("HOME"), ".config", ".p3pgroup"), "Group Host", "")
	pi.Endpoint = core.Endpoint(os.Getenv("PRIVATEINFO_ROOT_ENDPOINT"))
	pi.MessageCallback = append(pi.MessageCallback, botMsgHandler)
	if !pi.IsAccountReady() {
		pi.Create("Group Host", "no@no.no", 4096)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	log.Println("p3p is running...")
	log.Println("Group Server is available at:", pi.Endpoint)

	for sig := range c {
		log.Println("Closing [", sig, "] ...")
		return
	}
}
