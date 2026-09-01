package seed

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dis70rt/flowback/internal/config"
	"github.com/redis/go-redis/v9"
)

type menuState int

const (
	stateMainMenu menuState = iota
	stateDBMenu
	stateCardMenu
	stateSubMenu
)

type CustomerData struct {
	RzpID string
	Name  string
	Email string
	Tier  string
}

type tuiModel struct {
	state     menuState
	cursor    int
	selected  string
	customers []CustomerData
}

func fetchCustomers(db *sql.DB) []CustomerData {
	rows, err := db.Query("SELECT razorpay_customer_id, name, email, value_tier FROM customers ORDER BY created_at ASC")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var customers []CustomerData
	for rows.Next() {
		var c CustomerData
		var rzpID, name, email, tier sql.NullString
		if err := rows.Scan(&rzpID, &name, &email, &tier); err == nil {
			c.RzpID = rzpID.String
			c.Name = name.String
			c.Email = email.String
			c.Tier = tier.String
			customers = append(customers, c)
		}
	}
	return customers
}

func initialModel(db *sql.DB) tuiModel {
	return tuiModel{
		state:     stateMainMenu,
		customers: fetchCustomers(db),
	}
}

func (m tuiModel) Init() tea.Cmd { return nil }

func (m tuiModel) currentChoices() []string {
	switch m.state {
	case stateMainMenu:
		return []string{
			"System & Database Management",
			"Simulate Card Failures",
			"Simulate Subscription Failures",
			"Exit",
		}
	case stateDBMenu:
		return []string{
			"Seed Database with 6 Archetypes",
			"Flush Redis (Clear Ghost Tasks)",
			"Back to Main Menu",
		}
	case stateCardMenu, stateSubMenu:
		if len(m.customers) == 0 {
			return []string{"(No customers found in DB. Go back and Seed first.)", "Back to Main Menu"}
		}
		var choices []string
		for _, c := range m.customers {
			choices = append(choices, fmt.Sprintf("%s (%s) - %s", c.Name, c.Tier, c.Email))
		}
		choices = append(choices, "Back to Main Menu")
		return choices
	}
	return []string{}
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	choices := m.currentChoices()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.selected = "Exit"
			return m, tea.Quit
		case "esc":
			m.state = stateMainMenu
			m.cursor = 0
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			selection := choices[m.cursor]

			if selection == "Back to Main Menu" {
				m.state = stateMainMenu
				m.cursor = 0
				return m, nil
			}

			if selection == "(No customers found in DB. Go back and Seed first.)" {
				return m, nil
			}

			if m.state == stateMainMenu {
				switch selection {
				case "System & Database Management":
					m.state = stateDBMenu
				case "Simulate Card Failures":
					m.state = stateCardMenu
				case "Simulate Subscription Failures":
					m.state = stateSubMenu
				case "Exit":
					m.selected = "Exit"
					return m, tea.Quit
				}
				m.cursor = 0
				return m, nil
			}

			m.selected = selection
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m tuiModel) View() string {
	var title string
	switch m.state {
	case stateMainMenu:
		title = "--- Flowback Developer TUI ---"
	case stateDBMenu:
		title = "--- System & Database Management ---"
	case stateCardMenu:
		title = "--- Card Payment Failures ---"
	case stateSubMenu:
		title = "--- Subscription Failures ---"
	}

	s := fmt.Sprintf("\n%s\n\n", title)

	for i, choice := range m.currentChoices() {
		cursor := "  "
		if m.cursor == i {
			cursor = "-> "
		}
		s += fmt.Sprintf("%s%s\n", cursor, choice)
	}
	s += "\n[Press esc/q to quit or return]\n"
	return s
}

// RunTUI is the main entry point for the interactive terminal dashboard
func RunTUI(db *sql.DB, cfg *config.Config) {
	stateMemo := stateMainMenu

	for {
		m := initialModel(db) 
		m.state = stateMemo

		p := tea.NewProgram(m, tea.WithAltScreen())
		result, err := p.Run()
		if err != nil {
			log.Fatalf("Error running TUI: %v", err)
		}

		finalModel := result.(tuiModel)
		if finalModel.selected == "" || finalModel.selected == "Exit" {
			break
		}

		stateMemo = finalModel.state

		if finalModel.selected == "Seed Database with 6 Archetypes" {
			fmt.Println("\n[*] Running Database Seeder...")
			RunDBSeeder(db, cfg.RazorpayKeyID, cfg.RazorpayKeySecret)

			fmt.Println("\n[Press Enter to return to menu...]")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}

		if finalModel.selected == "Flush Redis (Clear Ghost Tasks)" {
			fmt.Println("\n[*] Connecting to Redis to flush all queues...")
			rdb := redis.NewClient(&redis.Options{
				Addr: cfg.RedisAddr,
			})
			
			err := rdb.FlushAll(context.Background()).Err()
			if err != nil {
				fmt.Printf("[!] Failed to flush redis: %v\n", err)
			} else {
				fmt.Println("[+] SUCCESS: Flushed the entire Redis database! All ghost tasks have been deleted.")
			}
			
			fmt.Println("\n[Press Enter to return to menu...]")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}

		engine := New(cfg.RazorpayKeyID, cfg.RazorpayKeySecret)
		var task Task

		var targetCustomer CustomerData
		if finalModel.cursor < len(finalModel.customers) {
			targetCustomer = finalModel.customers[finalModel.cursor]
		} else {
			continue 
		}

		// Set dynamic amounts based on value tier
		amt := 50000 
		if targetCustomer.Name == "Enterprise CEO" {
			amt = 95000000
		} else if targetCustomer.Tier == "HIGH" {
			amt = 999000
		} else if targetCustomer.Tier == "LOW" {
			amt = 15000
		}

		if finalModel.state == stateCardMenu {
			task = &FailedCardTask{
				CustomerID:  targetCustomer.RzpID,
				Email:       targetCustomer.Email,
				AmountPaise: amt,
				CardNumber:  CardErrorInsufficientFund,
			}
		} else if finalModel.state == stateSubMenu {
			task = &FailedSubscriptionTask{
				CustomerID:  targetCustomer.RzpID,
				Email:       targetCustomer.Email,
				AmountPaise: amt,
			}
		}

		if task != nil {
			fmt.Printf("\n[*] Executing task for %s (RZP ID: %s)\n\n", targetCustomer.Name, targetCustomer.RzpID)
			engine.Execute(task)

			fmt.Println("\n[Press Enter to return to menu...]")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}
}
