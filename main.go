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
  "math/rand"
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
  pokedex map[string]pokemonJSON
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
	EncounterMethodRates any	
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
		VersionDetails any 	
  } `json:"pokemon_encounters"`
}


type pokemonJSON struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []any  `json:"past_abilities"`
	PastTypes     []any  `json:"past_types"`
	Species       struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
  Sprites any
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
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
      callback: commandExplore,
    },
    "catch" : {
      name: "catch",
      description: "Catch Pokemon",
      callback: commandCatch,
    },
    "inspect" : {
      name: "inspect",
      description: "Inspect Pokemon",
      callback: commandInspect,
    },
    "pokedex" : {
      name: "pokedex",
      description: "Pokedex of all caught Pokemon",
      callback: commandPokedex,
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

func commandExplore(c *config, params ...string) error{
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

func commandCatch(c *config, params ...string) error{
  if params == nil{
    return errors.New("Please pass in a valid pokemon to catch")
  }
  if c.pokedex == nil{
    c.pokedex = make(map[string]pokemonJSON)
  }
  pokemon_name := strings.Join(params[:], " ")
  _, ok := c.pokedex[pokemon_name]
  if !ok{
      fmt.Println(fmt.Sprintf("Throwing a Pokeball at %s...", pokemon_name))
      catch_url := "https://pokeapi.co/api/v2/pokemon/" + pokemon_name
      _, ok := c.Cache.Get(catch_url)
      if !ok{
        res, err := http.Get(catch_url)
        if err != nil{
          return errors.New("Problem hitting the API")
        }
        defer res.Body.Close()
        body, _ := io.ReadAll(res.Body)
        c.Cache.Add(catch_url, body)
      }
      var pokemonjsondata pokemonJSON
      stream_obj, _ := c.Cache.Get(catch_url)
      json.Unmarshal(stream_obj, &pokemonjsondata)
      base_exp := pokemonjsondata.BaseExperience
      chance := rand.Intn(base_exp)
      if chance > (base_exp / 2){
        c.pokedex[pokemon_name] = pokemonjsondata
        fmt.Println(fmt.Sprintf("%s was caught!", pokemon_name))
        fmt.Println("You may now inspect it with the inspect command.")
      }else{
        fmt.Println(fmt.Sprintf("%s escaped!", pokemon_name))
      }
      return nil
  }else{
    fmt.Println(fmt.Sprintf("%s is already in pokedex", pokemon_name))
    return nil
  }

}


func commandInspect(c *config, params ...string) error{
  if params == nil{
    return errors.New("Please pass in a valid pokemon to inspect")
  }
  if c.pokedex == nil{
    c.pokedex = make(map[string]pokemonJSON)
  }
  pokemon_name := strings.Join(params[:], " ")
  val, ok := c.pokedex[pokemon_name]
  if ok{

    fmt.Println(fmt.Sprintf("Name: %s",val.Name))
    fmt.Println(fmt.Sprintf("Height: %v",val.Height))
    fmt.Println(fmt.Sprintf("Weight: %v",val.Weight))
    fmt.Println("Stats:")
    for _, stat := range val.Stats{
      fmt.Println(fmt.Sprintf("    -%s: %v",stat.Stat.Name, stat.BaseStat))
    }
    fmt.Println("Types:")
    for _, poketype := range val.Types{
      fmt.Println(fmt.Sprintf("    -%s", poketype.Type.Name))
    }
    fmt.Println()
    return nil
  }else{
    fmt.Println(fmt.Sprintf("you have not caught %s", pokemon_name))
    return nil
  }

}


func commandPokedex(c *config, params ...string) error{
  if params != nil{
    return errors.New("pokedex alone is sufficient to view the pokedex")
  }
  if c.pokedex == nil{
    return errors.New("pokedex is empty...\ncatch some pokemon to inspect")
  }
  fmt.Println("Your Pokedex:")
  for _, pokemon := range c.pokedex{
    fmt.Println(fmt.Sprintf("    - %s", pokemon.Name))
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
