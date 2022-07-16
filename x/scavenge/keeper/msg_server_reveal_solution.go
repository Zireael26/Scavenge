package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/Zireael26/scavenge/x/scavenge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"
)

func (k msgServer) RevealSolution(goCtx context.Context, msg *types.MsgRevealSolution) (*types.MsgRevealSolutionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// concatenate a solution and a scavenger address and convert it to bytes
	var solutionScavengerBytes = []byte(msg.Solution + msg.Creator)
	// find the hash of solution and address
	var solutionScavengerHash = sha256.Sum256(solutionScavengerBytes)
	// convert hash to a string
	var solutionScavengerHashString = hex.EncodeToString(solutionScavengerHash[:])

	// try getting a commit using the solution and address
	_, isFound := k.GetCommit(ctx, solutionScavengerHashString)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Commit with that hash does not exist")
	}

	// find a hash of the solution
	var solutionHash = sha256.Sum256([]byte(msg.Solution))
	// encode solution hash to string
	var solutionHashString = hex.EncodeToString(solutionHash[:])

	var scavenge types.Scavenge
	// get Scavenge from store using solutionHash
	scavenge, isFound = k.GetScavenge(ctx, solutionHashString)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Scavenge with this solution hash does not exist")
	}

	// check that scavenger property contains a valid address
	_, err := sdk.AccAddressFromBech32(scavenge.Scavenger)
	// return an erro if the scavenge has already been solved
	if err == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Scavenge has already been solved")
	}

	// save the scavenger's address into the scavenge
	scavenge.Scavenger = msg.Creator
	// save the correct solution to the scavenge
	scavenge.Solution = msg.Solution

	// get module account address
	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
	// parse scavenger's address from string to sdk.AccAddress
	scavenger, err := sdk.AccAddressFromBech32(scavenge.Scavenger)
	if err != nil {
		panic(err)
	}

	// parse tokens from string to sdk.Coins
	reward, err := sdk.ParseCoinsNormalized(scavenge.Reward)
	if err != nil {
		panic(err)
	}

	// send tokens from module account to scavenger account
	k.bankKeeper.SendCoins(ctx, moduleAcct, scavenger, reward)

	// save the updated scavenge
	k.SetScavenge(ctx, scavenge)

	return &types.MsgRevealSolutionResponse{}, nil
}
