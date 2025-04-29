package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	for {

		file, err := os.Create("../nifi_data/input.csv")
		if err != nil {
			panic(err)
		}

		writer := csv.NewWriter(file)
		header := []string{"name", "favorite_number", "favorite_color"}
		err = writer.Write(header)
		if err != nil {
			panic(err)
		}

		log.Println("Updating file with:")

		for range 5 {
			number := strconv.Itoa(rand.Intn(100))

			row := []string{
				"Jane Smith" + number,
				number,
				"color" + number,
			}

			err = writer.Write(row)
			if err != nil {
				panic(err)
			}

			log.Printf("%s,%s,%s", row[0], row[1], row[2])
		}

		writer.Flush()
		file.Close()

		time.Sleep(5 * time.Second)
	}
}
