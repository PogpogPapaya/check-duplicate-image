package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(6)
	go findDuplication("Datasets/train/unripe", true, wg)
	go findDuplication("Datasets/train/ripe", true, wg)
	go findDuplication("Datasets/train/medium", true, wg)
	go findDuplication("Datasets/validation/unripe", true, wg)
	go findDuplication("Datasets/validation/ripe", true, wg)
	go findDuplication("Datasets/validation/medium", true, wg)
}



func findDuplication(dirPath string, deleteDuplication bool, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Looking for duplication on %s\n", dirPath)
	dirFiles, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatalln("unable to read directory")
	}

	filesPath := make([]string, len(dirFiles))
	for i, file := range dirFiles {
		filePath := fmt.Sprintf("%s/%s", dirPath, file.Name())
		filesPath[i] = filePath
	}

	sums := make(map[string][]string, len(filesPath))

	for i, filePath := range filesPath {
		log.Printf("Checking file %d out of %d\n", i, len(filesPath))

		// Calculate check sum
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalln("unable to open file " + filePath, err)
		}

		h := md5.New()
		if _, err := io.Copy(h, file); err != nil {
			log.Fatalln("unable copy file data to hash", err, filePath)
		}

		sum := h.Sum(nil)
		hexSum := fmt.Sprintf("%x", sum)

		if err := file.Close(); err != nil {
			log.Fatalln("unable to close file", err, filePath)
		}

		if len(sums[hexSum]) != 0 {
			sums[hexSum] = append(sums[hexSum], filePath)
		} else {
			sums[hexSum] = []string{filePath}
		}
	}

	log.Println("Result")
	dupCount := 0
	disCount := 0
	delCount := 0
	for _, sum := range sums {
		count := len(sum)
		if count >= 1 {
			disCount++
		}
		if count > 1 {
			log.Printf("%s has %d duplication\n", sum[0], count-1)
			for i := 1; i<count; i++ {
				log.Printf("\t%s\n", sum[i])
				if deleteDuplication {
					if err := os.Remove(sum[i]); err != nil {
						log.Fatalln("unable to remove file")
					}
					delCount++
				}
			}
			dupCount += len(sum) - 1
		}
	}
	log.Println("Total Files:", len(dirFiles))
	log.Println("Distinct Files:", disCount)
	log.Println("Duplication Files:", dupCount)
	log.Println("Deleted:", delCount)
}