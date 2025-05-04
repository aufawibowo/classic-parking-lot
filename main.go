package main

import (
	"bufio"
	"fmt"
	"github.com/aufawibowo/classic-parking-lot/parking"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <commands.txt>", os.Args[0])
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lot *parking.ParkingLot
	scan := bufio.NewScanner(file)

	for scan.Scan() {
		fields := strings.Fields(scan.Text())
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "create_parking_lot":
			n, _ := strconv.Atoi(fields[1])
			lot = parking.NewParkingLot(n)

		case "park":
			slot, ok := lot.Park(fields[1])
			if ok {
				fmt.Printf("Allocated slot number: %d\n", slot)
			} else {
				fmt.Println("Sorry, parking lot is full")
			}

		case "leave":
			hours, _ := strconv.Atoi(fields[2])
			slot, fee, ok := lot.Leave(fields[1], hours)
			if !ok {
				fmt.Printf("Registration number %s not found\n", fields[1])
			} else {
				fmt.Printf("Registration number %s with Slot Number %d is free with Charge $%d\n",
					fields[1], slot, fee)
			}

		case "status":
			fmt.Println("Slot No. Registration No.")
			for _, s := range lot.Status() {
				fmt.Printf("%d %s\n", s[0], s[1])
			}
		}
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
}
