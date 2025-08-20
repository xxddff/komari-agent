package cmd

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/komari-monitor/komari-agent/cmd/flags"
	monitoring "github.com/komari-monitor/komari-agent/monitoring/unit"
	"github.com/komari-monitor/komari-agent/server"
	"github.com/komari-monitor/komari-agent/update"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "komari-agent",
	Short: "komari agent",
	Long:  `komari agent`,
	Run: func(cmd *cobra.Command, args []string) {
		// 从环境变量读取值
		if flags.CFAccessClientID == "" {
			flags.CFAccessClientID = viper.GetString("cf_access_client_id")
		}
		if flags.CFAccessClientSecret == "" {
			flags.CFAccessClientSecret = viper.GetString("cf_access_client_secret")
		}

		log.Println("Komari Agent", update.CurrentVersion)
		log.Println("Github Repo:", update.Repo)
		// Auto discovery
		if flags.AutoDiscoveryKey != "" {
			err := handleAutoDiscovery()
			if err != nil {
				log.Printf("Auto-discovery failed: %v", err)
				os.Exit(1)
			}
		}
		diskList, err := monitoring.DiskList()
		if err != nil {
			log.Println("Failed to get disk list:", err)
		}
		log.Println("Monitoring Mountpoints:", diskList)
		interfaceList, err := monitoring.InterfaceList()
		if err != nil {
			log.Println("Failed to get interface list:", err)
		}
		log.Println("Monitoring Interfaces:", interfaceList)

		// 忽略不安全的证书
		if flags.IgnoreUnsafeCert {
			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
		// 自动更新
		if !flags.DisableAutoUpdate {
			err := update.CheckAndUpdate()
			if err != nil {
				log.Println("[ERROR]", err)
			}
			go update.DoUpdateWorks()
		}
		go server.DoUploadBasicInfoWorks()
		for {
			server.UpdateBasicInfo()
			server.EstablishWebSocketConnection()
		}
	},
}

func Execute() {
	for i, arg := range os.Args {
		if arg == "-autoUpdate" || arg == "--autoUpdate" {
			log.Println("WARNING: The -autoUpdate flag is deprecated in version 0.0.9 and later. Use --disable-auto-update to configure auto-update behavior.")
			// 从参数列表中移除该参数，防止cobra解析错误
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			break
		}
	}

	if err := RootCmd.Execute(); err != nil {
		log.Println(err)
	}
}

func init() {
	// 设置环境变量前缀
	viper.SetEnvPrefix("KOMARI")
	viper.AutomaticEnv()

	RootCmd.PersistentFlags().StringVarP(&flags.Token, "token", "t", "", "API token")
	//RootCmd.MarkPersistentFlagRequired("token")
	RootCmd.PersistentFlags().StringVarP(&flags.Endpoint, "endpoint", "e", "", "API endpoint")
	RootCmd.MarkPersistentFlagRequired("endpoint")
	RootCmd.PersistentFlags().StringVar(&flags.AutoDiscoveryKey, "auto-discovery", "", "Auto discovery key for the agent")
	RootCmd.PersistentFlags().BoolVar(&flags.DisableAutoUpdate, "disable-auto-update", false, "Disable automatic updates")
	RootCmd.PersistentFlags().BoolVar(&flags.DisableWebSsh, "disable-web-ssh", false, "Disable remote control(web ssh and rce)")
	RootCmd.PersistentFlags().BoolVar(&flags.MemoryModeAvailable, "memory-mode-available", false, "Report memory as available instead of used.")
	RootCmd.PersistentFlags().Float64VarP(&flags.Interval, "interval", "i", 1.0, "Interval in seconds")
	RootCmd.PersistentFlags().BoolVarP(&flags.IgnoreUnsafeCert, "ignore-unsafe-cert", "u", false, "Ignore unsafe certificate errors")
	RootCmd.PersistentFlags().IntVarP(&flags.MaxRetries, "max-retries", "r", 3, "Maximum number of retries")
	RootCmd.PersistentFlags().IntVarP(&flags.ReconnectInterval, "reconnect-interval", "c", 5, "Reconnect interval in seconds")
	RootCmd.PersistentFlags().IntVar(&flags.InfoReportInterval, "info-report-interval", 5, "Interval in minutes for reporting basic info")
	RootCmd.PersistentFlags().StringVar(&flags.IncludeNics, "include-nics", "", "Comma-separated list of network interfaces to include")
	RootCmd.PersistentFlags().StringVar(&flags.ExcludeNics, "exclude-nics", "", "Comma-separated list of network interfaces to exclude")
	RootCmd.PersistentFlags().StringVar(&flags.IncludeMountpoints, "include-mountpoint", "", "Semicolon-separated list of mount points to include for disk statistics")
	RootCmd.PersistentFlags().IntVar(&flags.MonthRotate, "month-rotate", 0, "Month reset for network statistics (0 to disable)")
	RootCmd.PersistentFlags().StringVar(&flags.CFAccessClientID, "cf-access-client-id", "", "Cloudflare Access Client ID")
	RootCmd.PersistentFlags().StringVar(&flags.CFAccessClientSecret, "cf-access-client-secret", "", "Cloudflare Access Client Secret")

	// 绑定环境变量
	viper.BindPFlag("token", RootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("endpoint", RootCmd.PersistentFlags().Lookup("endpoint"))
	viper.BindPFlag("cf_access_client_id", RootCmd.PersistentFlags().Lookup("cf-access-client-id"))
	viper.BindPFlag("cf_access_client_secret", RootCmd.PersistentFlags().Lookup("cf-access-client-secret"))

	RootCmd.PersistentFlags().ParseErrorsWhitelist.UnknownFlags = true
}
