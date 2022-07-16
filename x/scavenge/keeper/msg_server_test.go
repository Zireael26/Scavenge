package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/Zireael26/scavenge/testutil/keeper"
	"github.com/Zireael26/scavenge/x/scavenge/keeper"
	"github.com/Zireael26/scavenge/x/scavenge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.ScavengeKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
