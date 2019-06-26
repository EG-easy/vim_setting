Vim�UnDo� r�μ���(��fk��9�`�\�;��*����   �           "                       ]	�f    _�                     %        ����                                                                                                                                                                                                                                                                                                                            %           %           V        ]	�L     �   $   %          3var DefaultNodeHome = os.ExpandEnv("$HOME/.votumd")5�_�                    %        ����                                                                                                                                                                                                                                                                                                                            %           %           V        ]	�M     �   $   %           5�_�                            ����                                                                                                                                                                                                                                                                                                                            %           %           V        ]	�O    �               �   package main       import (   	"encoding/json"   	"fmt"   	"io"   	"io/ioutil"   	"os"   	"path/filepath"   
	"strings"       	"github.com/spf13/cobra"   	"github.com/spf13/viper"       6	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"   )	sdk "github.com/cosmos/cosmos-sdk/types"       3	abci "github.com/tendermint/tendermint/abci/types"   .	cfg "github.com/tendermint/tendermint/config"   /	dbm "github.com/tendermint/tendermint/libs/db"   1	tmtypes "github.com/tendermint/tendermint/types"       *	"github.com/tendermint/tendermint/crypto"   /	"github.com/tendermint/tendermint/libs/common"   ,	"github.com/tendermint/tendermint/libs/log"   #	"github.com/tendermint/tmlibs/cli"       &	"github.com/cosmos/cosmos-sdk/client"   %	"github.com/cosmos/cosmos-sdk/codec"   &	"github.com/cosmos/cosmos-sdk/server"   &	"github.com/cosmos/cosmos-sdk/x/auth"   &	"github.com/cosmos/cosmos-sdk/x/bank"       $	app "github.com/EG-easy/votumchain"   )       const (   	flagOverwrite = "overwrite"   )       func main() {   	fmt.Println("start votumd!")   #	cobra.EnableCommandSorting = false       	cdc := app.MakeCodec()   "	ctx := server.NewDefaultContext()       	rootCmd := &cobra.Command{   #		Use:               "cotumchaind",   6		Short:             "votumchain App Daemon (server)",   5		PersistentPreRunE: server.PersistentPreRunEFn(ctx),   	}       &	rootCmd.AddCommand(InitCmd(ctx, cdc))   3	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc))       =	server.AddCommands(ctx, cdc, rootCmd, newApp, appExporter())       ?	executor := cli.PrepareBaseCmd(rootCmd, "CT", DefaultNodeHome)   	err := executor.Execute()   	if err != nil {   		panic(err)   	}   }       Rfunc newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {   (	return app.NewVotumChainApp(logger, db)   }       'func appExporter() server.AppExporter {   V	return func(logger log.Logger, db dbm.DB, _ io.Writer, _ int64, _ bool, _ []string) (   7		json.RawMessage, []tmtypes.GenesisValidator, error) {   *		dapp := app.NewVotumChainApp(logger, db)   +		return dapp.ExportAppStateAndValidators()   	}   }       Dfunc InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {   	cmd := &cobra.Command{   		Use:   "init",   M		Short: "Initialize genesis config, priv-validator file, and p2p-node file",   		Args:  cobra.NoArgs,   2		RunE: func(_ *cobra.Command, _ []string) error {   			config := ctx.Config   0			config.SetRoot(viper.GetString(cli.HomeFlag))       1			chainID := viper.GetString(client.FlagChainID)   			if chainID == "" {   =				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))   			}   >			_, pk, err := gaiaInit.InitializeNodeValidatorFiles(config)   			if err != nil {   				return err   			}       			var appState json.RawMessage   "			genFile := config.GenesisFile()       C			if !viper.GetBool(flagOverwrite) && common.FileExists(genFile) {   F				return fmt.Errorf("genesis.json file already exists: %v", genFile)   			}       			genesis := app.GenesisState{   )				AuthData: auth.DefaultGenesisState(),   )				BankData: bank.DefaultGenesisState(),   			}       8			appState, err = codec.MarshalJSONIndent(cdc, genesis)   			if err != nil {   				return err   			}       2			_, _, validator, err := SimpleAppGenTx(cdc, pk)   			if err != nil {   				return err   			}       w			if err = gaiaInit.ExportGenesisFile(genFile, chainID, []tmtypes.GenesisValidator{validator}, appState); err != nil {   				return err   			}       V			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)       q			fmt.Printf("Initialized nsd configuration and bootstrapping files in %s ...\n", viper.GetString(cli.HomeFlag))   			return nil   		},   	}       K	cmd.Flags().String(cli.HomeFlag, DefaultNodeHome, "node's home directory")   l	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")   P	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")       	return cmd   }       Qfunc AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {   	cmd := &cobra.Command{   9		Use:   "add-genesis-account [address] [coins[,coins]]",   /		Short: "Adds an account to the genesis file",   		Args:  cobra.ExactArgs(2),   		Long: strings.TrimSpace(`   VAdds accounts to the genesis file so that you can start a chain with coins in the CLI:       d$ votumchaind add-genesis-account cosmos1tse7r2fadvlrrgau3pa0ss7cqh55wrv6y9alwh 1000STAKE,1000mycoin   `),   5		RunE: func(_ *cobra.Command, args []string) error {   1			addr, err := sdk.AccAddressFromBech32(args[0])   			if err != nil {   				return err   			}   (			coins, err := sdk.ParseCoins(args[1])   			if err != nil {   				return err   			}   			coins.Sort()        			var genDoc tmtypes.GenesisDoc   			config := ctx.Config   "			genFile := config.GenesisFile()   #			if !common.FileExists(genFile) {   K				return fmt.Errorf("%s does not exist, run `gaiad init` first", genFile)   			}   /			genContents, err := ioutil.ReadFile(genFile)   			if err != nil {   			}       A			if err = cdc.UnmarshalJSON(genContents, &genDoc); err != nil {   				return err   			}        			var appState app.GenesisState   G			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {   				return err   			}       /			for _, stateAcc := range appState.Accounts {   &				if stateAcc.Address.Equals(addr) {   Q					return fmt.Errorf("the application state already contains account %v", addr)   				}   			}       .			acc := auth.NewBaseAccountWithAddress(addr)   			acc.Coins = coins   6			appState.Accounts = append(appState.Accounts, &acc)   1			appStateJSON, err := cdc.MarshalJSON(appState)   			if err != nil {   				return err   			}       ^			return gaiaInit.ExportGenesisFile(genFile, genDoc.ChainID, genDoc.Validators, appStateJSON)   		},   	}   	return cmd   }       9func SimpleAppGenTx(cdc *codec.Codec, pk crypto.PubKey) (   U	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {   .	addr, secret, err := server.GenerateCoinKey()   	if err != nil {   		return   	}       $	bz, err := cdc.MarshalJSON(struct {   #		Addr sdk.AccAddress `json:"addr"`   		}{addr})   	if err != nil {   		return   	}       	appGenTx = json.RawMessage(bz)       ?	bz, err = cdc.MarshalJSON(map[string]string{"secret": secret})   	if err != nil {   		return   	}       	cliPrint = json.RawMessage(bz)       &	validator = tmtypes.GenesisValidator{   		PubKey: pk,   		Power:  10,   	}       	return   }5�_�                    :   /    ����                                                                                                                                                                                                                                                                                                                            %           %           V        ]	�V     �   9   ;   �      ?	executor := cli.PrepareBaseCmd(rootCmd, "CT", DefaultNodeHome)5�_�                    :   2    ����                                                                                                                                                                                                                                                                                                                            %           %           V        ]	�Y    �               �   package main       import (   	"encoding/json"   	"fmt"   	"io"   	"io/ioutil"   	"path/filepath"   
	"strings"       	"github.com/spf13/cobra"   	"github.com/spf13/viper"       6	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"   )	sdk "github.com/cosmos/cosmos-sdk/types"       3	abci "github.com/tendermint/tendermint/abci/types"   .	cfg "github.com/tendermint/tendermint/config"   /	dbm "github.com/tendermint/tendermint/libs/db"   1	tmtypes "github.com/tendermint/tendermint/types"       *	"github.com/tendermint/tendermint/crypto"   /	"github.com/tendermint/tendermint/libs/common"   ,	"github.com/tendermint/tendermint/libs/log"   #	"github.com/tendermint/tmlibs/cli"       &	"github.com/cosmos/cosmos-sdk/client"   %	"github.com/cosmos/cosmos-sdk/codec"   &	"github.com/cosmos/cosmos-sdk/server"   &	"github.com/cosmos/cosmos-sdk/x/auth"   &	"github.com/cosmos/cosmos-sdk/x/bank"       $	app "github.com/EG-easy/votumchain"   )       const (   	flagOverwrite = "overwrite"   )       func main() {   	fmt.Println("start votumd!")   #	cobra.EnableCommandSorting = false       	cdc := app.MakeCodec()   "	ctx := server.NewDefaultContext()       	rootCmd := &cobra.Command{   #		Use:               "cotumchaind",   6		Short:             "votumchain App Daemon (server)",   5		PersistentPreRunE: server.PersistentPreRunEFn(ctx),   	}       &	rootCmd.AddCommand(InitCmd(ctx, cdc))   3	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc))       =	server.AddCommands(ctx, cdc, rootCmd, newApp, appExporter())       C	executor := cli.PrepareBaseCmd(rootCmd, "CT", app.DefaultNodeHome)   	err := executor.Execute()   	if err != nil {   		panic(err)   	}   }       Rfunc newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {   (	return app.NewVotumChainApp(logger, db)   }       'func appExporter() server.AppExporter {   V	return func(logger log.Logger, db dbm.DB, _ io.Writer, _ int64, _ bool, _ []string) (   7		json.RawMessage, []tmtypes.GenesisValidator, error) {   *		dapp := app.NewVotumChainApp(logger, db)   +		return dapp.ExportAppStateAndValidators()   	}   }       Dfunc InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {   	cmd := &cobra.Command{   		Use:   "init",   M		Short: "Initialize genesis config, priv-validator file, and p2p-node file",   		Args:  cobra.NoArgs,   2		RunE: func(_ *cobra.Command, _ []string) error {   			config := ctx.Config   0			config.SetRoot(viper.GetString(cli.HomeFlag))       1			chainID := viper.GetString(client.FlagChainID)   			if chainID == "" {   =				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))   			}   >			_, pk, err := gaiaInit.InitializeNodeValidatorFiles(config)   			if err != nil {   				return err   			}       			var appState json.RawMessage   "			genFile := config.GenesisFile()       C			if !viper.GetBool(flagOverwrite) && common.FileExists(genFile) {   F				return fmt.Errorf("genesis.json file already exists: %v", genFile)   			}       			genesis := app.GenesisState{   )				AuthData: auth.DefaultGenesisState(),   )				BankData: bank.DefaultGenesisState(),   			}       8			appState, err = codec.MarshalJSONIndent(cdc, genesis)   			if err != nil {   				return err   			}       2			_, _, validator, err := SimpleAppGenTx(cdc, pk)   			if err != nil {   				return err   			}       w			if err = gaiaInit.ExportGenesisFile(genFile, chainID, []tmtypes.GenesisValidator{validator}, appState); err != nil {   				return err   			}       V			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)       q			fmt.Printf("Initialized nsd configuration and bootstrapping files in %s ...\n", viper.GetString(cli.HomeFlag))   			return nil   		},   	}       K	cmd.Flags().String(cli.HomeFlag, DefaultNodeHome, "node's home directory")   l	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")   P	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")       	return cmd   }       Qfunc AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {   	cmd := &cobra.Command{   9		Use:   "add-genesis-account [address] [coins[,coins]]",   /		Short: "Adds an account to the genesis file",   		Args:  cobra.ExactArgs(2),   		Long: strings.TrimSpace(`   VAdds accounts to the genesis file so that you can start a chain with coins in the CLI:       d$ votumchaind add-genesis-account cosmos1tse7r2fadvlrrgau3pa0ss7cqh55wrv6y9alwh 1000STAKE,1000mycoin   `),   5		RunE: func(_ *cobra.Command, args []string) error {   1			addr, err := sdk.AccAddressFromBech32(args[0])   			if err != nil {   				return err   			}   (			coins, err := sdk.ParseCoins(args[1])   			if err != nil {   				return err   			}   			coins.Sort()        			var genDoc tmtypes.GenesisDoc   			config := ctx.Config   "			genFile := config.GenesisFile()   #			if !common.FileExists(genFile) {   K				return fmt.Errorf("%s does not exist, run `gaiad init` first", genFile)   			}   /			genContents, err := ioutil.ReadFile(genFile)   			if err != nil {   			}       A			if err = cdc.UnmarshalJSON(genContents, &genDoc); err != nil {   				return err   			}        			var appState app.GenesisState   G			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {   				return err   			}       /			for _, stateAcc := range appState.Accounts {   &				if stateAcc.Address.Equals(addr) {   Q					return fmt.Errorf("the application state already contains account %v", addr)   				}   			}       .			acc := auth.NewBaseAccountWithAddress(addr)   			acc.Coins = coins   6			appState.Accounts = append(appState.Accounts, &acc)   1			appStateJSON, err := cdc.MarshalJSON(appState)   			if err != nil {   				return err   			}       ^			return gaiaInit.ExportGenesisFile(genFile, genDoc.ChainID, genDoc.Validators, appStateJSON)   		},   	}   	return cmd   }       9func SimpleAppGenTx(cdc *codec.Codec, pk crypto.PubKey) (   U	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {   .	addr, secret, err := server.GenerateCoinKey()   	if err != nil {   		return   	}       $	bz, err := cdc.MarshalJSON(struct {   #		Addr sdk.AccAddress `json:"addr"`   		}{addr})   	if err != nil {   		return   	}       	appGenTx = json.RawMessage(bz)       ?	bz, err = cdc.MarshalJSON(map[string]string{"secret": secret})   	if err != nil {   		return   	}       	cliPrint = json.RawMessage(bz)       &	validator = tmtypes.GenesisValidator{   		PubKey: pk,   		Power:  10,   	}       	return   }5�_�                    �   "    ����                                                                                                                                                                                                                                                                                                                            %           %           V        ]	�a     �      �   �      K	cmd.Flags().String(cli.HomeFlag, DefaultNodeHome, "node's home directory")5�_�                     �       ����                                                                                                                                                                                                                                                                                                                            %           %           V        ]	�e    �               �   package main       import (   	"encoding/json"   	"fmt"   	"io"   	"io/ioutil"   	"path/filepath"   
	"strings"       	"github.com/spf13/cobra"   	"github.com/spf13/viper"       6	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"   )	sdk "github.com/cosmos/cosmos-sdk/types"       3	abci "github.com/tendermint/tendermint/abci/types"   .	cfg "github.com/tendermint/tendermint/config"   /	dbm "github.com/tendermint/tendermint/libs/db"   1	tmtypes "github.com/tendermint/tendermint/types"       *	"github.com/tendermint/tendermint/crypto"   /	"github.com/tendermint/tendermint/libs/common"   ,	"github.com/tendermint/tendermint/libs/log"   #	"github.com/tendermint/tmlibs/cli"       &	"github.com/cosmos/cosmos-sdk/client"   %	"github.com/cosmos/cosmos-sdk/codec"   &	"github.com/cosmos/cosmos-sdk/server"   &	"github.com/cosmos/cosmos-sdk/x/auth"   &	"github.com/cosmos/cosmos-sdk/x/bank"       $	app "github.com/EG-easy/votumchain"   )       const (   	flagOverwrite = "overwrite"   )       func main() {   	fmt.Println("start votumd!")   #	cobra.EnableCommandSorting = false       	cdc := app.MakeCodec()   "	ctx := server.NewDefaultContext()       	rootCmd := &cobra.Command{   #		Use:               "cotumchaind",   6		Short:             "votumchain App Daemon (server)",   5		PersistentPreRunE: server.PersistentPreRunEFn(ctx),   	}       &	rootCmd.AddCommand(InitCmd(ctx, cdc))   3	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc))       =	server.AddCommands(ctx, cdc, rootCmd, newApp, appExporter())       C	executor := cli.PrepareBaseCmd(rootCmd, "CT", app.DefaultNodeHome)   	err := executor.Execute()   	if err != nil {   		panic(err)   	}   }       Rfunc newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {   (	return app.NewVotumChainApp(logger, db)   }       'func appExporter() server.AppExporter {   V	return func(logger log.Logger, db dbm.DB, _ io.Writer, _ int64, _ bool, _ []string) (   7		json.RawMessage, []tmtypes.GenesisValidator, error) {   *		dapp := app.NewVotumChainApp(logger, db)   +		return dapp.ExportAppStateAndValidators()   	}   }       Dfunc InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {   	cmd := &cobra.Command{   		Use:   "init",   M		Short: "Initialize genesis config, priv-validator file, and p2p-node file",   		Args:  cobra.NoArgs,   2		RunE: func(_ *cobra.Command, _ []string) error {   			config := ctx.Config   0			config.SetRoot(viper.GetString(cli.HomeFlag))       1			chainID := viper.GetString(client.FlagChainID)   			if chainID == "" {   =				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))   			}   >			_, pk, err := gaiaInit.InitializeNodeValidatorFiles(config)   			if err != nil {   				return err   			}       			var appState json.RawMessage   "			genFile := config.GenesisFile()       C			if !viper.GetBool(flagOverwrite) && common.FileExists(genFile) {   F				return fmt.Errorf("genesis.json file already exists: %v", genFile)   			}       			genesis := app.GenesisState{   )				AuthData: auth.DefaultGenesisState(),   )				BankData: bank.DefaultGenesisState(),   			}       8			appState, err = codec.MarshalJSONIndent(cdc, genesis)   			if err != nil {   				return err   			}       2			_, _, validator, err := SimpleAppGenTx(cdc, pk)   			if err != nil {   				return err   			}       w			if err = gaiaInit.ExportGenesisFile(genFile, chainID, []tmtypes.GenesisValidator{validator}, appState); err != nil {   				return err   			}       V			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)       q			fmt.Printf("Initialized nsd configuration and bootstrapping files in %s ...\n", viper.GetString(cli.HomeFlag))   			return nil   		},   	}       O	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")   l	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")   P	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")       	return cmd   }       Qfunc AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {   	cmd := &cobra.Command{   9		Use:   "add-genesis-account [address] [coins[,coins]]",   /		Short: "Adds an account to the genesis file",   		Args:  cobra.ExactArgs(2),   		Long: strings.TrimSpace(`   VAdds accounts to the genesis file so that you can start a chain with coins in the CLI:       d$ votumchaind add-genesis-account cosmos1tse7r2fadvlrrgau3pa0ss7cqh55wrv6y9alwh 1000STAKE,1000mycoin   `),   5		RunE: func(_ *cobra.Command, args []string) error {   1			addr, err := sdk.AccAddressFromBech32(args[0])   			if err != nil {   				return err   			}   (			coins, err := sdk.ParseCoins(args[1])   			if err != nil {   				return err   			}   			coins.Sort()        			var genDoc tmtypes.GenesisDoc   			config := ctx.Config   "			genFile := config.GenesisFile()   #			if !common.FileExists(genFile) {   K				return fmt.Errorf("%s does not exist, run `gaiad init` first", genFile)   			}   /			genContents, err := ioutil.ReadFile(genFile)   			if err != nil {   			}       A			if err = cdc.UnmarshalJSON(genContents, &genDoc); err != nil {   				return err   			}        			var appState app.GenesisState   G			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {   				return err   			}       /			for _, stateAcc := range appState.Accounts {   &				if stateAcc.Address.Equals(addr) {   Q					return fmt.Errorf("the application state already contains account %v", addr)   				}   			}       .			acc := auth.NewBaseAccountWithAddress(addr)   			acc.Coins = coins   6			appState.Accounts = append(appState.Accounts, &acc)   1			appStateJSON, err := cdc.MarshalJSON(appState)   			if err != nil {   				return err   			}       ^			return gaiaInit.ExportGenesisFile(genFile, genDoc.ChainID, genDoc.Validators, appStateJSON)   		},   	}   	return cmd   }       9func SimpleAppGenTx(cdc *codec.Codec, pk crypto.PubKey) (   U	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {   .	addr, secret, err := server.GenerateCoinKey()   	if err != nil {   		return   	}       $	bz, err := cdc.MarshalJSON(struct {   #		Addr sdk.AccAddress `json:"addr"`   		}{addr})   	if err != nil {   		return   	}       	appGenTx = json.RawMessage(bz)       ?	bz, err = cdc.MarshalJSON(map[string]string{"secret": secret})   	if err != nil {   		return   	}       	cliPrint = json.RawMessage(bz)       &	validator = tmtypes.GenesisValidator{   		PubKey: pk,   		Power:  10,   	}       	return   }5��