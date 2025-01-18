package listeners

import (
	"context"
	"time"

	"github.com/Dev-Miniplays/Tickets-Worker"
	cmdcontext "github.com/Dev-Miniplays/Tickets-Worker/bot/command/context"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/constants"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/dbclient"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/listeners/messagequeue"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/logic"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/utils"
	"github.com/TicketsBot/common/sentry"
	"github.com/rxdn/gdl/gateway/payloads/events"
	gdlUtils "github.com/rxdn/gdl/utils"
)

// Remove user permissions when they leave
func OnMemberLeave(worker *worker.Context, e events.GuildMemberRemove) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3) // TODO: Propagate context
	defer cancel()

	if err := dbclient.Client.Permissions.RemoveSupport(ctx, e.GuildId, e.User.Id); err != nil {
		sentry.Error(err)
	}

	if err := utils.ToRetriever(worker).Cache().DeleteCachedPermissionLevel(ctx, e.GuildId, e.User.Id); err != nil {
		sentry.Error(err)
	}

	// auto close
	settings, err := dbclient.Client.AutoClose.Get(ctx, e.GuildId)
	if err != nil {
		sentry.Error(err)
	} else {
		// check setting is enabled
		if settings.Enabled && settings.OnUserLeave != nil && *settings.OnUserLeave {
			// get open tickets by user
			tickets, err := dbclient.Client.Tickets.GetOpenByUser(ctx, e.GuildId, e.User.Id)
			if err != nil {
				sentry.Error(err)
			} else {
				for _, ticket := range tickets {
					isExcluded, err := dbclient.Client.AutoCloseExclude.IsExcluded(ctx, e.GuildId, ticket.Id)
					if err != nil {
						sentry.Error(err)
						continue
					}

					if isExcluded {
						continue
					}

					// verify ticket exists + prevent potential panic
					if ticket.ChannelId == nil {
						return
					}

					// get premium status
					premiumTier, err := utils.PremiumClient.GetTierByGuildId(ctx, ticket.GuildId, true, worker.Token, worker.RateLimiter)
					if err != nil {
						sentry.Error(err)
						return
					}

					ctx, cancel := context.WithTimeout(context.Background(), constants.TimeoutCloseTicket)

					cc := cmdcontext.NewAutoCloseContext(ctx, worker, e.GuildId, *ticket.ChannelId, worker.BotId, premiumTier)
					logic.CloseTicket(ctx, cc, gdlUtils.StrPtr(messagequeue.AutoCloseReason), true)

					cancel()
				}
			}
		}
	}
}
