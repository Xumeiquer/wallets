package cmd

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/Xumeiquer/wallets/server/middleware/content"
	"github.com/Xumeiquer/wallets/server/middleware/headers"
	"github.com/Xumeiquer/wallets/server/middleware/persistence"
	"github.com/Xumeiquer/wallets/server/router"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	version = "v0.0.1"
)

var rootCmd = &cobra.Command{
	Use:   "serve",
	Short: "Wallet tracking webapp & server :: " + version,
	Long:  "",
	Run:   runRootCmd,
}

func init() {
	rootCmd.Flags().String("bind", "127.0.0.1", "ip to bind the server")
	rootCmd.Flags().Int("port", 9090, "port to listen connection from")
	rootCmd.Flags().String("database", "wallet.db", "path where resides the database")
	rootCmd.Flags().Bool("debug", false, "debug")
}

var StaticContent embed.FS

func Execute(assets embed.FS) {
	StaticContent = assets
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runRootCmd(cmd *cobra.Command, args []string) {
	if !areFlagsValid(cmd.Flags()) {
		log.Fatal(errors.New("ERROR: flags are not valid"))
	}

	ip, _ := cmd.Flags().GetString("bind")
	port, _ := cmd.Flags().GetInt("port")
	dbPath, _ := cmd.Flags().GetString("database")
	debug, _ := cmd.Flags().GetBool("debug")

	db, err := persistence.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}
	defer db.Close()

	customHeaders := headers.Headers{}

	if debug {
		customHeaders["Cache-Control"] = "no-cache, no-store, must-revalidate"
		customHeaders["Pragma"] = "no-cache"
		customHeaders["Expires"] = "0"
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(persistence.Persistence(db))
	e.Use(headers.Set(customHeaders))
	e.Use(content.Static(StaticContent))

	router.SetRouting(e)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", ip, port)))
}

func areFlagsValid(flags *pflag.FlagSet) bool {
	debug, _ := flags.GetBool("debug")
	if debug {
		fmt.Println("Checking arguments")
	}
	validIP := func(ip string) bool {
		if debug {
			fmt.Printf("IP: %s -> %t\n", ip, net.ParseIP(ip) != nil)
		}
		return net.ParseIP(ip) != nil
	}
	validPort := func(port int) bool {
		if debug {
			fmt.Printf("Port: %d -> %t\n", port, port > 0 && port < 65535)
		}
		return port > 0 && port < 65535
	}
	validPath := func(path string) bool {
		if debug {
			fmt.Printf("DB Path: %s -> ", path)
		}
		pathInfo, err := os.Stat(path)
		if err != nil {
			if debug {
				fmt.Printf("%t\n", strings.HasSuffix(err.Error(), "no such file or directory"))
			}
			return strings.HasSuffix(err.Error(), "no such file or directory")
		}
		if debug {
			fmt.Printf("%t\n", !pathInfo.IsDir())
		}
		return !pathInfo.IsDir()
	}

	ip, _ := flags.GetString("bind")
	port, _ := flags.GetInt("port")
	dbPath, _ := flags.GetString("database")

	return validIP(ip) && validPort(port) && validPath(dbPath)
}
