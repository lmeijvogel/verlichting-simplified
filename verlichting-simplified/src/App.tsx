import { useCallback, useEffect } from "react";
import "./App.css";

import { Updater, useImmer } from "use-immer";
import { Switch } from "./Switch";
import { SwitchesList } from "./SwitchesList";

type State = {
    switchableEntities: {
        scenes: Switch[] | "loading";
        switches: Switch[] | "loading";
        lights: Switch[] | "loading";
        states: Switch[] | "loading";
    };
};

function UseSwitches(setState: Updater<State>, url: string, updateBasePath: string, field: keyof State["switchableEntities"]) {
    const fetchData = () => {
        fetch(url)
            .then((response) => response.json())
            .then((json) =>
                setState((draft: State) => {
                    draft.switchableEntities[field] = json;
                })
            );
    };

    const updateEntity = useCallback((newSwitch: Switch) => {
        setState((draft) => {
            const relevantCollection = draft.switchableEntities[field];

            if (relevantCollection === "loading") return;

            const theSwitch = relevantCollection.find((s) => s.id === newSwitch.id);

            if (!theSwitch) return;

            theSwitch.state = newSwitch.state;
        });
    }, []);

    const onToggleEntity = useCallback((theSwitch: Switch, newState: boolean) => {
        fetch(`${updateBasePath}/${theSwitch.id}/${newState ? "on" : "off"}`, {
            method: "POST"
        })
            .then((response) => response.json())
            .then(updateEntity);
    }, [updateEntity]);

    useEffect(fetchData, []);

    useEffect(() => {
        const interval = setInterval(fetchData, 10000);

        return () => clearInterval(interval);
    });

    return [onToggleEntity];
}

function App() {
    const [state, setState] = useImmer<State>({
        switchableEntities: { switches: "loading", lights: "loading", states: "loading", scenes: "loading" }
    });

    const [onActivateScene] = UseSwitches(setState, "/api/scenes", "api/start_scene", "scenes");
    const [onToggleSwitch] = UseSwitches(setState, "/api/switches", "/api/set_switch", "switches");
    const [onToggleLight] = UseSwitches(setState, "/api/lights", "/api/set_light", "lights");
    const [onToggleState] = UseSwitches(setState, "/api/states", "/api/set_state", "states");

    return (
        <div className="App">
            <h1>Scenes</h1>
            <div className="card">
                {state.switchableEntities.scenes === "loading" ? (
                    "Loading..."
                ) : (
                    <SwitchesList switches={state.switchableEntities.scenes} onToggleSwitch={onActivateScene} />
                )}
            </div>
            <h1>Schakelaars</h1>
            <div className="card">
                {state.switchableEntities.switches === "loading" ? (
                    "Loading..."
                ) : (
                    <SwitchesList switches={state.switchableEntities.switches} onToggleSwitch={onToggleSwitch} />
                )}
            </div>
            <h1>Lichten</h1>
            <div className="card">
                {state.switchableEntities.lights === "loading" ? (
                    "Loading..."
                ) : (
                    <SwitchesList switches={state.switchableEntities.lights} onToggleSwitch={onToggleLight} />
                )}
            </div>
            <h1>States</h1>
            <div className="card">
                {state.switchableEntities.states === "loading" ? (
                    "Loading..."
                ) : (
                    <SwitchesList switches={state.switchableEntities.states} onToggleSwitch={onToggleState} />
                )}
            </div>
        </div>
    );
}

export default App;
