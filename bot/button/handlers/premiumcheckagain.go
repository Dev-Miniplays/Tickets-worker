package handlers

import (
	"time"

	"github.com/Dev-Miniplays/Tickets-Worker/bot/button/registry"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/button/registry/matcher"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/command/context"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/customisation"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/dbclient"
	prem "github.com/Dev-Miniplays/Tickets-Worker/bot/premium"
	"github.com/Dev-Miniplays/Tickets-Worker/bot/utils"
	"github.com/Dev-Miniplays/Tickets-Worker/i18n"
	"github.com/TicketsBot/common/permission"
	"github.com/TicketsBot/common/premium"
)

type PremiumCheckAgain struct{}

func (h *PremiumCheckAgain) Matcher() matcher.Matcher {
	return &matcher.SimpleMatcher{
		CustomId: "premium_check_again",
	}
}

func (h *PremiumCheckAgain) Properties() registry.Properties {
	return registry.Properties{
		Flags:   registry.SumFlags(registry.GuildAllowed),
		Timeout: time.Second * 5,
	}
}

func (h *PremiumCheckAgain) Execute(ctx *context.ButtonContext) {
	// Get permission level
	permissionLevel, err := ctx.UserPermissionLevel(ctx)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	if permissionLevel < permission.Admin {
		ctx.Reply(customisation.Red, i18n.Error, i18n.MessageNoPermission)
		return
	}

	ctx.EditWith(customisation.Green, i18n.MessagePremiumChecking, i18n.MessagePremiumPleaseWait)

	if err := utils.PremiumClient.DeleteCachedTier(ctx, ctx.GuildId()); err != nil {
		ctx.HandleError(err)
		return
	}

	if ctx.PremiumTier() > premium.None {
		ctx.EditWith(customisation.Green, i18n.Success, i18n.MessagePremiumSuccessAfterCheck)

		// Re-enable panels
		if err := dbclient.Client.Panel.EnableAll(ctx, ctx.GuildId()); err != nil {
			ctx.HandleError(err)
			return
		}
	} else {
		entitlement, err := dbclient.Client.LegacyPremiumEntitlements.GetUserTier(ctx, ctx.UserId(), premium.PatreonGracePeriod)
		if err != nil {
			ctx.HandleError(err)
			return
		}

		if entitlement == nil {
			ctx.Edit(prem.BuildPatreonNotLinkedMessage(ctx))
		} else {
			res, err := prem.BuildPatreonSubscriptionFoundMessage(ctx, entitlement)
			if err != nil {
				ctx.HandleError(err)
				return
			}

			ctx.Edit(res)
		}
	}
}
