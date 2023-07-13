import { useCallback, useEffect } from "react";
import "./App.css";

import { useImmer } from "use-immer";
import { Scene } from "./Scene";
import { ScenesList } from "./ScenesList";
import { Switch } from "./Switch";
import { SwitchesList } from "./SwitchesList";

type State = {
    scenes: Scene[] | "loading";
    switches: Switch[] | "loading";
};

function App() {
    const [state, setState] = useImmer<State>({ scenes: "loading", switches: "loading" });

    const fetchScenes = useCallback(() => {
        fetch("/api/scenes")
            .then((response) => response.json())
            .then((json) =>
                setState((draft: State) => {
                    draft.scenes = json.scenes;
                })
            );
    }, []);

    const fetchSwitches = useCallback(() => {
        fetch("/api/switches")
            .then((response) => response.json())
            .then((json) =>
                setState((draft: State) => {
                    draft.switches = json.switches;
                })
            );
    }, []);

    useEffect(fetchScenes, []);
    useEffect(fetchSwitches, []);

    useEffect(() => {
        const interval = setInterval(fetchScenes, 10000);

        return () => clearInterval(interval);
    });

    const onActivateScene = useCallback((scene: Scene) => {
        fetch(`/api/start_scene/${scene.id}`, {
            method: "POST"
        })
            .then((response) => response.json())
            .then(updateScene);
    }, []);

    const onToggleSwitch = useCallback((theSwitch: Switch, newState: boolean) => {
        fetch(`/api/set_switch/${theSwitch.id}/${newState ? "on" : "off"}`, {
            method: "POST"
        })
            .then((response) => response.json())
            .then(updateSwitch);
    }, []);

    return (
        <div className="App">
            <h1>Verlichting</h1>
            <div className="card">
                {state.scenes === "loading" ? (
                    "Loading..."
                ) : (
                    <ScenesList scenes={state.scenes} onActivateScene={onActivateScene} />
                )}
            </div>
            <h1>Schakelaars</h1>
            <div className="card">
                {state.switches === "loading" ? (
                    "Loading..."
                ) : (
                    <SwitchesList switches={state.switches} onToggleSwitch={onToggleSwitch} />
                )}
            </div>
        </div>
    );

    function updateScene(newScene: Scene) {
        setState((draft) => {
            if (draft.scenes === "loading") return;

            const scene = draft.scenes.find((scene) => scene.id === newScene.id);

            if (!scene) return;

            scene.lastActivated = newScene.lastActivated;
        });
    }

    function updateSwitch(newSwitch: Switch) {
        setState((draft) => {
            if (draft.switches === "loading") return;

            const theSwitch = draft.switches.find((s) => s.id === newSwitch.id);

            if (!theSwitch) return;

            theSwitch.state = newSwitch.state;
        });
    }
}

export default App;
