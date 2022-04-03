package types_test

import (
	"testing"

	"github.com/lavanet/lava/x/servicer/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

				StakeMapList: []types.StakeMap{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				SpecStakeStorageList: []types.SpecStakeStorage{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				BlockDeadlineForCallback: types.BlockDeadlineForCallback{
					Deadline: types.BlockNum{Num: 0},
				},
				UnstakingServicersAllSpecsList: []types.UnstakingServicersAllSpecs{
					{
						Id: 0,
					},
					{
						Id: 1,
					},
				},
				UnstakingServicersAllSpecsCount: 2,
				CurrentSessionStart: &types.CurrentSessionStart{
					Block: types.BlockNum{Num: 0},
				},
				PreviousSessionBlocks: &types.PreviousSessionBlocks{
					BlocksNum: 6,
				},
				SessionStorageForSpecList: []types.SessionStorageForSpec{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				EarliestSessionStart: &types.EarliestSessionStart{
					Block: types.BlockNum{Num: 0},
				},
				UniquePaymentStorageUserServicerList: []types.UniquePaymentStorageUserServicer{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				UserPaymentStorageList: []types.UserPaymentStorage{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				SessionPaymentsList: []types.SessionPayments{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated stakeMap",
			genState: &types.GenesisState{
				StakeMapList: []types.StakeMap{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated specStakeStorage",
			genState: &types.GenesisState{
				SpecStakeStorageList: []types.SpecStakeStorage{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated unstakingServicersAllSpecs",
			genState: &types.GenesisState{
				UnstakingServicersAllSpecsList: []types.UnstakingServicersAllSpecs{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid unstakingServicersAllSpecs count",
			genState: &types.GenesisState{
				UnstakingServicersAllSpecsList: []types.UnstakingServicersAllSpecs{
					{
						Id: 1,
					},
				},
				UnstakingServicersAllSpecsCount: 0,
			},
			valid: false,
		},
		{
			desc: "duplicated sessionStorageForSpec",
			genState: &types.GenesisState{
				SessionStorageForSpecList: []types.SessionStorageForSpec{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated uniquePaymentStorageUserServicer",
			genState: &types.GenesisState{
				UniquePaymentStorageUserServicerList: []types.UniquePaymentStorageUserServicer{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated userPaymentStorage",
			genState: &types.GenesisState{
				UserPaymentStorageList: []types.UserPaymentStorage{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated sessionPayments",
			genState: &types.GenesisState{
				SessionPaymentsList: []types.SessionPayments{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}