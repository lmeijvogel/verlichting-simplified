package main

import (
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

type HasID interface {
	GetID() string
}

func main() {
	env, err := godotenv.Read()

	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	ha := NewAdapter(env["API_HOST"], env["API_TOKEN"])

	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/api/scenes", func(c *gin.Context) {
		entities, err := getEntities(ha.GetScenes, getAllowedSceneIds())

		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.IndentedJSON(http.StatusOK, entities)
	})

	router.POST("/api/start_scene/:id/on", func(c *gin.Context) { startScene(c, ha) })

	router.GET("/api/switches", func(c *gin.Context) {
		entities, err := getEntities(ha.GetSwitches, getAllowedSwitchIds())
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.IndentedJSON(http.StatusOK, entities)
	})

	router.POST("/api/set_switch/:id/:new_state", func(c *gin.Context) { setSwitchOrState(c, getAllowedSwitchIds(), ha.SetSwitch) })

	router.GET("/api/states", func(c *gin.Context) {
		entities, err := getEntities(ha.GetStates, getAllowedStateIds())
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.IndentedJSON(http.StatusOK, entities)
	})

	router.POST("/api/set_state/:id/:new_state", func(c *gin.Context) { setSwitchOrState(c, getAllowedStateIds(), ha.SetState) })

	router.GET("/api/lights", func(c *gin.Context) {
		entities, err := getEntities(ha.GetLights, getAllowedLightIds())
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.IndentedJSON(http.StatusOK, entities)
	})
	router.POST("/api/set_light/:id/:new_state", func(c *gin.Context) { setSwitchOrState(c, getAllowedLightIds(), ha.SetLight) })

	router.Run(":3123")
}

func getEntities[T HasID](getter func() ([]T, error), whitelist []string) ([]T, error) {
	typedEntities, err := getter()

	if err != nil {
		return nil, err
	}

	allowedEntities := []T{}

	for i := range typedEntities {
		typedEntity := typedEntities[i]

		if slices.Contains(whitelist, typedEntity.GetID()) {
			allowedEntities = append(allowedEntities, typedEntity)
		}
	}

	sort.Slice(allowedEntities, func(i, j int) bool {
		equals := func(requiredId string) func(string) bool {
			return func(entityId string) bool {
				return entityId == requiredId
			}
		}

		positionI := slices.IndexFunc(whitelist, equals(allowedEntities[i].GetID()))
		positionJ := slices.IndexFunc(whitelist, equals(allowedEntities[j].GetID()))

		return positionI < positionJ
	})

	return allowedEntities, nil
}

func startScene(c *gin.Context, ha HomeAssistantAdapter) {
	scene_id := c.Param("id")

	if !slices.Contains(getAllowedSceneIds(), scene_id) {
		c.AbortWithStatus(403)
		return
	}

	updatedScene, err := ha.StartScene(scene_id)

	if err != nil {
		c.AbortWithError(500, err)
		log.Print(err)
		return
	}

	c.IndentedJSON(http.StatusOK, updatedScene)
}

func setSwitchOrState(c *gin.Context,
	entityWhitelist []string,
	setter func(string, AllowedStateValue) (*Switch, error)) {
	entityId := c.Param("id")
	newStateParam := c.Param("new_state")

	if !slices.Contains(entityWhitelist, entityId) {
		c.AbortWithStatus(403)
		return
	}

	var newState AllowedStateValue

	newState = Off
	if newStateParam == "on" {
		newState = On
	}

	updatedEntity, err := setter(entityId, newState)

	if err != nil {
		c.AbortWithError(500, err)
		log.Print(err)
		return
	}

	c.IndentedJSON(http.StatusOK, updatedEntity)
}

func getAllowedSceneIds() []string {
	return []string{"script.uit", "script.ochtend", "script.koken", "script.avond", "script.alles_aan"}
}

func getAllowedSwitchIds() []string {
	return []string{}
}

func getAllowedStateIds() []string {
	return []string{"input_boolean.vacation_mode", "input_boolean.auto_on_single"}
}

func getAllowedLightIds() []string {
	return []string{"switch.salontafel", "switch.eettafel", "switch.schemerlamp_voor", "switch.schemerlampen_dressoir", "switch.tv_meubel", "switch.doorgang_keuken", "light.keuken", "switch.aanrecht"}
}
