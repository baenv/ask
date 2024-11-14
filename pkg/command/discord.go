package command

import (
	"sum/pkg/adapter"
	"sum/pkg/command/reg"
	"sum/pkg/logger"
	"sum/pkg/repo"

	"github.com/bwmarrin/discordgo"
)

// discord represents a Discord command handler
type discord struct {
	session *discordgo.Session
	reg     *reg.Discord
}

// NewDiscord creates a new Discord command handler
func NewDiscord(repo repo.Repository, s *discordgo.Session, a adapter.IAdapter, logger logger.Logger) ICommand {
	return &discord{
		session: s,
		reg:     reg.NewDiscord(repo, logger),
	}
}

// AddHandler adds the sum command handler to the Discord session
func (d *discord) AddHandler() {
	d.session.AddHandler(d.reg.Handle)
	d.session.AddHandler(d.reg.HandleSubmit)
}

// RegisterReg registers the reg command with the Discord API
func (d *discord) RegisterReg() {
	d.session.ApplicationCommandCreate(d.session.State.User.ID, "", d.reg.Info())
}

// RegisterLs registers the ls command with the Discord API
func (d *discord) RegisterLs() {}

// RegisterAi registers the ai command with the Discord API
func (d *discord) RegisterAi() {}

// RegisterStart registers the start command with the Discord API
func (d *discord) RegisterStart() {}

// RegisterSum registers the sum command with the Discord API
func (d *discord) RegisterSum() {}
