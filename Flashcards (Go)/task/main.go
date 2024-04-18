package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	cardsList      []card
	logs           string
	importFileName = flag.String("import_from", "", "(Optional) Filename to import cards from")
	exportFileName = flag.String("export_to", "", "(Optional) Filename to export cards to")
)

type card struct {
	term       string
	definition string
	mistakes   int
}

func main() {
	flag.Parse()
	if *importFileName != "" {
		readCardsFromFile(*importFileName)
	}
	for {
		printOutput("Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):")
		switch input := readLine(); input {
		case "add":
			printOutput("The card:")
			term := readLine()
			for _, ok := duplicateTermsChecker(&cardsList, term); ok; _, ok = duplicateTermsChecker(&cardsList, term) {
				printOutput(fmt.Sprintf("The term \"%s\" already exists. Try again:\n", term))
				term = readLine()
			}
			printOutput("The definition of the card:")
			definition := readLine()
			for _, ok := duplicateDefinitionsChecker(&cardsList, definition); ok; _, ok = duplicateDefinitionsChecker(&cardsList, definition) {
				printOutput(fmt.Sprintf("The definition \"%s\" already exists. Try again:\n", definition))
				definition = readLine()
			}
			cardsList = append(cardsList, card{term: term, definition: definition})
			printOutput(fmt.Sprintf("\"The pair (\"%s\":\"%s\") has been added.\"", term, definition))
		case "remove":
			printOutput("Which card?")
			term := readLine()
			removeCard(term)
		case "import":
			printOutput("File name:")
			filePath := readLine()
			readCardsFromFile(filePath)
		case "export":
			printOutput("File name:")
			filePath := readLine()
			writeCardsToFile(filePath)
		case "ask":
			printOutput("How many times to ask?")
			count, err := strconv.Atoi(readLine())
			if err != nil {
				log.Fatal(err)
			}
			for i := 0; i < count; i++ {
				index := i
				if i >= len(cardsList) {
					index = i % len(cardsList)
				}
				printOutput(fmt.Sprintf("Print the definition of \"%s\"\n", cardsList[index].term))
				answer := readLine()
				if cardsList[index].definition == answer {
					printOutput("Correct!")
				} else {
					if key, ok := duplicateDefinitionsChecker(&cardsList, answer); ok {
						printOutput(fmt.Sprintf("Wrong. The right answer is \"%s\", but your definition is correct for \"%s\".\n", cardsList[index].definition, cardsList[key].term))
						cardsList[index].mistakes = cardsList[index].mistakes + 1
						continue
					}
					printOutput(fmt.Sprintf("Wrong. The right answer is \"%s\".\n", cardsList[index].definition))
					cardsList[index].mistakes = cardsList[index].mistakes + 1
				}
			}
		case "log":
			printOutput("File name:")
			filePath := readLine()
			err := os.WriteFile(filePath, []byte(logs), 0644)
			if err != nil {
				log.Fatal(err)
			}
			printOutput("The log has been saved.")
		case "reset stats":
			for k, _ := range cardsList {
				cardsList[k].mistakes = 0
			}
			printOutput("Card statistics have been reset.")
		case "hardest card":
			mistakes := checkMistakes()
			if len(mistakes) == 0 {
				printOutput("There are no cards with errors.")
			}
			if len(mistakes) == 1 {
				printOutput(fmt.Sprintf("The hardest card is \"%s\". You have %s errors answering it.", cardsList[mistakes[0]].term, cardsList[mistakes[0]].mistakes))
			}
			if len(mistakes) > 1 {
				termsList := ""
				mistakesSum := 0
				for _, v := range mistakes {
					mistakesSum = mistakesSum + cardsList[v].mistakes
					termsList = termsList + "\"" + cardsList[v].term + "\"" + " "
				}
				termsList = termsList + fmt.Sprintf("You have %s errors answering it.", strconv.Itoa(mistakesSum))
				printOutput(fmt.Sprintf("The hardest cards are %v", termsList))
			}
		case "exit":
			if *exportFileName != "" {
				writeCardsToFile(*exportFileName)
			}
			printOutput(fmt.Sprintf("Bye bye!"))
			os.Exit(0)
		}
	}
}

func checkMistakes() []int {
	biggestNumber := 0
	var biggestMistakesIndexes []int
	for _, v := range cardsList {
		if v.mistakes > biggestNumber {
			biggestNumber = v.mistakes
		}
	}
	if biggestNumber == 0 {
		return biggestMistakesIndexes
	}
	for k, v := range cardsList {
		if v.mistakes == biggestNumber {
			biggestMistakesIndexes = append(biggestMistakesIndexes, k)
		}
	}
	return biggestMistakesIndexes
}

func printOutput(str string) {
	fmt.Println(str)
	logs = logs + fmt.Sprintln(str)
}

func writeCardsToFile(filename string) {
	data := ""
	for _, v := range cardsList {
		data = data + v.term + "=" + v.definition + "=" + strconv.Itoa(v.mistakes) + "\n"
	}
	err := os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		log.Fatal(err)
	}
	printOutput(fmt.Sprintf("%v cards have been saved.", len(cardsList)))
}

func readCardsFromFile(filename string) {
	lines := 0
	cardsList = []card{}
	reader, err := os.ReadFile(filename)
	if err != nil {
		printOutput("File not found.")
	}
	str := strings.Split(string(reader), "\n")
	if len(str) == 0 {
		return
	}
	for _, v := range str[:len(str)-1] {
		mistakes, _ := strconv.Atoi(strings.Split(v, "=")[2])
		cardsList = append(cardsList, card{term: strings.Split(v, "=")[0], definition: strings.Split(v, "=")[1], mistakes: mistakes})
		lines++
	}
	printOutput(fmt.Sprintf("%v cards have been loaded.", lines))
}

func removeCard(term string) {
	for k, v := range cardsList {
		if v.term == term {
			cardsList = append(cardsList[:k], cardsList[k+1:]...)
			printOutput("The card has been removed.")
		}
	}
	printOutput(fmt.Sprintf("Can't remove \"%s\": there is no such card.", term))
}

func duplicateTermsChecker(cards *[]card, valueToCheck string) (int, bool) {
	for k, v := range *cards {
		if v.term == valueToCheck {
			return k, true
		}
	}
	return -1, false
}

func duplicateDefinitionsChecker(cards *[]card, valueToCheck string) (int, bool) {
	for k, v := range *cards {
		if v.definition == valueToCheck {
			return k, true
		}
	}
	return -1, false
}

func readLine() (line string) {
	reader := bufio.NewReader(os.Stdin)
	var err error
	line, err = reader.ReadString('\n')
	if err != nil {
		printOutput(fmt.Sprintln("Error reading input:", err))
		os.Exit(1)
	}
	line = strings.TrimSpace(line)
	logs = logs + fmt.Sprintln(line)
	return line
}
