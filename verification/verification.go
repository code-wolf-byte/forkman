package verification

import (
	"fmt"

	"github.com/avvo-na/devil-guard/config"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// This is the list of commands that the bot will register
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "echo",
		Description: "Ping the bot",
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"echo": echo,
}

func EnableModule(s *discordgo.Session) error {
	// Write the module config
	err := config.WriteEnableModule("verification")
	if err != nil {
		return fmt.Errorf("failed to write module config: %w", err)
	}

	// Register the commands
	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(
			config.AppCfg.DiscordAppID,
			config.AppCfg.DiscordDevGuildID,
			v,
		)
		// TODO: add better handling if one command fails to register
		if err != nil {
			log.Error().Err(err).Msg("Failed to register command")
		}
	}

	// This is a map of command names to their handlers. When a command is
	// received, the bot will check if the command name is in this map. If it
	// is, the bot will call the handler function with the session and the
	// interaction.
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	log.Info().Msg("Enabled verification module")
	return nil
}

// TODO: will be called from a rest API
func DisableModule(s *discordgo.Session) error {
	// Write the module config
	err := config.WriteDisableModule("verification")
	if err != nil {
		return fmt.Errorf("failed to write module config: %w", err)
	}

	// NOTE:
	// Grab all regeistered and cross check with the commands
	// n^2 so its a bit slow :() idk if it matters
	registeredCommands, err := s.ApplicationCommands(config.AppCfg.DiscordAppID, config.AppCfg.DiscordDevGuildID)
	for _, v := range registeredCommands {
		for _, c := range commands {
			if c.Name == v.Name {
				err := s.ApplicationCommandDelete(
					config.AppCfg.DiscordAppID,
					config.AppCfg.DiscordDevGuildID,
					v.ID,
				)
				if err != nil {
					log.Error().Err(err).Msg("Failed to delete command")
				}
			}
		}
	}

	log.Info().Msg("Disabled verification module")
	return nil
}

func echo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
	log.Info().Interface("command", i.ApplicationCommandData()).Msg("Responded to interaction request")
}
