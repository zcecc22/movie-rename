package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ryanbradynd05/go-tmdb"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var moviePatterns = []string{`(?i)(\d{4})`}

func strCleanupNonWord(str string) string {
	regExp := regexp.MustCompile(`\W+`)
	return strings.Trim(regExp.ReplaceAllString(str, " "), ` 	`)
}

func strCleanupSymbols(str string) string {
	regExp := regexp.MustCompile(`[/]`)
	return regExp.ReplaceAllString(str, "_")
}

func renameMovie(path string, movieName string, releaseYear string) (string, error) {
	ext := filepath.Ext(path)
	dir := filepath.Dir(path)
	newName := strCleanupSymbols(movieName + " (" + releaseYear + ")" + ext)
	newPath := filepath.Join(dir, newName)

	if _, err := os.Stat(newPath); err == nil {
		return newName, errors.New("Destination file already exists.")
	}

	return newName, os.Rename(path, newPath)
}

func MovieInfo(filename string) (string, string, error) {
	for _, curPattern := range moviePatterns {
		if matched, err := regexp.MatchString(curPattern, filename); err == nil && matched {
			movieName := ""
			releaseYear := ""
			err = nil

			regExp := regexp.MustCompile(curPattern)
			movieName = strCleanupNonWord(regExp.Split(filename, -1)[0])
			releaseYear = regExp.FindStringSubmatch(filename)[1]

			return movieName, releaseYear, err
		} else if err != nil {
			return "", "", err
		}
	}
	return "", "", errors.New("No filename pattern match found.")
}

func main() {
	flag.Parse()
	TMDb := tmdb.Init("3d6ca007d0677db4a3444067691b6b6a")

	for _, path := range flag.Args() {
		fmt.Println("\nCurrent Name: ", filepath.Base(path))
		movieName, releaseYear, err := MovieInfo(filepath.Base(path))

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Show: ", movieName, " Year: ", releaseYear)

		moviesList, err := TMDb.SearchMovie(movieName, nil)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if len(moviesList.Results) == 0 {
			fmt.Println(errors.New("No match found."))
			continue
		}

		var selectedMovie tmdb.MovieShort

		if len(moviesList.Results) >= 1 && len(moviesList.Results[0].ReleaseDate) >= 4 &&
			moviesList.Results[0].ReleaseDate[:4] == releaseYear {

			selectedMovie = moviesList.Results[0]

		} else {

			for index, movie := range moviesList.Results {
				fmt.Println("\nTMDB Match: ", index+1, " Movie: ", movie.Title, " Year: ", movie.ReleaseDate)
			}

			fmt.Println("\nSelect match:")
			i := -1
			fmt.Scanf("%d", &i)
			selectedMovie = moviesList.Results[i-1]
		}

		newName, err := renameMovie(path, selectedMovie.Title, selectedMovie.ReleaseDate[:4])

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("New Name: ", newName)
	}
}
