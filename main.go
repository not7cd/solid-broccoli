package main

import (
	"bufio"
	"database/sql"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"solidbroccoli/config"
	"solidbroccoli/graph"
	"strings"
	"syscall"
)

var token string

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", cfg.SqlitePath)
	if err != nil {
		log.Fatal(err)
	}

	graph.Repository = graph.NewSQLiteRepository(db)
	if err := graph.Repository.Migrate(); err != nil {
		log.Fatal(err)
	}

	discord(cfg)
}

func discord(cfg config.Config) {
	token = cfg.DiscordToken

	if token == "" {
		log.Println("No discord token provided")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	//// Register guildCreate as a callback for the guildCreate events.
	//dg.AddHandler(guildCreate)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		log.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Solid-broccoli is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func cli() {
	in := bufio.NewReader(os.Stdin)
	for {
		line, err := in.ReadString('\n')
		if err != nil {
			log.Println("Bad input")
		}
		// sanitize it a bit
		line = strings.TrimSpace(line)
		if line[0] == config.PrefixCmd {
			log.Printf("Got command\n")
			graph.ProcessQuery(line[1:])
		}
	}
}
