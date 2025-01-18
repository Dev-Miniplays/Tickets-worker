package settings

import (
	"time"

	"github.com/Dev-Miniplays/Tickets-Worker/bot/command"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/command/registry"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/customisation"
	"github.com/Dev-Miniplays/Tickets-Worker/i18n"
	"github.com/TicketsBot/common/permission"
	"github.com/rxdn/gdl/objects/interaction"
)

type PanelCommand struct {
}

func (PanelCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:             "panel",
		Description:      i18n.HelpPanel,
		Type:             interaction.ApplicationCommandTypeChatInput,
		PermissionLevel:  permission.Admin,
		Category:         command.Settings,
		DefaultEphemeral: true,
		Timeout:          time.Second * 3,
	}
}

func (c PanelCommand) GetExecutor() interface{} {
	return c.Execute
}

func (PanelCommand) Execute(ctx registry.CommandContext) {
	ctx.Reply(customisation.Green, i18n.TitlePanel, i18n.MessagePanel, ctx.GuildId())
}
