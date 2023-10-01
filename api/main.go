package main

import (
	"log"
	"net/http"

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
	router.GET("/api/scenes", func(c *gin.Context) { getEntities(c, ha.GetScenes, getAllowedSceneIds()) })
	router.POST("/api/start_scene/:id", func(c *gin.Context) { startScene(c, ha) })

	router.GET("/api/switches", func(c *gin.Context) { getEntities(c, ha.GetSwitches, getAllowedSwitchIds()) })
	router.POST("/api/set_switch/:id/:new_state", func(c *gin.Context) { setSwitchOrState(c, getAllowedSwitchIds(), ha.SetSwitch) })

	router.GET("/api/states", func(c *gin.Context) { getEntities(c, ha.GetStates, getAllowedStateIds()) })
	router.POST("/api/set_state/:id/:new_state", func(c *gin.Context) { setSwitchOrState(c, getAllowedStateIds(), ha.SetState) })

	router.Run(":3123")
}

func getEntities[T HasID](c *gin.Context, getter func() ([]T, error), whitelist []string) {
	typedEntities, err := getter()

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	allowedEntities := []T{}

	for i := range typedEntities {
		typedEntity := typedEntities[i]

		if slices.Contains(whitelist, typedEntity.GetID()) {
			allowedEntities = append(allowedEntities, typedEntity)
		}
	}

	c.IndentedJSON(http.StatusOK, allowedEntities)
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
	return []string{"scene.uit", "scene.ochtend", "scene.middag", "scene.avond", "scene.nacht"}
}

func getAllowedSwitchIds() []string {
	return []string{"switch.elektrische_deken", "switch.mechanische_ventilatie", "switch.babyfoon", "switch.tv_meubel"}
}

func getAllowedStateIds() []string {
	return []string{"input_boolean.vacation_mode", "input_boolean.elektrische_deken_in_gebruik", "input_boolean.kerstboom"}
}
