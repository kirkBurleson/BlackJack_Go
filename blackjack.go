package main
import ("fmt";"math/rand";"time";"strconv";"os")

const decks = 1
const suits = 4
var shoe [52 * decks]int
var cut_spot = 52 * decks - 10
var current_spot = 0
var seed = time.Now().UTC().UnixNano()
var rng = rand.New(rand.NewSource(seed))
var dealer = [12]int {}
var player = [4][12]int { {}, {}, {}, {} }
var bets = [4]float64 {}
var bet = 0.0
var money = 50.0
var show_one = true
var doubled_down = false
var no_blackjack = [12]int {}
var hand_total = 0

func main() {
	initialize()

	for {
		fmt.Print("$",money," Bet: ")
		fmt.Scan(&bet)
		if bet == -1 {
			break
		} else if bet < 1 || bet > money {
			continue
		}

		fmt.Println("****************")

		// initial deal
		player[0][0] = next_card()
		player[0][1] = next_card()
		dealer[0] = next_card()
		dealer[1] = next_card()

		fmt.Println("Dealer:", hand_to_string(dealer))
		show_one = false
				
		// players turn
		for i := 0; i < count_hands(player); i++ {
			bets[i] = bet

			if player[i][0] == 0 { break }
			hand_total = get_total(&player[i])

choice:
			if doubled_down == true {
				doubled_down = false
				continue
			}

			// show old hands
			for j := i; j > 0; j-- {
				show_hand(player[j-1], get_total(&player[j-1]), j-1)
			}

			show_hand(player[i], hand_total, i)

			if hand_total == 0 { continue }		

			fmt.Println()
			fmt.Print("(H)it (S)tand (D)ouble down s(P)lit (Q)uit:")

			choice := ""
			fmt.Scan(&choice)			
			switch choice {

			case "s":
				fallthrough
			case "S":
				bets[i] = bet
				continue

			case "h":
				fallthrough
			case "H":
				player[i][next_card_position(player[i])] = next_card()
				hand_total = get_total(&player[i])
				goto choice

			case "d":
				fallthrough
			case "D":
				if bet * 2 <= money {
					if hand_total == 9 || hand_total == 10 || hand_total == 11 {
						if count_cards(player[i]) == 2 {
							bets[i] = bet * 2
							
							player[i][next_card_position(player[i])] = next_card()
							doubled_down = true
							hand_total = get_total(&player[i])
						} else {
							fmt.Println("You can only double down on the first 2 cards dealt to you.")
						}
					} else {
						fmt.Println("You can only double down on 9, 10, or 11.")
					}
				} else {
					fmt.Println("You don't have enough money to cover the bet.")
				}

				goto choice

			case "p":
				fallthrough
			case "P":
				if bet * 2 <= money {
					if count_hands(player) < 4 {
						if count_cards(player[i]) == 2 {
							if is_aces(player[i]) == true {
								// make sure aces are 14 instead of 1
								player[i][0] = 14
								player[i][1] = 14
							}

							if player[i][0] == player[i][1] {
								player[i+1][0] = player[i][1]
								player[i+1][1] = next_card()
								
								player[i][1] = next_card()

								// split aces can't make blackjack
								no_blackjack[i] = 1
								no_blackjack[i+1] = 1

								bets[i+1] = bet

								hand_total = get_total(&player[i])
							} else {
								fmt.Println("You can only split 2 of a kind.")
							}
						} else {
							fmt.Println("You can only split the first 2 cards dealt to you.")
						}
					} else {
						fmt.Println("You can only have 4 hands at a time.")
					}
				} else {
					fmt.Println("You don't have enough money to cover the bet.")
				}

				goto choice

			case "q":
				fallthrough
			case "Q":
				fmt.Println("Bye.")
				os.Exit(0)

			default:
				goto choice
			}
		}
		
		// dealer's turn
		fmt.Println()
		for {
			total := get_total(&dealer)
			if total == 0 {
				fmt.Println("Dealer:", hand_to_string(dealer), "-->", "BUST")
				break
			} else if total == 21 && count_cards(dealer) == 2 {
				fmt.Println("Dealer:", hand_to_string(dealer), "-->", "BLACKJACK")
				break
			} else {
				fmt.Println("Dealer:", hand_to_string(dealer), "-->", total)
				
				if 	total < 17 || has_soft_17(dealer) {
					dealer[next_card_position(dealer)] = next_card()
				} else {
					fmt.Println("Dealer stands at:", total)
					break
				}
			}
		}

		// collect winnings
		winnings := 0.0
		dealer_score := get_total(&dealer)

		for i, hand := range player {
			bet = bets[i]

			// no more hands
			if hand[0] == 0 { break }

			player_score := get_total(&hand)
			
			// blackjack
			if no_blackjack[i] == 0 {
				if has_blackjack(hand) == true && has_blackjack(dealer) == false { winnings += bet + (bet / 2); continue }
				if has_blackjack(hand) == false && has_blackjack(dealer) == true { winnings -= bet; continue }
				if has_blackjack(hand) == true && has_blackjack(dealer) == true { continue }
			}

			// bust
			if player_score == 0 && dealer_score != 0 { winnings -= bet; continue }
			if dealer_score == 0 && player_score != 0 { winnings += bet; continue }
			if player_score == 0 && dealer_score == 0 { continue }

			// highest
			if player_score > dealer_score { winnings += bet; continue }
			if player_score < dealer_score { winnings -= bet; continue }
			if player_score == dealer_score { continue }

			fmt.Println("*** SCORING ERROR: UNKNOWN LOGIC ***")
			
		}
		
		money += winnings

		fmt.Println("---------------")
		if winnings == 0 {
			fmt.Println("Push")
		} else if winnings > 0 {
			fmt.Println("Winner")
		} else {
			fmt.Println("Loser")
		}

		if money < 1 {
			fmt.Println("\nYou're broke. Bye.")
			os.Exit(0)
		}

		clear_state()
	}
}

func is_aces(hand [12]int) bool {
	return hand[0] == 1 || hand[0] == 14 && hand[1] == 1 || hand[1] == 14
}

func clear_state() {
	show_one = true
	doubled_down = false
	clear(&dealer)
	no_blackjack[0] = 0
	no_blackjack[1] = 0
	no_blackjack[2] = 0
	no_blackjack[3] = 0
	bets[0] = 0.0
	bets[1] = 0.0
	bets[2] = 0.0
	bets[3] = 0.0

	for i := 0; i < len(player); i++ {
		clear(&player[i])
	}
}

func show_hand(hand [12]int, hand_total int, i int) {
	if hand_total == 0 {
		fmt.Println("Hand " + strconv.Itoa(i+1) + ":", hand_to_string(player[i]), "-->", "BUST")
	} else if hand_total == 21  && count_cards(player[i]) == 2 && no_blackjack[i] == 0 {
		fmt.Println("Hand " + strconv.Itoa(i+1) + ":", hand_to_string(player[i]), "-->", "BLACKJACK")
	} else {
		fmt.Println("Hand " + strconv.Itoa(i+1) + ":", hand_to_string(player[i]), "-->", hand_total)
	}
}

func count_cards(hand [12]int) int {
	total := 0
	for _, card := range hand {
		if card != 0 { total += 1 }
	}
	return total
}

func has_blackjack(hand [12]int) bool {
	return get_total(&hand) == 21 && next_card_position(hand) == 2
}

func clear(hand *[12]int) {
	for i := 0; i < len(hand); i++ {
		(*hand)[i] = 0
	}
}

func has_soft_17(hand [12]int) bool {
	if next_card_position(hand) == 2 && (hand[0] == 14 || hand[1] == 14) && (hand[0] == 6 || hand[1] == 6) {
		return true
	}

	return false
}

func next_card_position(hand [12]int) int {
	pos := 0
	for _, card := range hand {
		if card != 0 {
			pos += 1
		}
	}

	return pos
}

func count_hands_index(hands [4][12]int) int {
	total := -1
	for _, cards := range hands {
		if cards[0] != 0 { total += 1 }
	}
	return total
}

func count_hands(hands [4][12]int) int {
	total := 0
	for _, cards := range hands {
		if cards[0] != 0 { total += 1 }
	}
	return total
}

func hand_to_string(hand [12]int) string{
	text := ""
	for i, card := range hand {
		if card == 0 { break }

		if card > 0 {
			switch card {
			case 1:
				fallthrough
			case 14:
				text += "A"

			case 13:
				text += "K"

			case 12:
				text += "Q"

			case 11:
				text += "J"

			default:
				text += strconv.Itoa(card)
			}

			// dealer shows 1 card at initial deal
			if show_one == true { return text }

			// add space if not last card
			if i < len(hand) - 2 {
				if hand[i+1] > 0 {
					text += " "
				}
			}
		}
	}
	return text
}

func index_of(hand [12]int, target int) (int, bool) {
	for i, card := range hand {
		if card == target {
			return i, true
		}
	}

	return -1, false
}

func get_total(hand *[12]int) int {
	total := 0
	for _, card := range hand {
		if card == 14 {
			total += 11
		} else if card > 9 {
			total += 10
		} else {
			total += card
		}
	}
	
	if total > 21 {
		// try using ace as a 1 instead of 11
		if change_ace(hand) == true { return get_total(hand) }

		// no aces
		return 0
	}

	return total
}

func change_ace(hand *[12]int) bool {
	pos, has_ace := index_of(*hand, 14)
	if has_ace == true {
		(*hand)[pos] = 1
		return true
	}

	return false
}

func initialize() {
	fill_shoe()
	shuffle()
}

func fill_shoe() {
	for i := 0; i < decks; i++ {
		x := i * 52
		for j := 0; j < suits; j++ {
			y := j * 13
			for k := 0; k < 13; k++ {
				shoe[x + y + k] = k + 2
			}
		}
	}
}

func random(max int)  int {	
	number := rng.Intn(max)
	return number
}

func shuffle() {
	end := 52 * decks
	for end > 1 {
		pos := random(end)
		shoe[pos], shoe[end - 1] = shoe[end - 1], shoe[pos]
		end--
	}
}

func next_card() int {
	if current_spot == cut_spot {
		shuffle()
		current_spot = 0
	}

	card := shoe[current_spot]
	current_spot++

	return card
}