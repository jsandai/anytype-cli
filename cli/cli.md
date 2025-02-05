├── cli
├── cli.md
├── cmd/
│ ├── auth/
│ │ ├── auth.go
│ │ ├── login/
│ │ │ └── login.go
│ │ └── logout/
│ │ └── logout.go
│ ├── daemon/
│ │ └── daemon.go
│ ├── root.go
│ ├── server/
│ │ ├── server.go
│ │ ├── start/
│ │ │ └── start.go
│ │ ├── status/
│ │ │ └── status.go
│ │ └── stop/
│ │ └── stop.go
│ ├── shell/
│ │ └── shell.go
│ ├── space/
│ │ ├── autoapprove/
│ │ │ └── autoapprove.go
│ │ └── space.go
│ └── token/
│ ├── create/
│ │ └── create.go
│ └── token.go
├── daemon/
│ ├── daemon.go
│ ├── daemon_client.go
│ └── taskmanager.go
├── internal/
│ ├── auth.go
│ ├── client.go
│ ├── keyring.go
│ ├── space.go
│ ├── stream.go
│ └── token.go
├── main.go
└── tasks/
├── autoapprove.go
└── server.go

# File: cli.md

```markdown

```

# End of file: cli.md

# File: main.go

```text
package main

import (
	"github.com/anyproto/anytype-cli/cmd"
	"github.com/anyproto/anytype-cli/internal"
)

func main() {
	defer internal.CloseGRPCConnection()
	cmd.Execute()
}

```

# End of file: main.go

# File: cmd/root.go

```text
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/cmd/auth"
	"github.com/anyproto/anytype-cli/cmd/daemon"
	"github.com/anyproto/anytype-cli/cmd/server"
	"github.com/anyproto/anytype-cli/cmd/shell"
	"github.com/anyproto/anytype-cli/cmd/space"
	"github.com/anyproto/anytype-cli/cmd/token"
)

var rootCmd = &cobra.Command{
	Use:   "anyctl <command> <subcommand> [flags]",
	Short: "Anytype CLI",
	Long:  "Seamlessly interact with Anytype from the command line",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		auth.NewAuthCmd(),
		server.NewServerCmd(),
		shell.NewShellCmd(rootCmd),
		space.NewSpaceCmd(),
		token.NewTokenCmd(),
		daemon.NewDaemonCmd(),
	)
}

```

# End of file: cmd/root.go

# File: cmd/token/token.go

```text
package token

import (
	"github.com/spf13/cobra"

	tokenCreateCmd "github.com/anyproto/anytype-cli/cmd/token/create"
)

func NewTokenCmd() *cobra.Command {
	tokenCmd := &cobra.Command{
		Use:   "token <command>",
		Short: "Manage API tokens for authenticating requests to the REST API",
	}

	tokenCmd.AddCommand(tokenCreateCmd.NewCreateCmd())

	return tokenCmd
}

```

# End of file: cmd/token/token.go

# File: cmd/token/create/create.go

```text
package create

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/internal"
)

func NewCreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new API token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.CreateToken(); err != nil {
				return fmt.Errorf("X Failed to create token: %w", err)
			}

			fmt.Println("✓ Token created successfully.")
			return nil
		},
	}

	createCmd.Flags().String("mnemonic", "", "Provide mnemonic (12 words) for authentication")

	return createCmd
}

```

# End of file: cmd/token/create/create.go

# File: cmd/auth/auth.go

```text
package auth

import (
	"github.com/spf13/cobra"

	authLoginCmd "github.com/anyproto/anytype-cli/cmd/auth/login"
	authLogoutCmd "github.com/anyproto/anytype-cli/cmd/auth/logout"
)

func NewAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate with Anytype",
	}

	authCmd.AddCommand(authLoginCmd.NewLoginCmd())
	authCmd.AddCommand(authLogoutCmd.NewLogoutCmd())

	return authCmd
}

```

# End of file: cmd/auth/auth.go

# File: cmd/auth/logout/logout.go

```text
package logout

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/internal"
)

func NewLogoutCmd() *cobra.Command {
	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out and remove stored mnemonic from keychain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.Logout(); err != nil {
				return fmt.Errorf("X Failed to log out: %w", err)
			}
			fmt.Println("✓ Successfully logged out")
			return nil
		},
	}

	return logoutCmd
}

```

# End of file: cmd/auth/logout/logout.go

# File: cmd/auth/login/login.go

```text
package login

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/daemon"
	"github.com/anyproto/anytype-cli/internal"
)

func NewLoginCmd() *cobra.Command {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to your Anytype vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			mnemonic, _ := cmd.Flags().GetString("mnemonic")
			rootPath, _ := cmd.Flags().GetString("path")

			statusResp, err := daemon.SendTaskStatus("server")
			if err != nil || statusResp.Status != "running" {
				return fmt.Errorf("server is not running")
			}

			if err := internal.Login(mnemonic, rootPath); err != nil {
				return fmt.Errorf("X Failed to log in: %w", err)
			}
			fmt.Println("✓ Successfully logged in")
			return nil

		},
	}

	loginCmd.Flags().String("mnemonic", "", "Provide mnemonic (12 words) for authentication")
	loginCmd.Flags().String("path", "", "Provide custom root path for wallet recovery")

	return loginCmd
}

```

# End of file: cmd/auth/login/login.go

# File: cmd/shell/shell.go

```text
package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func NewShellCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Start the Anytype interactive shell",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Starting Anytype interactive shell. Type 'exit' to quit.")
			return runShell(rootCmd)
		},
	}
}

func runShell(rootCmd *cobra.Command) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">>> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		line = strings.TrimSpace(line)

		if line == "exit" || line == "quit" {
			fmt.Println("Goodbye!")
			return nil
		}

		if line == "" {
			continue // ignore empty input
		}

		args := strings.Split(line, " ")
		rootCmd.SetArgs(args)

		if err := rootCmd.Execute(); err != nil {
			fmt.Println("Command error:", err)
		}
	}
}

```

# End of file: cmd/shell/shell.go

# File: cmd/server/server.go

```text
package server

import (
	"github.com/spf13/cobra"

	serverStartCmd "github.com/anyproto/anytype-cli/cmd/server/start"
	serverStatusCmd "github.com/anyproto/anytype-cli/cmd/server/status"
	serverStopCmd "github.com/anyproto/anytype-cli/cmd/server/stop"
)

func NewServerCmd() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server <command>",
		Short: "Manage the Anytype local server",
	}

	serverCmd.AddCommand(serverStartCmd.NewStartCmd())
	serverCmd.AddCommand(serverStopCmd.NewStopCmd())
	serverCmd.AddCommand(serverStatusCmd.NewStatusCmd())

	return serverCmd
}

```

# End of file: cmd/server/server.go

# File: cmd/server/start/start.go

```text
package start

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/daemon"
	"github.com/anyproto/anytype-cli/internal"
)

func NewStartCmd() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the Anytype local server",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := daemon.SendTaskStart("server", nil)
			if err != nil {
				return fmt.Errorf("failed to start server: %w", err)
			}
			fmt.Println("✓ Server started successfully via daemon. Response:", resp.Status)
			time.Sleep(2 * time.Second) // wait for server to start

			mnemonic, err := internal.GetStoredMnemonic()
			if err == nil && mnemonic != "" {
				fmt.Println("ℹ Keychain mnemonic found. Attempting to login...")
				if err := internal.LoginAccount(mnemonic, ""); err != nil {
					fmt.Println("X Failed to login using keychain mnemonic:", err)
				} else {
					fmt.Println("✓ Successfully logged in using keychain mnemonic.")
				}
			} else {
				fmt.Println("ℹ No keychain mnemonic found. Please login using 'anyctl login'.")
			}
			return nil
		},
	}

	return startCmd
}

```

# End of file: cmd/server/start/start.go

# File: cmd/server/stop/stop.go

```text
package stop

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/daemon"
)

func NewStopCmd() *cobra.Command {
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the Anytype local server",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := daemon.SendTaskStop("server", nil)
			if err != nil {
				return fmt.Errorf("failed to stop server task: %w", err)
			}
			fmt.Println("✓ Server task stopped successfully. Response:", resp.Status)
			return nil
		},
	}

	return stopCmd
}

```

# End of file: cmd/server/stop/stop.go

# File: cmd/server/status/status.go

```text
package status

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/daemon"
)

func NewStatusCmd() *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get the status of the Anytype local server",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := daemon.SendTaskStatus("server")
			if err != nil {
				return fmt.Errorf("failed to get server status: %w", err)
			}
			fmt.Println("ℹ Server status:", resp.Status)
			return nil
		},
	}
	return statusCmd
}

```

# End of file: cmd/server/status/status.go

# File: cmd/daemon/daemon.go

```text
package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/daemon"
)

const (
	defaultDaemonAddr = "127.0.0.1:31010"
)

func NewDaemonCmd() *cobra.Command {
	var addr string

	daemonCmd := &cobra.Command{
		Use:   "daemon",
		Short: "Run the Anytype background daemon",
		Long:  "Run the Anytype daemon that manages background tasks (should be run as a system service).",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := cmd.Flags().GetString("addr")
			if err != nil {
				return err
			}
			fmt.Println("ℹ Starting daemon on", addr)
			return daemon.StartManager(addr)
		},
	}

	daemonCmd.Flags().StringVar(&addr, "addr", defaultDaemonAddr, "Address for the daemon to listen on")
	return daemonCmd
}

```

# End of file: cmd/daemon/daemon.go

# File: cmd/space/space.go

```text
package space

import (
	"github.com/spf13/cobra"

	spaceAutoapproveCmd "github.com/anyproto/anytype-cli/cmd/space/autoapprove"
)

func NewSpaceCmd() *cobra.Command {
	spaceCmd := &cobra.Command{
		Use:   "space <command>",
		Short: "Manage spaces",
	}

	spaceCmd.AddCommand(spaceAutoapproveCmd.NewAutoapproveCmd())

	return spaceCmd
}

```

# End of file: cmd/space/space.go

# File: cmd/space/autoapprove/autoapprove.go

```text
package autoapprove

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/daemon"
)

func NewAutoapproveCmd() *cobra.Command {
	var spaceID string
	var role string

	autoapproveCmd := &cobra.Command{
		Use:   "autoapprove",
		Short: "Start autoapproval of join requests for a space (runs in background)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if spaceID == "" {
				return fmt.Errorf("space id is required (use --space)")
			}
			if role == "" {
				return fmt.Errorf("role is required (use --role)")
			}

			params := map[string]string{
				"space": spaceID,
				"role":  role,
			}
			resp, err := daemon.SendTaskStart("autoapprove", params)
			if err != nil {
				return fmt.Errorf("failed to start autoapprove task: %w", err)
			}
			fmt.Printf("Autoapprove task started for space %s with role %s. Response: %s\n", spaceID, role, resp.Status)
			return nil
		},
	}

	autoapproveCmd.Flags().StringVar(&spaceID, "space", "", "ID of the space to monitor")
	autoapproveCmd.Flags().StringVar(&role, "role", "", "Role to assign to approved join requests (e.g., Editor, Viewer)")

	return autoapproveCmd
}

```

# End of file: cmd/space/autoapprove/autoapprove.go

# File: tasks/server.go

```text
package tasks

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// ServerTask is a background task that starts the server process.
// It spawns the server and waits until the given context is canceled.
func ServerTask(ctx context.Context) error {
	grpcPort := "31007"
	grpcWebPort := "31008"

	cmd := exec.Command("../dist/server")
	cmd.Env = append(os.Environ(),
		"ANYTYPE_GRPC_ADDR=127.0.0.1:"+grpcPort,
		"ANYTYPE_GRPCWEB_ADDR=127.0.0.1:"+grpcWebPort,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Run a goroutine to wait for the process to exit.
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Wait until either the task context is canceled or the process exits.
	select {
	case <-ctx.Done():
		syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		return <-done
	case err := <-done:
		return err
	}
}

```

# End of file: tasks/server.go

# File: tasks/autoapprove.go

```text
package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/anyproto/anytype-cli/internal"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func AutoapproveTask(ctx context.Context, spaceID, role string) error {
	var permissions model.ParticipantPermissions
	switch role {
	case "Editor":
		permissions = model.ParticipantPermissions_Writer
	case "Viewer":
		fallthrough
	default:
		permissions = model.ParticipantPermissions_Reader
	}

	token, err := internal.GetStoredToken()
	if err != nil || token == "" {
		return fmt.Errorf("failed to get stored token; are you logged in?")
	}

	er, err := internal.ListenForEvents(token)
	if err != nil {
		return fmt.Errorf("failed to start event listener: %w", err)
	}

	// Optionally, monitor the server status.
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				status, err := internal.IsGRPCServerRunning()
				if err != nil || !status {
					return
				}
			}
		}
	}()

	// Main loop: poll for join request events and approve them.
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			joinReq, err := internal.WaitForJoinRequestEvent(er, spaceID)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			if err := internal.ApproveJoinRequest(token, joinReq.SpaceId, joinReq.Identity, permissions); err != nil {
				fmt.Println("Failed to approve join request: %v", err)
			} else {
				fmt.Println("Successfully approved join request for identity %s", joinReq.Identity)
			}
		}
	}
}

```

# End of file: tasks/autoapprove.go

# File: internal/token.go

```text
package internal

func CreateToken() error {
	// TODO: implement
	return nil
}

```

# End of file: internal/token.go

# File: internal/auth.go

```text
package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/anyproto/anytype-heart/pb"
)

// LoginAccount performs the common steps for logging in with a given mnemonic and root path.
func LoginAccount(mnemonic, rootPath string) error {
	if rootPath == "" {
		rootPath = "/Users/jmetrikat/Library/Application Support/anytype/alpha/data"
	}

	client, err := GetGRPCClient()
	if err != nil {
		return fmt.Errorf("error connecting to gRPC server: %w", err)
	}

	// Create a context for the initial calls.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set initial parameters.
	_, err = client.InitialSetParameters(ctx, &pb.RpcInitialSetParametersRequest{
		Platform: "Mac",
		Version:  "0.0.0-test",
		Workdir:  "/Users/jmetrikat/Library/Application Support/anytype",
	})
	if err != nil {
		return fmt.Errorf("failed to set initial parameters: %w", err)
	}

	// Recover the wallet.
	_, err = client.WalletRecover(ctx, &pb.RpcWalletRecoverRequest{
		Mnemonic: mnemonic,
		RootPath: rootPath,
	})
	if err != nil {
		return fmt.Errorf("wallet recovery failed: %w", err)
	}

	// Create a session.
	resp, err := client.WalletCreateSession(ctx, &pb.RpcWalletCreateSessionRequest{
		Auth: &pb.RpcWalletCreateSessionRequestAuthOfMnemonic{
			Mnemonic: mnemonic,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	sessionToken := resp.Token
	err = SaveToken(sessionToken)
	fmt.Println("ℹ Session token:", sessionToken)
	if err != nil {
		return fmt.Errorf("failed to save session token: %w", err)
	}

	// Start listening for session events.
	er, err := ListenForEvents(sessionToken)
	if err != nil {
		return fmt.Errorf("failed to start event listener: %w", err)
	}

	// Recover the account.
	ctx = ClientContextWithAuth(sessionToken)
	_, err = client.AccountRecover(ctx, &pb.RpcAccountRecoverRequest{})
	if err != nil {
		return fmt.Errorf("account recovery failed: %w", err)
	}

	// Wait for the account ID.
	accountID, err := WaitForAccountID(er)
	if err != nil {
		return fmt.Errorf("error waiting for account ID: %w", err)
	}
	fmt.Println("ℹ Account ID:", accountID)

	// Select the account.
	_, err = client.AccountSelect(ctx, &pb.RpcAccountSelectRequest{
		DisableLocalNetworkSync: false,
		Id:                      accountID,
		JsonApiListenAddr:       "127.0.0.1:31009",
		RootPath:                rootPath,
	})
	if err != nil {
		return fmt.Errorf("failed to select account: %w", err)
	}

	return nil
}

func Login(mnemonic, rootPath string) error {
	usedStoredMnemonic := false
	if mnemonic == "" {
		mnemonic, err := GetStoredMnemonic()
		if err == nil && mnemonic != "" {
			fmt.Println("Using stored mnemonic from keychain.")
			usedStoredMnemonic = true
		} else {
			fmt.Print("Enter mnemonic (12 words): ")
			reader := bufio.NewReader(os.Stdin)
			mnemonic, _ = reader.ReadString('\n')
			mnemonic = strings.TrimSpace(mnemonic)
		}
	}

	if len(strings.Split(mnemonic, " ")) != 12 {
		return fmt.Errorf("mnemonic must be 12 words")
	}

	err := LoginAccount(mnemonic, rootPath)
	if err != nil {
		return fmt.Errorf("failed to log in: %w", err)
	}

	if !usedStoredMnemonic {
		if err := SaveMnemonic(mnemonic); err != nil {
			fmt.Println("Warning: failed to save mnemonic in keychain:", err)
		} else {
			fmt.Println("✓ Mnemonic saved to keychain.")
		}
	}

	return nil
}

func Logout() error {
	client, err := GetGRPCClient()
	if err != nil {
		fmt.Println("Failed to connect to gRPC server:", err)
	}

	token, err := GetStoredToken()
	if err != nil {
		return fmt.Errorf("failed to get stored token: %w", err)
	}

	ctx := ClientContextWithAuth(token)
	resp, err := client.AccountStop(ctx, &pb.RpcAccountStopRequest{
		RemoveData: false,
	})
	if err != nil {
		return fmt.Errorf("failed to log out: %w", err)
	}
	if resp.Error.Code != pb.RpcAccountStopResponseError_NULL {
		fmt.Println("Failed to log out:", resp.Error.Description)
	}

	resp2, err := client.WalletCloseSession(ctx, &pb.RpcWalletCloseSessionRequest{Token: token})
	if err != nil {
		return fmt.Errorf("failed to close session: %w", err)
	}
	if resp2.Error.Code != pb.RpcWalletCloseSessionResponseError_NULL {
		fmt.Println("Failed to close session:", resp2.Error.Description)
	}

	if err := DeleteStoredMnemonic(); err != nil {
		return fmt.Errorf("failed to delete stored mnemonic: %w", err)
	}
	fmt.Println("✓ Successfully logged out. Stored mnemonic removed.")

	return nil
}

```

# End of file: internal/auth.go

# File: internal/stream.go

```text
package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	pb "github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

// Singleton instance of EventReceiver
var (
	eventReceiverInstance *EventReceiver
	erOnce                sync.Once
)

// EventReceiver is a universal receiver that collects all incoming event messages.
type EventReceiver struct {
	lock   *sync.Mutex
	events []*pb.EventMessage
}

// ListenForEvents ensures a single EventReceiver instance is used.
func ListenForEvents(token string) (*EventReceiver, error) {
	var err error
	erOnce.Do(func() {
		eventReceiverInstance, err = startListeningForEvents(token)
	})
	if err != nil {
		return nil, err
	}
	return eventReceiverInstance, nil
}

// ListenForEvents starts the gRPC stream for events using the provided token.
// It returns an EventReceiver that will store all incoming events.
func startListeningForEvents(token string) (*EventReceiver, error) {
	client, err := GetGRPCClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get gRPC client: %w", err)
	}

	req := &pb.StreamRequest{
		Token: token,
	}
	stream, err := client.ListenSessionEvents(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to start event stream: %w", err)
	}

	er := &EventReceiver{
		lock:   &sync.Mutex{},
		events: make([]*pb.EventMessage, 0),
	}

	// Start a goroutine to continuously receive events from the stream.
	go func() {
		for {
			event, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Println("🔄 Event stream ended, reconnecting...")
				break
			}
			if err != nil {
				// Check for intentional close
				if err.Error() == "rpc error: code = Canceled desc = grpc: the client connection is closing" {
					break
				}
				fmt.Errorf("X Event stream error: %w\n", err)
				break
			}

			er.lock.Lock()
			er.events = append(er.events, event.Messages...)
			er.lock.Unlock()
		}
	}()

	return er, nil
}

// WaitForAccountID continuously checks the stored events until an accountShow event is found.
// It returns the account ID from that event.
func WaitForAccountID(er *EventReceiver) (string, error) {
	for {
		er.lock.Lock()
		// Process recent events first.
		for i := len(er.events) - 1; i >= 0; i-- {
			m := er.events[i]
			if m == nil {
				continue
			}
			if v := m.GetAccountShow(); v != nil && v.GetAccount() != nil {
				accountID := v.GetAccount().Id
				// Mark event as processed.
				er.events[i] = nil
				er.lock.Unlock()
				return accountID, nil
			}
		}
		er.lock.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

// WaitForJoinRequestEvent continuously polls the event receiver until it finds a join request for the specified space.
// It returns the join request details.
func WaitForJoinRequestEvent(er *EventReceiver, spaceID string) (*model.NotificationRequestToJoin, error) {
	for {
		er.lock.Lock()
		for i := len(er.events) - 1; i >= 0; i-- {
			m := er.events[i]
			if m == nil {
				continue
			}
			// Check for a notificationSend event with a join request.
			if ns := m.GetNotificationSend(); ns != nil && ns.Notification != nil && ns.Notification.GetRequestToJoin() != nil {
				req := ns.Notification.GetRequestToJoin()
				if req.SpaceId == spaceID {
					// Mark event as processed.
					er.events[i] = nil
					er.lock.Unlock()
					return req, nil
				}
			}
		}
		er.lock.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

```

# End of file: internal/stream.go

# File: internal/client.go

```text
package internal

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/grpc/metadata"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	clientInstance service.ClientCommandsClient
	grpcConn       *grpc.ClientConn
	once           sync.Once
)

// GetGRPCClient initializes (if needed) and returns the shared gRPC client
func GetGRPCClient() (service.ClientCommandsClient, error) {
	var err error

	// Ensure we only initialize once (singleton)
	once.Do(func() {
		grpcConn, err = grpc.NewClient("dns:///127.0.0.1:31007", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println("Failed to connect to gRPC server:", err)
			return
		}
		clientInstance = service.NewClientCommandsClient(grpcConn)
	})

	if err != nil {
		return nil, err
	}
	return clientInstance, nil
}

// CloseGRPCConnection ensures the connection is properly closed
func CloseGRPCConnection() {
	if grpcConn != nil {
		grpcConn.Close()
	}
}

// IsGRPCServerRunning checks if the gRPC server is reachable
func IsGRPCServerRunning() (bool, error) {
	client, err := GetGRPCClient()
	if err != nil {
		return false, err
	}

	_, err = client.AppGetVersion(context.Background(), &pb.RpcAppGetVersionRequest{})
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func ClientContextWithAuth(token string) context.Context {
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs("token", token))
}

```

# End of file: internal/client.go

# File: internal/keyring.go

```text
package internal

import "github.com/zalando/go-keyring"

const (
	// keyringService is the identifier for the CLI in the OS keychain.
	keyringService = "anytype-cli"
	// keyringMnemonicUser is the key under which the mnemonic is stored.
	keyringMnemonicUser = "mnemonic"
	// keyringTokenUser is the key under which the session token is stored.
	keyringTokenUser = "session-token"
)

// SaveMnemonic stores the mnemonic securely in the OS keychain.
func SaveMnemonic(mnemonic string) error {
	return keyring.Set(keyringService, keyringMnemonicUser, mnemonic)
}

// GetStoredMnemonic retrieves the mnemonic from the OS keychain.
func GetStoredMnemonic() (string, error) {
	return keyring.Get(keyringService, keyringMnemonicUser)
}

// DeleteStoredMnemonic removes the mnemonic from the OS keychain.
func DeleteStoredMnemonic() error {
	return keyring.Delete(keyringService, keyringMnemonicUser)
}

// SaveToken stores the session token securely in the OS keychain.
func SaveToken(token string) error {
	return keyring.Set(keyringService, keyringTokenUser, token)
}

// GetStoredToken retrieves the session token from the OS keychain.
func GetStoredToken() (string, error) {
	return keyring.Get(keyringService, keyringTokenUser)
}

// DeleteStoredToken removes the session token from the OS keychain.
func DeleteStoredToken() error {
	return keyring.Delete(keyringService, keyringTokenUser)
}

```

# End of file: internal/keyring.go

# File: internal/space.go

```text
package internal

import (
	"fmt"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func ApproveJoinRequest(token, spaceID, identity string, permissions model.ParticipantPermissions) error {
	client, err := GetGRPCClient()
	if err != nil {
		return err
	}
	ctx := ClientContextWithAuth(token)
	req := &pb.RpcSpaceRequestApproveRequest{
		SpaceId:     spaceID,
		Identity:    identity,
		Permissions: permissions,
	}
	_, err = client.SpaceRequestApprove(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to approve join request: %w", err)
	}
	return nil
}

```

# End of file: internal/space.go

# File: daemon/daemon_client.go

```text
package daemon

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	defaultDaemonAddr = "http://127.0.0.1:31010"
)

// SendTaskStart sends a start request for a given task.
func SendTaskStart(task string, params map[string]string) (*TaskResponse, error) {
	reqData := TaskRequest{Task: task, Params: params}
	b, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(defaultDaemonAddr+"/task/start", "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var taskResp TaskResponse
	err = json.Unmarshal(body, &taskResp)
	if err != nil {
		return nil, err
	}
	if taskResp.Status == "error" {
		return &taskResp, errors.New(taskResp.Error)
	}
	return &taskResp, nil
}

// SendTaskStop sends a stop request for a given task.
func SendTaskStop(task string, params map[string]string) (*TaskResponse, error) {
	reqData := TaskRequest{Task: task, Params: params}
	b, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(defaultDaemonAddr+"/task/stop", "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var taskResp TaskResponse
	err = json.Unmarshal(body, &taskResp)
	if err != nil {
		return nil, err
	}
	if taskResp.Status == "error" {
		return &taskResp, errors.New(taskResp.Error)
	}
	return &taskResp, nil

}

// SendTaskStatus sends a status request for a given task.
func SendTaskStatus(task string) (*TaskResponse, error) {
	resp, err := http.Get(defaultDaemonAddr + "/task/status?task=" + task)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var taskResp TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return nil, err
	}
	if taskResp.Status == "error" {
		return &taskResp, errors.New(taskResp.Error)
	}
	return &taskResp, nil
}

```

# End of file: daemon/daemon_client.go

# File: daemon/taskmanager.go

```text
package daemon

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Task is a background function that runs until the given context is canceled.
type Task func(ctx context.Context) error

// TaskManager tracks running background tasks.
type TaskManager struct {
	mu    sync.Mutex
	tasks map[string]context.CancelFunc
}

// defaultTaskManager is the singleton instance.
var defaultTaskManager = NewTaskManager()

// NewTaskManager returns a new task manager.
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]context.CancelFunc),
	}
}

// StartTask starts a new task with a unique ID.
// It returns an error if a task with that ID is already running.
func (tm *TaskManager) StartTask(id string, task Task) error {
	tm.mu.Lock()
	if _, exists := tm.tasks[id]; exists {
		tm.mu.Unlock()
		return errors.New("task already running")
	}
	ctx, cancel := context.WithCancel(context.Background())
	tm.tasks[id] = cancel
	tm.mu.Unlock()

	go func() {
		if err := task(ctx); err != nil {
			fmt.Printf("Task %s exited with error: %v", id, err)
		}
		tm.mu.Lock()
		delete(tm.tasks, id)
		tm.mu.Unlock()
	}()
	return nil
}

// StopTask cancels a running task by its ID.
func (tm *TaskManager) StopTask(id string) error {
	tm.mu.Lock()
	cancel, exists := tm.tasks[id]
	tm.mu.Unlock()
	if !exists {
		return errors.New("task not found")
	}
	cancel()
	return nil
}

// StopAll stops every running task.
func (tm *TaskManager) StopAll() {
	tm.mu.Lock()
	for id, cancel := range tm.tasks {
		cancel()
		delete(tm.tasks, id)
	}
	tm.mu.Unlock()
}

// GetTaskManager returns the singleton instance.
func GetTaskManager() *TaskManager {
	return defaultTaskManager
}

```

# End of file: daemon/taskmanager.go

# File: daemon/daemon.go

```text
package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anyproto/anytype-cli/tasks"
)

// TaskRequest is used by the HTTP API.
type TaskRequest struct {
	Task   string            `json:"task"`
	Params map[string]string `json:"params,omitempty"`
}

// TaskResponse is returned by the HTTP API.
type TaskResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// Manager wraps the in‑memory task manager and exposes HTTP handlers.
type Manager struct {
	taskManager *TaskManager
	mux         *http.ServeMux
}

// NewManager returns a new Manager.
func NewManager() *Manager {
	m := &Manager{
		taskManager: GetTaskManager(),
		mux:         http.NewServeMux(),
	}
	m.routes()
	return m
}

// routes sets up the HTTP endpoints.
func (m *Manager) routes() {
	m.mux.HandleFunc("/task/start", m.handleStartTask)
	m.mux.HandleFunc("/task/stop", m.handleStopTask)
	m.mux.HandleFunc("/task/status", m.handleStatusTask)
}

// ServeHTTP satisfies the http.Handler interface.
func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

// handleStartTask processes a POST request to start a task.
func (m *Manager) handleStartTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var err error
	switch req.Task {
	case "server":
		err = m.taskManager.StartTask("server", tasks.ServerTask)
	case "autoapprove":
		spaceID, ok := req.Params["space"]
		if !ok || spaceID == "" {
			http.Error(w, "missing 'space' param", http.StatusBadRequest)
			return
		}
		role := req.Params["role"]
		taskID := "autoapprove-" + spaceID
		err = m.taskManager.StartTask(taskID, func(ctx context.Context) error {
			return tasks.AutoapproveTask(ctx, spaceID, role)
		})
	default:
		http.Error(w, "unknown task", http.StatusBadRequest)
		return
	}

	resp := TaskResponse{}
	if err != nil {
		resp.Status = "error"
		resp.Error = err.Error()
	} else {
		resp.Status = "started"
	}
	json.NewEncoder(w).Encode(resp)
}

// handleStopTask processes a POST request to stop a task.
func (m *Manager) handleStopTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var err error
	switch req.Task {
	case "server":
		err = m.taskManager.StopTask("server")
	case "autoapprove":
		spaceID, ok := req.Params["space"]
		if !ok || spaceID == "" {
			http.Error(w, "missing 'space' param", http.StatusBadRequest)
			return
		}
		taskID := "autoapprove-" + spaceID
		err = m.taskManager.StopTask(taskID)
	default:
		http.Error(w, "unknown task", http.StatusBadRequest)
		return
	}

	resp := TaskResponse{}
	if err != nil {
		resp.Status = "error"
		resp.Error = err.Error()
	} else {
		resp.Status = "stopped"
	}
	json.NewEncoder(w).Encode(resp)
}

// handleStatusTask processes a GET request to check a task’s status.
func (m *Manager) handleStatusTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	// Expect a query parameter like ?task=server or ?task=autoapprove-<spaceID>
	taskID := r.URL.Query().Get("task")
	if taskID == "" {
		http.Error(w, "missing task parameter", http.StatusBadRequest)
		return
	}

	m.taskManager.mu.Lock()
	_, exists := m.taskManager.tasks[taskID]
	m.taskManager.mu.Unlock()

	resp := TaskResponse{}
	if exists {
		resp.Status = "running"
	} else {
		resp.Status = "stopped"
	}

	json.NewEncoder(w).Encode(resp)
}

// StartManager launches the daemon's HTTP server.
func StartManager(addr string) error {
	manager := NewManager()
	srv := &http.Server{
		Addr:              addr,
		Handler:           manager,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Channel to signal when the server is done.
	done := make(chan struct{})
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Errorf("daemon ListenAndServe: %v", err)
		}
		close(done)
	}()

	// Set up channel on which to send signal notifications.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-quit
	fmt.Println("Daemon is shutting down...")

	// First, tell the task manager to stop all tasks.
	GetTaskManager().StopAll()
	fmt.Println("All managed tasks have been stopped.")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Errorf("daemon forced to shutdown: %v", err)
	}

	<-done
	fmt.Println("Daemon exiting")
	return nil
}

```

# End of file: daemon/daemon.go
