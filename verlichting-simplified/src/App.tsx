import { useCallback, useEffect } from "react";
import "./App.css";

import { useImmer } from "use-immer";
import { Scene } from "./Scene";
import { ScenesList } from "./ScenesList";

type State = {
    scenes: Scene[] | "loading";
};

function App() {
    const [state, setState] = useImmer<State>({ scenes: "loading" });

    const fetchScenes = useCallback(() => {
        fetch("/api/scenes")
            .then((response) => response.json())
            .then((json) =>
                setState((draft: State) => {
                    draft.scenes = json.scenes;
                })
            );
    }, []);

    useEffect(fetchScenes, []);

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
}

export default App;
