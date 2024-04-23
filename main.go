package main

import (
  "fmt"
  "os"
  "bufio"
  "net/http"
"errors"
  "encoding/json"
  "github.com/Mukilan-Krishnakumar/pokedexcli/internal/pokecache"
  "time"
  "io"
  "strings"
)

type cliCommand struct{
  name string
  description string
  callback func(*config, ...string) error
}

type config struct{
  Next *string
  Previous *string
  Cache pokecache.Cache
}

type locationJSON struct{
  Count int `json:"count"`
  Next *string `json:"next"`
  Previous *string `json:"previous"`
  Results []struct{
    Name string `json:"name"`
    URL string `json:"url"`
  } `json:"results"`

}

type exploreJSON struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func getMap() map[string]cliCommand{
  getMap := map[string]cliCommand{
    "help" : {
        name:        "help",
        description: "Displays a help message",
        callback:    commandHelp,
    },
    "exit": {
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
    },
    "map" : {
      name: "map",
      description: "Map forward",
      callback: commandMap,
    },
    "mapb" : {
      name: "mapb",
      description: "Map backward",
      callback: commandMapB,
    },
    "explore" : {
      name: "explore",
      description: "Explore the region",
      callback: explore,
    },


  }
  return getMap
}

func commandHelp(c *config, params ...string) error{
  helpText := "Welcome to the Pokedex!\nUsage:\n"
  getMap := getMap()
  fmt.Println(helpText)
  for _, command := range getMap{
    fmt.Println(fmt.Sprintf("%s: %s",command.name, command.description))
  }
  fmt.Println()
  return nil
}


func commandExit(c *config, params ...string) error{
  os.Exit(0)
  return nil
}

func commandMap(c *config, params ...string) error{
  if c.Previous == nil{
    init_url := "https://pokeapi.co/api/v2/location-area/"
    _, ok := c.Cache.Get(init_url)
    if !ok{
      res, err := http.Get(init_url)
      if err != nil{
          return errors.New("Problem hitting the API")
        }
      defer res.Body.Close()
      body,_ := io.ReadAll(res.Body)
      c.Cache.Add(init_url, body)
    }
    var pokejsondata locationJSON
    stream_obj, _ := c.Cache.Get(init_url)
    json.Unmarshal(stream_obj, &pokejsondata)
    c.Previous = &init_url
    c.Next = pokejsondata.Next
    for _, location := range pokejsondata.Results{
      fmt.Println(location.Name)
    } 
    fmt.Println()
    return nil

  }
  if c.Next == nil{
    return nil
  }
  _, ok := c.Cache.Get(*c.Next)
  if !ok{
    res, err := http.Get(*c.Next)
    if err != nil{
      return errors.New("Problem hitting the API")
    }
    defer res.Body.Close()
    body,_ := io.ReadAll(res.Body)
    c.Cache.Add(*c.Next, body)
  }

  var pokejsondata locationJSON
  stream_obj, _ := c.Cache.Get(*c.Next)
  json.Unmarshal(stream_obj, &pokejsondata)

  c.Next = pokejsondata.Next
  c.Previous = pokejsondata.Previous
  for _, location := range pokejsondata.Results{
    fmt.Println(location.Name)
  }
  fmt.Println()
  return nil
}

func commandMapB(c *config, params ...string) error{
  if c.Previous == nil{
    return nil
  }
  _, ok := c.Cache.Get(*c.Previous)
  if !ok{
    res, err := http.Get(*c.Previous)
    if err != nil{
      return errors.New("Problem hitting the API")
    }
    defer res.Body.Close()
    body, _ := io.ReadAll(res.Body)
    c.Cache.Add(*c.Previous, body)
  }
  var pokejsondata locationJSON
  stream_obj, _ := c.Cache.Get(*c.Previous)
  json.Unmarshal(stream_obj, &pokejsondata)
  c.Next = pokejsondata.Next
  c.Previous = pokejsondata.Previous
  for _, location := range pokejsondata.Results{
    fmt.Println(location.Name)
  }
  fmt.Println()
  return nil
}

func explore(c *config, params ...string) error{
  if params == nil{
    return errors.New("Please pass in the location name to explore")
  }
  city_name := strings.Join(params[:], " ")
  location_url := "https://pokeapi.co/api/v2/location-area/" + city_name
  _, ok := c.Cache.Get(location_url)
  if !ok{
    res, err := http.Get(location_url)
    if err != nil{
      return errors.New("Problem hitting the API")
    }
    defer res.Body.Close()
    body, _ := io.ReadAll(res.Body)
    c.Cache.Add(location_url, body)
  }
  var explorejsondata exploreJSON
  stream_obj,_ := c.Cache.Get(location_url)
  json.Unmarshal(stream_obj, &explorejsondata)
  if len(explorejsondata.PokemonEncounters) < 1{
    return errors.New("Enter a valid location")
  }
  fmt.Println(fmt.Sprintf("Exploring %s...", city_name))

  fmt.Println("Found Pokemon:")
  for _, pokemon := range explorejsondata.PokemonEncounters{
    fmt.Println(fmt.Sprintf("- %s",pokemon.Pokemon.Name))
  }
  return nil
}

func main(){
  cache := pokecache.NewCache(15 * time.Second)
  prompt := "pokedex > "
  getMap := getMap() 
  var actualconf config
  config := &actualconf
  config.Cache = cache
  fmt.Printf(prompt)
  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    fmt.Println()
    command_slice := strings.Split(scanner.Text(), " ")
    val, ok := getMap[command_slice[0]] 
    if ok{
      if len(command_slice) > 1{
        err := val.callback(config, command_slice[1:]...)
        if err != nil{
          fmt.Println(err)
        }
      }else{
        err := val.callback(config)
        if err != nil{
          fmt.Println(err)
        }
      }
    }

    fmt.Printf(prompt)
  }
}
