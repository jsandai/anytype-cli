package core

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"

	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
)

// getDefaultDataPath returns the default data path for Anytype based on the operating system
func getDefaultDataPath() string {
	if dataPath := os.Getenv("DATA_PATH"); dataPath != "" {
		return dataPath
	}

	baseDir := getDefaultWorkDir()
	return filepath.Join(baseDir, "data")
}

// getDefaultWorkDir returns the default work directory for Anytype based on the operating system
func getDefaultWorkDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "anytype")
	case "windows":
		return filepath.Join(homeDir, "AppData", "Roaming", "anytype")
	default:
		return filepath.Join(homeDir, ".config", "anytype")
	}
}

// LoginAccount performs the common steps for logging in with a given mnemonic and root path.
func LoginAccount(mnemonic, rootPath, apiAddr string) error {
	if rootPath == "" {
		rootPath = getDefaultDataPath()
	}
	if apiAddr == "" {
		apiAddr = config.DefaultAPIAddress
	}

	var sessionToken string

	err := GRPCCallNoAuth(func(ctx context.Context, client service.ClientCommandsClient) error {
		_, err := client.InitialSetParameters(ctx, &pb.RpcInitialSetParametersRequest{
			Platform: runtime.GOOS,
			Version:  Version,
			Workdir:  getDefaultWorkDir(),
		})
		if err != nil {
			return fmt.Errorf("failed to set initial parameters: %w", err)
		}

		_, err = client.WalletRecover(ctx, &pb.RpcWalletRecoverRequest{
			Mnemonic: mnemonic,
			RootPath: rootPath,
		})
		if err != nil {
			return fmt.Errorf("wallet recovery failed: %w", err)
		}

		resp, err := client.WalletCreateSession(ctx, &pb.RpcWalletCreateSessionRequest{
			Auth: &pb.RpcWalletCreateSessionRequestAuthOfMnemonic{
				Mnemonic: mnemonic,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
		sessionToken = resp.Token
		return nil
	})

	if err != nil {
		return err
	}

	err = SaveToken(sessionToken)
	if err != nil {
		return fmt.Errorf("failed to save session token: %w", err)
	}

	er, err := ListenForEvents(sessionToken)
	if err != nil {
		return fmt.Errorf("failed to start event listener: %w", err)
	}

	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		_, err := client.AccountRecover(ctx, &pb.RpcAccountRecoverRequest{})
		if err != nil {
			return fmt.Errorf("account recovery failed: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	accountId, err := WaitForAccountId(er)
	if err != nil {
		return fmt.Errorf("error waiting for account Id: %w", err)
	}
	output.Info("Account Id: %s", accountId)

	var techSpaceId string
	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		resp, err := client.AccountSelect(ctx, &pb.RpcAccountSelectRequest{
			Id:                accountId,
			JsonApiListenAddr: apiAddr,
			RootPath:          rootPath,
		})
		if err != nil {
			return fmt.Errorf("failed to select account: %w", err)
		}
		if resp.Account != nil && resp.Account.Info != nil {
			techSpaceId = resp.Account.Info.TechSpaceId
		}
		return nil
	})
	if err != nil {
		return err
	}

	configMgr := config.GetConfigManager()
	if err := configMgr.Load(); err != nil {
		output.Warning("failed to load config: %v", err)
	}
	if err := configMgr.SetAccountId(accountId); err != nil {
		output.Warning("failed to save account Id: %v", err)
	}
	if techSpaceId != "" {
		if err := configMgr.SetTechSpaceId(techSpaceId); err != nil {
			output.Warning("failed to save tech space Id: %v", err)
		}
	}

	return nil
}

func ValidateMnemonic(mnemonic string) error {
	if mnemonic == "" {
		return fmt.Errorf("mnemonic cannot be empty")
	}

	words := strings.Fields(mnemonic)
	if len(words) != 12 {
		return fmt.Errorf("mnemonic must be exactly 12 words, got %d", len(words))
	}

	return nil
}

func Login(mnemonic, rootPath, apiAddr string) error {
	usedStoredMnemonic := false
	if mnemonic == "" {
		storedMnemonic, err := GetStoredMnemonic()
		if err == nil && storedMnemonic != "" {
			mnemonic = storedMnemonic
			output.Info("Using stored mnemonic from keychain.")
			usedStoredMnemonic = true
		} else {
			output.Print("Enter mnemonic (12 words): ")
			reader := bufio.NewReader(os.Stdin)
			mnemonic, _ = reader.ReadString('\n')
			mnemonic = strings.TrimSpace(mnemonic)
		}
	}

	if err := ValidateMnemonic(mnemonic); err != nil {
		return err
	}

	err := LoginAccount(mnemonic, rootPath, apiAddr)
	if err != nil {
		return fmt.Errorf("failed to log in: %w", err)
	}

	if !usedStoredMnemonic {
		if err := SaveMnemonic(mnemonic); err != nil {
			output.Warning("failed to save mnemonic in keychain: %v", err)
		} else {
			output.Success("Mnemonic saved to keychain.")
		}
	}

	return nil
}

func Logout() error {
	token, err := GetStoredToken()
	if err != nil {
		return fmt.Errorf("failed to get stored token: %w", err)
	}

	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		resp, err := client.AccountStop(ctx, &pb.RpcAccountStopRequest{
			RemoveData: false,
		})
		if err != nil {
			return fmt.Errorf("failed to log out: %w", err)
		}
		if resp.Error.Code != pb.RpcAccountStopResponseError_NULL {
			output.Warning("Failed to log out: %s", resp.Error.Description)
		}

		resp2, err := client.WalletCloseSession(ctx, &pb.RpcWalletCloseSessionRequest{Token: token})
		if err != nil {
			return fmt.Errorf("failed to close session: %w", err)
		}
		if resp2.Error.Code != pb.RpcWalletCloseSessionResponseError_NULL {
			output.Warning("Failed to close session: %s", resp2.Error.Description)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := DeleteStoredMnemonic(); err != nil {
		return fmt.Errorf("failed to delete stored mnemonic: %w", err)
	}

	if err := DeleteStoredToken(); err != nil {
		return fmt.Errorf("failed to delete stored token: %w", err)
	}

	configMgr := config.GetConfigManager()
	if err := configMgr.Delete(); err != nil {
		output.Warning("failed to clear config: %v", err)
	}

	CloseEventReceiver()

	return nil
}

// CreateWallet creates a new wallet with the given root path and returns the mnemonic and account Id
func CreateWallet(name, rootPath, apiAddr string) (string, string, error) {
	if rootPath == "" {
		rootPath = getDefaultDataPath()
	}
	if apiAddr == "" {
		apiAddr = config.DefaultAPIAddress
	}

	var mnemonic string
	var sessionToken string

	err := GRPCCallNoAuth(func(ctx context.Context, client service.ClientCommandsClient) error {
		_, err := client.InitialSetParameters(ctx, &pb.RpcInitialSetParametersRequest{
			Platform: runtime.GOOS,
			Version:  Version,
			Workdir:  getDefaultWorkDir(),
		})
		if err != nil {
			return fmt.Errorf("failed to set initial parameters: %w", err)
		}

		createResp, err := client.WalletCreate(ctx, &pb.RpcWalletCreateRequest{
			RootPath: rootPath,
		})
		if err != nil {
			return fmt.Errorf("wallet creation failed: %w", err)
		}
		mnemonic = createResp.Mnemonic

		sessionResp, err := client.WalletCreateSession(ctx, &pb.RpcWalletCreateSessionRequest{
			Auth: &pb.RpcWalletCreateSessionRequestAuthOfMnemonic{
				Mnemonic: mnemonic,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
		sessionToken = sessionResp.Token
		return nil
	})

	if err != nil {
		return "", "", err
	}

	err = SaveToken(sessionToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to save session token: %w", err)
	}

	_, err = ListenForEvents(sessionToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to start event listener: %w", err)
	}

	var accountId string
	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		resp, err := client.AccountCreate(ctx, &pb.RpcAccountCreateRequest{
			Name:              name,
			StorePath:         rootPath,
			JsonApiListenAddr: apiAddr,
		})
		if err != nil {
			return fmt.Errorf("account creation failed: %w", err)
		}
		accountId = resp.Account.Id
		return nil
	})
	if err != nil {
		return "", "", err
	}

	var techSpaceId string
	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		resp, err := client.AccountSelect(ctx, &pb.RpcAccountSelectRequest{
			Id:                accountId,
			JsonApiListenAddr: apiAddr,
			RootPath:          rootPath,
		})
		if err != nil {
			return fmt.Errorf("failed to select account: %w", err)
		}
		if resp.Account != nil && resp.Account.Info != nil {
			techSpaceId = resp.Account.Info.TechSpaceId
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}

	_ = SaveMnemonic(mnemonic)

	configMgr := config.GetConfigManager()
	if err := configMgr.Load(); err != nil {
		output.Warning("failed to load config: %v", err)
	}
	if err := configMgr.SetAccountId(accountId); err != nil {
		output.Warning("failed to save account Id: %v", err)
	}
	if techSpaceId != "" {
		if err := configMgr.SetTechSpaceId(techSpaceId); err != nil {
			output.Warning("failed to save tech space Id: %v", err)
		}
	}

	return mnemonic, accountId, nil
}
