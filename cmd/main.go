package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"trcli/internal/api"
	"trcli/internal/config"
	"trcli/internal/output"
)

const version = "2.1.0"

const usage = `te_demo - TradeRevolution CLI

USAGE:
  te_demo <command> [flags]

CONFIG:
  config                       Show / set configuration

ACCOUNT:
  accounts                     List accounts
  accounts state <id>          Account state / balance
  accounts positions <id>      Open positions
  accounts orders <id>         Active orders
  accounts orders-history <id> Orders history
  accounts executions <id>     Executions (fills)
  accounts instruments <id>    Available instruments
  accounts statements <id>     Statement report
  accounts risk-counters <id>  Risk rules counters

TRADING:
  order place <accountId>      Place a new order
  order cancel <orderId>       Cancel order
  order cancel-all <accountId> Cancel all orders
  order modify <orderId>       Modify order
  position close <positionId>  Close position
  position close-all <accountId> Close all positions
  position modify <positionId> Modify SL/TP on position

MARKET DATA:
  quotes                       Get quotes
  depth                        Depth of market
  history                      Historical bars
  daily-bar                    Daily bar
  last-bar                     Last bar
  trades                       Recent trades

OTHER:
  version                      Show version
  help                         Show this help

EXAMPLES:
  te_demo setup
  te_demo accounts
  te_demo accounts state 12345
  te_demo accounts positions 12345
  te_demo accounts instruments 12345 --type forex
  te_demo quotes --tradableInstrumentId 100 --accountId 12345
  te_demo order place 12345 --side buy --type market --qty 1 --tradableInstrumentId 100 --validity DAY
  te_demo order cancel 67890
  te_demo position close 111

VERSION: ` + version

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(usage)
		os.Exit(0)
	}

	switch args[0] {
	case "help", "--help", "-h":
		fmt.Println(usage)
	case "version", "--version", "-v":
		fmt.Printf("te_demo %s\n", version)
	case "config":
		cmdConfig(args[1:])
	case "accounts":
		cmdAccounts(args[1:])
	case "quotes":
		cmdGet("/quotes", args[1:])
	case "depth":
		cmdGet("/depth", args[1:])
	case "history":
		cmdGet("/history", args[1:])
	case "daily-bar":
		cmdGet("/dailyBar", args[1:])
	case "last-bar":
		cmdGet("/lastBar", args[1:])
	case "trades":
		cmdGet("/trades", args[1:])
	case "sessions":
		cmdGet("/sessions", args[1:])
	case "order":
		cmdOrder(args[1:])
	case "position":
		cmdPosition(args[1:])
	default:
		output.Error(fmt.Sprintf("Unknown command: %s", args[0]))
		fmt.Println("Run 'te_demo help' for usage.")
		os.Exit(1)
	}
}

// ── setup ─────────────────────────────────────────────────────────────────────

// ── config ────────────────────────────────────────────────────────────────────

func cmdConfig(args []string) {
	cfg, err := config.Load()
	die(err)

	if v := flag(args, "--url"); v != "" {
		cfg.BaseURL = v
		die(cfg.Save())
		output.Success(fmt.Sprintf("Base URL set to: %s", v))
		return
	}
	if v := flag(args, "--token"); v != "" {
		cfg.Token = v
		die(cfg.Save())
		output.Success("Token saved")
		return
	}
	output.Header("Configuration")
	fmt.Printf("  Base URL:  %s\n", cfg.BaseURL)
	if cfg.Token != "" {
		masked := cfg.Token
		if len(masked) > 16 {
			masked = masked[:8] + "..." + masked[len(masked)-4:]
		}
		fmt.Printf("  Token:     %s\n", masked)
	} else {
		fmt.Printf("  Token:     (not set — run: te_demo config --token YOUR_TOKEN)\n")
	}
}

// ── accounts ──────────────────────────────────────────────────────────────────

func cmdAccounts(args []string) {
	if len(args) == 0 {
		c := mustClient()
		result, err := c.Get("/accounts", nil)
		dieAPI(err)
		output.JSON(result)
		return
	}

	sub := args[0]
	rest := args[1:]

	switch sub {
	case "state":
		id := requireArg(rest, 0, "account ID")
		c := mustClient()
		result, err := c.Get(fmt.Sprintf("/accounts/%s/state", id), nil)
		dieAPI(err)
		output.JSON(result)
	case "positions":
		id := requireArg(rest, 0, "account ID")
		c := mustClient()
		result, err := c.Get(fmt.Sprintf("/accounts/%s/positions", id), queryParams(rest[1:]))
		dieAPI(err)
		output.JSON(result)
	case "orders":
		id := requireArg(rest, 0, "account ID")
		c := mustClient()
		result, err := c.Get(fmt.Sprintf("/accounts/%s/orders", id), queryParams(rest[1:]))
		dieAPI(err)
		output.JSON(result)
	case "orders-history":
		id := requireArg(rest, 0, "account ID")
		c := mustClient()
		result, err := c.Get(fmt.Sprintf("/accounts/%s/ordersHistory", id), queryParams(rest[1:]))
		dieAPI(err)
		output.JSON(result)
	case "executions":
		id := requireArg(rest, 0, "account ID")
		c := mustClient()
		result, err := c.Get(fmt.Sprintf("/accounts/%s/executions", id), queryParams(rest[1:]))
		dieAPI(err)
		output.JSON(result)
	case "instruments":
		id := requireArg(rest, 0, "account ID")
		c := mustClient()
		result, err := c.Get(fmt.Sprintf("/accounts/%s/instruments", id), queryParams(rest[1:]))
		dieAPI(err)
		output.JSON(result)
	case "statements":
		id := requireArg(rest, 0, "account ID")
		c := mustClient()
		result, err := c.Get(fmt.Sprintf("/accounts/%s/statements", id), queryParams(rest[1:]))
		dieAPI(err)
		output.JSON(result)
	case "risk-counters":
		id := requireArg(rest, 0, "account ID")
		c := mustClient()
		result, err := c.Get(fmt.Sprintf("/accounts/%s/riskRulesCounters", id), nil)
		dieAPI(err)
		output.JSON(result)
	default:
		// treat as account ID → show state
		c := mustClient()
		result, err := c.Get(fmt.Sprintf("/accounts/%s/state", sub), nil)
		dieAPI(err)
		output.JSON(result)
	}
}

// ── market data ───────────────────────────────────────────────────────────────

func cmdGet(path string, args []string) {
	c := mustClient()
	result, err := c.Get(path, queryParams(args))
	dieAPI(err)
	output.JSON(result)
}

// ── order ─────────────────────────────────────────────────────────────────────

func cmdOrder(args []string) {
	if len(args) == 0 {
		output.Error("Usage: te_demo order <place|cancel|cancel-all|modify> ...")
		os.Exit(1)
	}
	sub, rest := args[0], args[1:]

	switch sub {
	case "place":
		accountID := requireArg(rest, 0, "account ID")
		c := mustClient()
		body := buildOrderBody(rest[1:])
		result, err := c.Post(fmt.Sprintf("/accounts/%s/orders", accountID), body)
		dieAPI(err)
		output.JSON(result)
	case "cancel":
		orderID := requireArg(rest, 0, "order ID")
		c := mustClient()
		result, err := c.Delete(fmt.Sprintf("/orders/%s", orderID), nil)
		dieAPI(err)
		output.JSON(result)
	case "cancel-all":
		accountID := requireArg(rest, 0, "account ID")
		c := mustClient()
		path := fmt.Sprintf("/accounts/%s/orders", accountID)
		if p := queryParams(rest[1:]); len(p) > 0 {
			path += "?" + p.Encode()
		}
		result, err := c.Delete(path, nil)
		dieAPI(err)
		output.JSON(result)
	case "modify":
		orderID := requireArg(rest, 0, "order ID")
		c := mustClient()
		body := buildFromFlags(rest[1:], "qty", "price", "stopPrice", "validity",
			"expireDate", "stopLoss", "stopLossType", "takeProfit", "takeProfitType",
			"trStopOffset", "userComment")
		result, err := c.Patch(fmt.Sprintf("/orders/%s", orderID), body)
		dieAPI(err)
		output.JSON(result)
	default:
		output.Error(fmt.Sprintf("Unknown order subcommand: %s", sub))
		os.Exit(1)
	}
}

// ── position ──────────────────────────────────────────────────────────────────

func cmdPosition(args []string) {
	if len(args) == 0 {
		output.Error("Usage: te_demo position <close|close-all|modify> ...")
		os.Exit(1)
	}
	sub, rest := args[0], args[1:]

	switch sub {
	case "close":
		posID := requireArg(rest, 0, "position ID")
		c := mustClient()
		var body interface{}
		if qty := flag(rest[1:], "--qty"); qty != "" {
			body = map[string]interface{}{"qty": autoType(qty)}
		}
		result, err := c.Delete(fmt.Sprintf("/positions/%s", posID), body)
		dieAPI(err)
		output.JSON(result)
	case "close-all":
		accountID := requireArg(rest, 0, "account ID")
		c := mustClient()
		path := fmt.Sprintf("/accounts/%s/positions", accountID)
		if p := queryParams(rest[1:]); len(p) > 0 {
			path += "?" + p.Encode()
		}
		result, err := c.Delete(path, nil)
		dieAPI(err)
		output.JSON(result)
	case "modify":
		posID := requireArg(rest, 0, "position ID")
		c := mustClient()
		body := buildFromFlags(rest[1:], "stopLoss", "takeProfit", "trailingOffset")
		result, err := c.Patch(fmt.Sprintf("/positions/%s", posID), body)
		dieAPI(err)
		output.JSON(result)
	default:
		output.Error(fmt.Sprintf("Unknown position subcommand: %s", sub))
		os.Exit(1)
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

// mustClient returns a client using the saved token
func mustClient() *api.Client {
	cfg, err := config.Load()
	die(err)
	if cfg.Token == "" {
		output.Error("Token not set — run: te_demo config --token <your_token>")
		os.Exit(1)
	}
	return api.NewClient(cfg.BaseURL, cfg.Token)
}

func die(err error) {
	if err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
}

func dieAPI(err error) {
	if err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
}

func orDash(s string) string {
	if s == "" {
		return "(not set)"
	}
	return s
}

func flag(args []string, name string) string {
	for i, a := range args {
		if a == name && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(a, name+"=") {
			return strings.TrimPrefix(a, name+"=")
		}
	}
	return ""
}

func requireArg(args []string, idx int, name string) string {
	if idx >= len(args) || strings.HasPrefix(args[idx], "--") {
		output.Error(fmt.Sprintf("Required argument missing: %s", name))
		os.Exit(1)
	}
	return args[idx]
}

// queryParams converts --key value pairs to url.Values
func queryParams(args []string) url.Values {
	params := url.Values{}
	for i := 0; i < len(args)-1; i++ {
		if strings.HasPrefix(args[i], "--") {
			key := strings.TrimPrefix(args[i], "--")
			params.Set(key, args[i+1])
			i++
		}
	}
	return params
}

func buildFromFlags(args []string, keys ...string) map[string]interface{} {
	body := map[string]interface{}{}
	for _, k := range keys {
		if v := flag(args, "--"+k); v != "" {
			body[k] = autoType(v)
		}
	}
	return body
}

func buildOrderBody(args []string) map[string]interface{} {
	keys := []string{
		"tradableInstrumentId", "side", "type", "qty", "cashOrderQty",
		"validity", "expireDate", "price", "stopPrice",
		"stopLoss", "stopLossType", "takeProfit", "takeProfitType",
		"trStopOffset", "userComment",
	}
	return buildFromFlags(args, keys...)
}

func autoType(s string) interface{} {
	switch s {
	case "true":
		return true
	case "false":
		return false
	case "null":
		return nil
	}
	var i int64
	n, _ := fmt.Sscanf(s, "%d", &i)
	if n == 1 && fmt.Sprintf("%d", i) == s {
		return i
	}
	var f float64
	m, _ := fmt.Sscanf(s, "%g", &f)
	if m == 1 && strings.ContainsAny(s, ".eE") {
		return f
	}
	return s
}
