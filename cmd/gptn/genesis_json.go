/*
   This file is part of go-palletone.
   go-palletone is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-palletone is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
*/
/*
 * @author PalletOne core developer Albert·Gou <dev@pallet.one>
 * @date 2018
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/palletone/go-palletone/cmd/console"
	"github.com/palletone/go-palletone/cmd/utils"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/configure"
	mp "github.com/palletone/go-palletone/consensus/mediatorplugin"
	"github.com/palletone/go-palletone/core"
	"github.com/palletone/go-palletone/core/gen"
	"github.com/palletone/go-palletone/dag/errors"
	"gopkg.in/urfave/cli.v1"
)

const defaultGenesisJsonPath = "./ptn-genesis.json"

var (
	GenesisTimestampFlag = cli.Int64Flag{
		Name:  "genesistime",
		Usage: "Replace timestamp from genesis.json with current time plus this many seconds (experts only!)",
		//		Value: 0,
	}

	GenesisJsonPathFlag = cli.StringFlag{
		Name:  "genesispath",
		Usage: "Path to create a Genesis State at.",
		Value: defaultGenesisJsonPath,
	}

	createGenesisJsonCommand = cli.Command{
		Action:    utils.MigrateFlags(createGenesisJson),
		Name:      "newgenesis",
		Usage:     "Create a genesis json file template",
		ArgsUsage: "<genesisJsonPath>",
		Flags: []cli.Flag{
			GenesisJsonPathFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
Create a json file for the genesis state of a new chain 
with an existing account or a newly created account.

If a well-formed JSON file exists at the path, 
it will be replaced with an example Genesis State.`,
	}
)

func getTokenAccount(ctx *cli.Context) (string, error) {
	confirm, err := console.Stdin.PromptConfirm(
		"Do you want to create a new account as the holder of the token?")
	if err != nil {
		utils.Fatalf("%v", err)
	}

	var account string
	if !confirm {
		account, err = console.Stdin.PromptInput("Please enter an existing account address: ")
		if err != nil || len(account) == 0 || !common.IsValidAddress(account) {
			errStr := "Invalid Token Account Address!"
			utils.Fatalf(errStr)
			return "", errors.New(errStr)
		}
	} else {
		account, err = initialAccount(ctx)
		if err != nil {
			utils.Fatalf("%v", err)
			return "", err
		}
	}

	return account, nil
}

func createExampleMediators(ctx *cli.Context, mcLen int) []*mp.MediatorConf {
	exampleMediators := make([]*mp.MediatorConf, mcLen, mcLen)
	for i := 0; i < mcLen; i++ {
		account, password, _ := createExampleAccount(ctx)
		secStr, pubStr := core.CreateInitDKS()

		exampleMediators[i] = &mp.MediatorConf{
			Address:     account,
			Password:    password,
			InitPrivKey: secStr,
			InitPubKey:  pubStr,
		}
	}

	return exampleMediators
}

// createGenesisJson, Create a json file for the genesis state of a new chain.
func createGenesisJson(ctx *cli.Context) error {
	var (
		genesisFile *os.File
		err         error
	)

	account, err := getTokenAccount(ctx)
	if err != nil {
		return err
	}

	// comment by Albert·Gou
	//account, _, err := createExampleAccount(ctx)
	//if err != nil {
	//	return err
	//}

	mcs := createExampleMediators(ctx, core.DefaultMediatorCount)
	nodeStr, err := getNodeInfo(ctx)
	if err != nil {
		return err
	}

	genesisState := createExampleGenesis(account, mcs, nodeStr)

	var genesisJson []byte
	genesisJson, err = json.MarshalIndent(genesisState, "", "  ")
	if err != nil {
		utils.Fatalf("%v", err)
		return err
	}

	genesisOut := getGenesisPath(ctx)

	err = os.MkdirAll(filepath.Dir(genesisOut), os.ModePerm)
	if err != nil {
		utils.Fatalf("%v", err)
		return err
	}

	genesisFile, err = os.Create(genesisOut)
	defer genesisFile.Close()
	if err != nil {
		utils.Fatalf("%v", err)
		return err
	}

	_, err = genesisFile.Write(genesisJson)
	if err != nil {
		utils.Fatalf("%v", err)
		return err
	}

	fmt.Println("Creating example genesis state in file " + genesisOut)

	modifyConfig(ctx, mcs)

	return nil
}

func modifyConfig(ctx *cli.Context, mediators []*mp.MediatorConf) error {
	cfg := &FullConfig{Node: defaultNodeConfig()}
	configPath := getConfigPath(ctx)

	// 加载配置文件中的配置信息到 cfg中
	err := loadConfig(configPath, cfg)
	if err != nil {
		utils.Fatalf("%v", err)
		return err
	}

	cfg.MediatorPlugin.EnableStaleProduction = true
	cfg.MediatorPlugin.Mediators = mediators

	err = makeConfigFile(cfg, configPath)
	if err != nil {
		utils.Fatalf("%v", err)
		return err
	}

	fmt.Println("Rewriting config file at:", configPath)

	return nil
}

func getGenesisPath(ctx *cli.Context) string {
	// Make sure we have a valid genesis JSON
	genesisOut := ctx.Args().First()

	// If no path is specified, the default path is used
	if len(genesisOut) == 0 {
		// utils.Fatalf("Must supply path to genesis JSON file")
		genesisOut = defaultGenesisJsonPath
	}

	return genesisOut
}

// initialAccount, create a initial account for a new account
func initialAccount(ctx *cli.Context) (string, error) {
	address, err := newAccount(ctx)
	if err != nil {
		utils.Fatalf("%v", err)
	}

	fmt.Printf("Initial token holder's account address: %s\n", address.String())

	return address.Str(), nil
}

func createExampleAccount(ctx *cli.Context) (addrStr, password string, err error) {
	password = mp.DefaultPassword
	address, err := createAccount(ctx, password)
	addrStr = address.Str()
	return
}

// createExampleGenesis, create the genesis state of new chain with the specified account
func createExampleGenesis(tokenAccount string, mediators []*mp.MediatorConf, nodeInfo string) *core.Genesis {
	SystemConfig := core.SystemConfig{
		DepositRate:               core.DefaultDepositRate,
		FoundationAddress:         core.DefaultFoundationAddress,
		DepositAmountForMediator:  core.DefaultDepositAmountForMediator,
		DepositAmountForJury:      core.DefaultDepositAmountForJury,
		DepositAmountForDeveloper: core.DefaultDepositAmountForDeveloper,
		DepositPeriod:             core.DefaultDepositPeriod,
	}

	initParams := core.NewChainParams()

	return &core.Genesis{
		Alias:        core.DefaultAlias,
		Version:      configure.Version,
		TokenAmount:  core.DefaultTokenAmount,
		TokenDecimal: core.DefaultTokenDecimal,
		ChainID:      core.DefaultChainID,
		TokenHolder:  tokenAccount,

		Text:                      core.DefaultText,
		SystemConfig:              SystemConfig,
		InitialParameters:         initParams,
		ImmutableParameters:       core.NewImmutChainParams(),
		InitialTimestamp:          gen.InitialTimestamp(initParams.MediatorInterval),
		InitialActiveMediators:    core.DefaultMediatorCount,
		InitialMediatorCandidates: initialMediatorCandidates(mediators, nodeInfo),
	}
}

func initialMediatorCandidates(mediators []*mp.MediatorConf, nodeInfo string) []*core.InitialMediator {
	mcLen := len(mediators)
	initialMediators := make([]*core.InitialMediator, mcLen)
	for i := 0; i < mcLen; i++ {
		im := core.NewInitialMediator()
		im.AddStr = mediators[i].Address
		im.InitPubKey = mediators[i].InitPubKey
		im.Node = nodeInfo
		initialMediators[i] = im
	}

	return initialMediators
}
