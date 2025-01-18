package setup

import (
	"time"

	"github.com/Dev-Miniplays/Tickets-Worker/bot/command"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/command/registry"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/customisation"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/dbclient"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/utils"
	"github.com/Dev-Miniplays/Tickets-Worker/i18n"
	"github.com/TicketsBot/common/permission"
	"github.com/rxdn/gdl/objects/interaction"
	"github.com/rxdn/gdl/rest/request"
)

type TranscriptsSetupCommand struct{}

func (TranscriptsSetupCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:            "transcripts",
		Description:     i18n.HelpSetup,
		Type:            interaction.ApplicationCommandTypeChatInput,
		Aliases:         []string{"transcript", "archives", "archive"},
		PermissionLevel: permission.Admin,
		Category:        command.Settings,
		Arguments: command.Arguments(
			command.NewRequiredArgument("channel", "The channel that ticket transcripts should be sent to", interaction.OptionTypeChannel, i18n.SetupTranscriptsInvalid),
		),
		Timeout: time.Second * 5,
	}
}

func (c TranscriptsSetupCommand) GetExecutor() interface{} {
	return c.Execute
}

func (TranscriptsSetupCommand) Execute(ctx registry.CommandContext, channelId uint64) {
	if _, err := ctx.Worker().GetChannel(channelId); err != nil {
		if restError, ok := err.(request.RestError); ok && restError.IsClientError() {
			ctx.Reply(customisation.Red, i18n.Error, i18n.SetupTranscriptsInvalid, ctx.ChannelId)
		} else {
			ctx.HandleError(err)
		}

		return
	}

	if err := dbclient.Client.ArchiveChannel.Set(ctx, ctx.GuildId(), utils.Ptr(channelId)); err == nil {
		ctx.Reply(customisation.Green, i18n.TitleSetup, i18n.SetupTranscriptsComplete, channelId)
	} else {
		ctx.HandleError(err)
	}
}
