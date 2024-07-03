package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

type Switch struct {
	ID           string `json:"id"`
	FriendlyName string `json:"friendlyName"`
	State        string `json:"state"`
}

type HaScene struct {
	ID            string    `json:"id"`
	FriendlyName  string    `json:"friendlyName"`
	LastActivated time.Time `json:"lastActivated"`
}

type Scene struct {
	ID           string `json:"id"`
	FriendlyName string `json:"friendlyName"`
	State        string `json:"state"`
}

type AllowedStateValue string

const (
	On  AllowedStateValue = "on"
	Off                   = "off"
)

type haEntity struct {
	EntityId   string             `json:"entity_id"`
	State      string             `json:"state"`
	Attributes haEntityAttributes `json:"attributes"`
}

type haEntityAttributes struct {
	FriendlyName string `json:"friendly_name"`
}

type HomeAssistantAdapter struct {
	host      string
	api_token string
}

func (theSwitch Switch) GetID() string {
	return theSwitch.ID
}

func (scene Scene) GetID() string {
	return scene.ID
}

func NewAdapter(host string, api_token string) HomeAssistantAdapter {
	return HomeAssistantAdapter{host, api_token}
}

func (ha HomeAssistantAdapter) GetScenes() ([]Scene, error) {
	scenes, err := getTypedEntities(ha, deserializeScene)

	if err != nil {
		return nil, err
	}

	currentScene := slices.MaxFunc(scenes, func(one HaScene, two HaScene) int {
		if one.LastActivated.Before(two.LastActivated) {
			return -1
		} else if one.LastActivated.After(two.LastActivated) {
			return 1
		}

		return 0
	})

	var result []Scene

	for i := range scenes {
		deserializedScene, err := toScene(scenes[i], currentScene.ID)

		if err != nil {
			continue
		}

		result = append(result, deserializedScene)
	}

	return result, nil
}

func toScene(haScene HaScene, currentSceneID string) (Scene, error) {
	isCurrent := haScene.ID == currentSceneID

	var state string

	if isCurrent {
		state = "on"
	} else {
		state = "off"
	}

	return Scene{haScene.ID, haScene.FriendlyName, state}, nil
}

func (ha HomeAssistantAdapter) GetStates() ([]Switch, error) {
	return getTypedEntities(ha, deserializeSwitch)
}

func (ha HomeAssistantAdapter) GetLights() ([]Switch, error) {
	return getTypedEntities(ha, deserializeSwitch)
}

func (ha HomeAssistantAdapter) StartScene(sceneId string) (*Scene, error) {
	url := fmt.Sprintf("%v/api/services/scene/turn_on", ha.host)

	data := fmt.Sprintf(`{"entity_id": "%v" }`, sceneId)

	response, err := performRequestWithMethodAndBody(url, ha, "POST", data)

	if err != nil {
		return nil, err
	}

	haScene, err := decodeResponse(response, deserializeScene)

	if err != nil {
		return nil, err
	}

	scene, err := toScene(*haScene, haScene.ID)

	if err != nil {
		return nil, err
	}

	return &scene, nil
}

func (ha HomeAssistantAdapter) GetSwitches() ([]Switch, error) {
	return getTypedEntities(ha, deserializeSwitch)
}

func (ha HomeAssistantAdapter) SetSwitch(switchId string, newState AllowedStateValue) (*Switch, error) {
	command := "turn_off"

	if newState == On {
		command = "turn_on"
	}
	url := fmt.Sprintf("%v/api/services/switch/%v", ha.host, command)

	data := fmt.Sprintf(`{"entity_id": "%v" }`, switchId)

	response, err := performRequestWithMethodAndBody(url, ha, "POST", data)

	if err != nil {
		return nil, err
	}

	theSwitch, err := decodeResponse(response, deserializeSwitch)

	return theSwitch, err
}

func (ha HomeAssistantAdapter) SetLight(switchId string, newState AllowedStateValue) (*Switch, error) {
	command := "turn_off"

	if newState == On {
		command = "turn_on"
	}

	var url string

	if strings.HasPrefix(switchId, "light") {
		url = fmt.Sprintf("%v/api/services/light/%v", ha.host, command)
	} else {
		url = fmt.Sprintf("%v/api/services/switch/%v", ha.host, command)
	}

	data := fmt.Sprintf(`{"entity_id": "%v" }`, switchId)

	response, err := performRequestWithMethodAndBody(url, ha, "POST", data)

	if err != nil {
		return nil, err
	}

	theSwitch, err := decodeResponse(response, deserializeSwitch)

	return theSwitch, err
}

func (ha HomeAssistantAdapter) SetState(stateId string, newState AllowedStateValue) (*Switch, error) {
	command := "turn_off"

	if newState == On {
		command = "turn_on"
	}
	url := fmt.Sprintf("%v/api/services/input_boolean/%v", ha.host, command)

	data := fmt.Sprintf(`{"entity_id": "%v" }`, stateId)

	response, err := performRequestWithMethodAndBody(url, ha, "POST", data)

	if err != nil {
		return nil, err
	}

	state, err := decodeResponse(response, deserializeSwitch)

	return state, err
}

func decodeResponse[T any](response *http.Response, deserializer func(entity haEntity) (*T, error)) (*T, error) {
	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var entities []haEntity

	err = json.Unmarshal(responseData, &entities)

	if err != nil {
		return nil, err
	}

	if len(entities) == 0 {
		return nil, nil
	}

	state, err := deserializer(entities[0])

	if err != nil {
		return nil, err
	}

	return state, err
}

func deserializeScene(entity haEntity) (*HaScene, error) {
	lastActivated, err := time.Parse(time.RFC3339, entity.State)

	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(entity.EntityId, "scene.") {
		return nil, errors.New("Entity is not a scene")
	}
	scene := HaScene{entity.EntityId, entity.Attributes.FriendlyName, lastActivated}

	return &scene, nil
}

func deserializeSwitch(entity haEntity) (*Switch, error) {
	theSwitch := Switch{entity.EntityId, entity.Attributes.FriendlyName, entity.State}

	return &theSwitch, nil
}

func performRequest(url string, ha HomeAssistantAdapter) (*http.Response, error) {

	return performRequestWithMethodAndBody(url, ha, "GET", "")
}

func performRequestWithMethodAndBody(url string, ha HomeAssistantAdapter, method string, body string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(body)))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", ha.api_token))

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Upstream server error, status %v", response.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

func getTypedEntities[T any](ha HomeAssistantAdapter, converter func(input haEntity) (*T, error)) ([]T, error) {
	url := fmt.Sprintf("%v/api/states", ha.host)

	response, err := performRequest(url, ha)

	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var entities []haEntity

	err = json.Unmarshal(responseData, &entities)

	if err != nil {
		return nil, err
	}

	typedEntities := []T{}

	for i := range entities {
		entity := entities[i]

		convertedEntity, err := converter(entity)

		if err != nil {
			continue
		}

		if convertedEntity != nil {
			typedEntities = append(typedEntities, *convertedEntity)
		}
	}

	return typedEntities, nil
}
