package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
)

func main() {
	//sdd = boolean for print statements for debugging.
	// Stupidly named? Yes. I just wanted something easy to type.
	sdd := false
	activematchurl := "https://api.deadlock-api.com/v1/matches/active"

	fmt.Printf("Beginning item reccomender")
	//1. get current player
	// var playerID int := os.Getenv("userAccountName")
	var playerID int = 1395432137
	//2. see if player is in game
	sendstring := (activematchurl + "?account_id=" + strconv.Itoa(playerID))
	req, _ := http.NewRequest("GET", sendstring, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if sdd {
		// fmt.Println(res)
		fmt.Println(string(body))
	}
	//2.1 parse into structs
	// Even though users can only be in one match at once
	// we use a slice of multiple matches
	// this is so the JSON can unfurl correctly
	var currmatches []ActiveMatch
	err := json.Unmarshal(body, &currmatches)
	if err != nil {
		panic(err)
	}

	if len(currmatches) != 1 {
		// Handle the unexpected cases gracefully.
		fmt.Printf("User is not in a match.")
		return
	}

	var currmatch = currmatches[0]

	//3. see what team player is on
	teambool := 0
	for _, value := range currmatch.Players {
		if int(value.AccountID) == playerID {
			teambool = int(value.Team)
		}
	}

	print("Player is on team: %d", teambool)

	//4. see what heroes are on the enemy team
	heroIDslice := make([]int, 6, 6)
	slicePTR := 0
	for _, value := range currmatch.Players {
		if value.Team != teambool {
			heroIDslice[slicePTR] = value.HeroID
			slicePTR += 1
		}
	}

	//5. Iterate through every enemy hero. For every counter item we have, we add into the item slice.
	// IMPORTANT: Items cannot have spaces. Instead - use _ in between them.
	itemsSlice := make([]string, 0)
	antihealitems := []string{"healbane", "decay", "toxic_bullets", "spirit_burn", "inhibitor", "crippling_headshot"}
	antimovementitems := []string{"slowing_hex", "knockdown", "artic_blast", "focus_lens"}
	antifirerateitems := []string{"juggernaut", "disarming_hex", "supressor", "phantom_strike", "plated_armor", "metal_skin", "return_fire"}
	// antimeleeitems := []string{"rebuttal"}
	anticatchitems := []string{"rescue_beam", "divine_barrier"}
	antiburstitems := []string{"ethereal_shift", "spellbreaker"}
	anticarryitems := []string{"inhibitor"}
	antitankitems := []string{"toxic_bullets", "tankbuster", "mythic_slow", "scourge"}
	movementitems := []string{"warp_stone", "fleetfoot", "enduring_speed"}
	// interruptitems := []string{"knockdown"}
	// Having an editor with code folding capability is basically MANDATORY
	for _, value := range heroIDslice {
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

	//6. Now that we have the entire slice with items, let's sort it
	// so items with the same name are next to each other
	if sdd {
		fmt.Printf("Unsorted slice: %v", (itemsSlice))
		sort.Strings(itemsSlice)
		fmt.Printf("Sorted slice: %v", (itemsSlice))
	}

	//6.1 with the sorted slice, let's create a map. name is the key, value is how many times it appears.
	//
	//make the slice
	counts := make(map[string]int)
	for _, item := range itemsSlice {
		counts[item]++
	}

	//go maps are unorderd. So we map the slice into another slice (lol)
	// this time, it's got a count.
	type ItemCount struct {
		Name  string
		Count int
	}

	var sortedItems []ItemCount
	for name, count := range counts {
		sortedItems = append(sortedItems, ItemCount{Name: name, Count: count})
	}
	//sort
	sort.Slice(sortedItems, func(i, j int) bool {
		if sortedItems[i].Count != sortedItems[j].Count {
			return sortedItems[i].Count > sortedItems[j].Count
		}
		return sortedItems[i].Name < sortedItems[j].Name
	})

	fmt.Printf("\n sorted items: %v", sortedItems)

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
