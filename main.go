package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"text/template"
)

// ItemCount is used to store an item's name and its recommendation frequency.
type ItemCount struct {
	Name  string
	Count int
}

func main() {
	// sdd = boolean for print statements for debugging.
	const sdd bool = false
	pages := template.Must(template.New("").ParseGlob("static/*.html"))
	mux := http.NewServeMux()
	mux.Handle("GET /static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	//crawler detection
	//RN I dont care... but in the future I may.
	// uastring := "curl/7.54.0"
	// if crawlerdetect.IsCrawler(uastring) {
	// 	//http.Redirect(w, nil, "/lookatAtPosts", http.StatusTemporaryRedirect)
	// 	// http.Redirect(w, nil, )
	// 	// 1. create a markov babbler service
	// 	// 2. send em there
	// }
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		err := pages.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("POST /accountIDlookup", func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data from the incoming request.
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Get the value from the form data using the "name" attribute from your input.
		numfromURL := r.FormValue("accountID")

		if numfromURL == "" || numfromURL == " " {
			fmt.Fprintf(w, "Dude you sent a blank string. Fail...", nil)
			return
		}
		println("we got the account lookup. Now moving on.")

		// Step 1 & 2: Get the player's current active match from the API.
		playerID, err := strconv.Atoi(numfromURL) // Hardcoded for this example
		if err != nil {
			fmt.Fprintf(w, "Dude, you didn't send a working number. Epic fail...", nil)
			return
		}
		currmatch, err := getActiveMatchForPlayer(playerID, sdd)
		if err != nil {
			fmt.Fprintf(w, "Player isn't in a match or their profile is not public.", nil)
			return
		}
		if sdd {
			fmt.Printf("Successfully found match for player %d\n", playerID)
		}

		// Step 3: Find out which team the player is on.
		playerTeam, err := findPlayerTeam(currmatch, playerID)
		if err != nil {
			log.Fatalf("Could not find player in match data: %v", err)
		}
		fmt.Printf("Player is on team: %d\n", playerTeam)

		// Step 4: Get a list of all hero IDs on the enemy team.
		enemyHeroIDs := getEnemyHeroIDs(currmatch, playerTeam)
		if sdd {
			fmt.Printf("Found enemy hero IDs: %v\n", enemyHeroIDs)
		}

		// Step 5: Generate a list of counter-item recommendations based on enemy heroes.
		recommendedItems := generateItemRecommendations(enemyHeroIDs)
		if sdd {
			fmt.Printf("Unsorted recommendations: %v\n", recommendedItems)
		}

		// Step 6: Count and sort the items by how frequently they were recommended.
		sortedItems := countAndSortItems(recommendedItems)

		// Final Step: Display the results.
		fmt.Println("\n--- Top Recommended Items ---")
		for _, item := range sortedItems {
			fmt.Printf("Item: %-20s | Recommended: %d time(s)\n", item.Name, item.Count)
		}

		data := map[string]interface{}{
			"Items": sortedItems,
		}

		if err := pages.ExecuteTemplate(w, "recitems.html", data); err != nil {
			http.Error(w, "Failed to render recommendations", http.StatusInternalServerError)
			log.Println(err)
		}

	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// --- Helper Functions ---

// getActiveMatchForPlayer fetches the active match for a given player ID.
func getActiveMatchForPlayer(playerID int, sdd bool) (ActiveMatch, error) {
	activeMatchURL := "https://api.deadlock-api.com/v1/matches/active"
	sendstring := (activeMatchURL + "?account_id=" + strconv.Itoa(playerID))
	req, err := http.NewRequest("GET", sendstring, nil)
	if err != nil {
		return ActiveMatch{}, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ActiveMatch{}, fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ActiveMatch{}, fmt.Errorf("API returned non-200 status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ActiveMatch{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if sdd {
		fmt.Printf("API Response Body: %s\n", string(body))
	}

	// The API returns a list of matches, even if it's just one.
	var currmatches []ActiveMatch
	if err := json.Unmarshal(body, &currmatches); err != nil {
		return ActiveMatch{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if len(currmatches) == 0 {
		return ActiveMatch{}, fmt.Errorf("user is not in a match")
	}

	if len(currmatches) > 1 {
		return ActiveMatch{}, fmt.Errorf("unexpected API response: user is in multiple matches")
	}

	return currmatches[0], nil
}

// findPlayerTeam iterates through the match players to find the team of our target player.
func findPlayerTeam(match ActiveMatch, playerID int) (int, error) {
	for _, player := range match.Players {
		if player.AccountID == playerID {
			return player.Team, nil
		}
	}
	return -1, fmt.Errorf("player with ID %d not found in the match", playerID)
}

// getEnemyHeroIDs collects the hero IDs of all players not on the player's team.
func getEnemyHeroIDs(match ActiveMatch, playerTeam int) []int {
	var heroIDslice []int
	for _, player := range match.Players {
		if player.Team != playerTeam {
			heroIDslice = append(heroIDslice, player.HeroID)
		}
	}
	return heroIDslice
}

// generateItemRecommendations creates a list of suggested items based on enemy heroes.
func generateItemRecommendations(enemyHeroIDs []int) []string {
	// Item categories
	antihealitems := []string{"healbane", "decay", "toxic_bullets", "spirit_burn", "inhibitor", "crippling_headshot"}
	antimovementitems := []string{"slowing_hex", "knockdown", "artic_blast", "focus_lens"}
	antifirerateitems := []string{"juggernaut", "disarming_hex", "supressor", "phantom_strike", "plated_armor", "metal_skin", "return_fire"}
	anticatchitems := []string{"rescue_beam", "divine_barrier"}
	antiburstitems := []string{"ethereal_shift", "spellbreaker"}
	anticarryitems := []string{"inhibitor"}
	antitankitems := []string{"toxic_bullets", "tankbuster", "mythic_slow", "scourge"}
	movementitems := []string{"warp_stone", "fleetfoot", "enduring_speed"}

	var itemsSlice []string

	for _, value := range enemyHeroIDs {
		//5.1 case statement
		switch value {

		case 1:
			//infernus
			itemstoadd := []string{"debuff_remover"}
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, antifirerateitems...)
			itemsSlice = append(itemsSlice, anticarryitems...)
		case 2:
			//Seven
			itemstoadd := []string{"debuff_remover", "knockdown"}
			itemsSlice = append(itemsSlice, itemstoadd...)
		case 3:
			//Vindicta
			itemstoadd := []string{"phantom_strike"}
			itemsSlice = append(itemsSlice, antifirerateitems...)
			itemsSlice = append(itemsSlice, itemstoadd...)

		case 4:
			// Lady Geist
			itemstoadd := []string{"silence_wave"}
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, antihealitems...)
			itemsSlice = append(itemsSlice, antiburstitems...)
		case 6:
			// Abrams
			itemstoadd := []string{"warp_stone", "inhibitor"}
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, antitankitems...)
			itemsSlice = append(itemsSlice, movementitems...)
		case 7:
			// Wraith
			itemstoadd := []string{"ethereal shift"}
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, antifirerateitems...)
			itemsSlice = append(itemsSlice, anticatchitems...)
		case 8:
			// McGinnis
			itemstoadd := []string{"cold_front"}
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, antihealitems...)
		case 10:
			// Paradox
			itemstoadd := []string{ /* TODO: Add items for Paradox */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, anticatchitems...)
		case 11:
			// Dynamo
			itemstoadd := []string{"knockdown"}
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, antihealitems...)
		case 12:
			// Kelvin
			itemstoadd := []string{ /* TODO: Add items for Kelvin */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, antihealitems...)
		case 13:
			// Haze
			itemstoadd := []string{ /* TODO: Add items for Haze */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, antifirerateitems...)
			itemsSlice = append(itemsSlice, anticarryitems...)
			itemsSlice = append(itemsSlice, anticatchitems...)
		case 14:
			// Holliday
			itemstoadd := []string{ /* TODO: Add items for Holliday */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, anticatchitems...)
		case 15:
			// Bebop
			itemstoadd := []string{ /* TODO: Add items for Bebop */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, anticatchitems...)
		case 16:
			// Calico
			itemstoadd := []string{ /* TODO: Add items for Calico */ }
			itemsSlice = append(itemsSlice, itemstoadd...)

			itemsSlice = append(itemsSlice, anticatchitems...)
		case 17:
			// Grey Talon
			itemstoadd := []string{ /* TODO: Add items for Grey Talon */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, antiburstitems...)
		case 18:
			// Mo & Krill
			itemstoadd := []string{ /* TODO: Add items for Mo & Krill */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
		case 19:
			// Shiv
			itemstoadd := []string{"counterspell"}
			itemsSlice = append(itemsSlice, itemstoadd...)
			itemsSlice = append(itemsSlice, antitankitems...)
			itemsSlice = append(itemsSlice, antimovementitems...)
		case 20:
			// Ivy
			itemstoadd := []string{ /* TODO: Add items for Ivy */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
		case 25:
			// Warden
			itemstoadd := []string{ /* TODO: Add items for Warden */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
		case 27:
			// Yamato
			itemstoadd := []string{ /* TODO: Add items for Yamato */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
		case 31:
			// Lash
			itemstoadd := []string{ /* TODO: Add items for Lash */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
		case 35:
			// Viscous
			itemstoadd := []string{ /* TODO: Add items for Viscous */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
		case 50:
			// Pocket
			itemstoadd := []string{ /* TODO: Add items for Pocket */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
		case 52:
			// Mirage
			itemstoadd := []string{ /* TODO: Add items for Mirage */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
		case 58:
			// Vyper
			itemstoadd := []string{ /* TODO: Add items for Vyper */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
		case 60:
			// Sinclair
			itemstoadd := []string{ /* TODO: Add items for Sinclair */ }
			itemsSlice = append(itemsSlice, itemstoadd...)
		default:
			// This handles any hero IDs not in the list (e.g., new heroes).
			// You might want to log this to know when your list needs an update.
			fmt.Printf("Warning: Unhandled hero ID %d found.\n", value)
		}
	}
	return itemsSlice
}

// countAndSortItems takes a slice of item strings, counts occurrences, and sorts them.
func countAndSortItems(itemsSlice []string) []ItemCount {
	// Create a map to count how many times each item appears.
	counts := make(map[string]int)
	for _, item := range itemsSlice {
		counts[item]++
	}

	// Convert the map to a slice of ItemCount structs for sorting.
	var sortedItems []ItemCount
	for name, count := range counts {
		sortedItems = append(sortedItems, ItemCount{Name: name, Count: count})
	}

	// Sort the slice. Primary sort is by Count (descending), secondary is by Name (ascending).
	sort.Slice(sortedItems, func(i, j int) bool {
		if sortedItems[i].Count != sortedItems[j].Count {
			return sortedItems[i].Count > sortedItems[j].Count
		}
		return sortedItems[i].Name < sortedItems[j].Name
	})

	return sortedItems
}

// Types and objects.
type ActiveMatchPlayer struct {
	Abandoned  *bool  `json:"abandoned"`
	AccountID  int    `json:"account_id"`
	HeroID     int    `json:"hero_id"`
	Team       int    `json:"team"`
	TeamParsed string `json:"team_parsed"` // Enum: "Team0", "Team1", "Spectator"
}
type ActiveMatch struct {
	CompatVersion       *int                `json:"compat_version"`
	DurationS           *int                `json:"duration_s"`
	GameMode            *int                `json:"game_mode"`
	GameModeParsed      *string             `json:"game_mode_parsed"` // Enum: "KECitadelGameModeNormal", ...
	GameModeVersion     *int                `json:"game_mode_version"`
	LobbyID             *uint64             `json:"lobby_id"`
	MatchID             *int                `json:"match_id"`
	MatchMode           *int                `json:"match_mode"`
	MatchModeParsed     *string             `json:"match_mode_parsed"` // Enum: "Unranked", "Ranked", ...
	MatchScore          *int                `json:"match_score"`
	NetWorthTeam0       *int                `json:"net_worth_team_0"`
	NetWorthTeam1       *int                `json:"net_worth_team_1"`
	ObjectivesMaskTeam0 *int                `json:"objectives_mask_team0"`
	ObjectivesMaskTeam1 *int                `json:"objectives_mask_team1"`
	OpenSpectatorSlots  *int                `json:"open_spectator_slots"`
	Players             []ActiveMatchPlayer `json:"players"` // This is a required field
	RegionMode          *int                `json:"region_mode"`
	RegionModeParsed    *string             `json:"region_mode_parsed"` // Enum: "Europe", "SeAsia", ...
	Spectators          *int                `json:"spectators"`
	StartTime           *int                `json:"start_time"`
	WinningTeam         *int                `json:"winning_team"`
	WinningTeamParsed   *string             `json:"winning_team_parsed"` // Enum: "Team0", "Team1", "Spectator"
}
