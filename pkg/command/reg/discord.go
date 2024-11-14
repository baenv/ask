package reg

import (
	"fmt"
	"net/url"
	"strings"
	"sum/pkg/logger"
	"sum/pkg/models"
	"sum/pkg/repo"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	repo   repo.Repository
	logger logger.Logger
}

func NewDiscord(repo repo.Repository, logger logger.Logger) *Discord {
	return &Discord{
		repo:   repo,
		logger: logger,
	}
}

func (d *Discord) Info() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "reg",
		Description: "Register a new user or server agent",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "server",
				Description: "Register agent for this server",
				Required:    false,
			},
		},
	}
}

func (d *Discord) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if the user is not the server owner
	if i.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You must be a server administrator to use this command.",
			},
		})
		return
	}

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if i.ApplicationCommandData().Name != "reg" {
		return
	}

	regType := "user"
	if len(i.ApplicationCommandData().Options) > 0 &&
		i.ApplicationCommandData().Options[0].Name == "server" {
		regType = "server"
	}

	// Create and show modal
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: fmt.Sprintf("reg_%s_%v", i.Member.User.ID, regType),
			Title:    "Agent Registration",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "agent_url",
							Label:       "Agent URL",
							Style:       discordgo.TextInputShort,
							Placeholder: "https://example.com/api",
							Required:    false,
							MaxLength:   200,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "api_token",
							Label:       "API Token",
							Style:       discordgo.TextInputShort,
							Placeholder: "Your API token",
							Required:    false,
							MaxLength:   100,
						},
					},
				},
			},
		},
	})

	if err != nil {
		d.logger.Error(err, "Failed to show registration modal")
		return
	}
}

func (d *Discord) HandleSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionModalSubmit {
		return
	}

	submitIDParts := strings.Split(i.ModalSubmitData().CustomID, "_")
	if len(submitIDParts) != 3 {
		return
	}

	if submitIDParts[0] != "reg" {
		return
	}

	if submitIDParts[1] != i.Member.User.ID {
		return
	}

	var agentURL, apiToken string
	if len(i.ModalSubmitData().Components) > 0 {
		if row, ok := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow); ok && len(row.Components) > 0 {
			if input, ok := row.Components[0].(*discordgo.TextInput); ok {
				agentURL = input.Value
			}
		}
	}
	if len(i.ModalSubmitData().Components) > 1 {
		if row, ok := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow); ok && len(row.Components) > 0 {
			if input, ok := row.Components[0].(*discordgo.TextInput); ok {
				apiToken = input.Value
			}
		}
	}

	// Validate input
	if agentURL != "" {
		if _, err := url.ParseRequestURI(agentURL); err != nil {
			d.respondWithError(s, i, "Invalid Agent URL format")
			return
		}
	}

	isServer := submitIDParts[2] == "server"
	d.handleRegistration(s, i, agentURL, apiToken, isServer)
}

func (d *Discord) handleRegistration(s *discordgo.Session, i *discordgo.InteractionCreate, agentURL, apiToken string, isServer bool) {
	user := models.User{
		UserID:   i.Member.User.ID,
		Username: i.Member.User.Username,
		Platform: models.PlatformDiscord,
	}

	thisGuild, err := s.State.Guild(i.GuildID)
	if err != nil {
		d.logger.Error(err, "Failed to get guild state")
		return
	}

	if isServer {
		server := models.Server{
			ServerID:   i.GuildID,
			ServerName: thisGuild.Name,
			Platform:   models.PlatformDiscord,
			OwnerID:    i.Member.User.ID,
		}

		if agentURL != "" && apiToken != "" {
			server.ServerAdminConfig = []models.ServerAdminConfig{
				{
					APIKey:      apiToken,
					EndpointURL: agentURL,
				},
			}
		}

		user.Servers = []models.Server{server}
	} else if agentURL != "" && apiToken != "" {
		user.UserAgentConfigs = []models.UserAgentConfig{
			{
				APIKey:      apiToken,
				EndpointURL: agentURL,
			},
		}
	}

	_, err = d.repo.User().Create(user)
	if err != nil {
		d.logger.Error(err, "Failed to create or update user/server with config")
		d.respondWithError(s, i, fmt.Sprintf("Failed to register %s. Please try again.", map[bool]string{true: "server", false: "user"}[isServer]))
		return
	}

	d.respondWithSuccess(s, i, fmt.Sprintf("%s registered successfully!", map[bool]string{true: "Server", false: "User"}[isServer]))
}

func (d *Discord) respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Error: %s", message),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func (d *Discord) respondWithSuccess(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
