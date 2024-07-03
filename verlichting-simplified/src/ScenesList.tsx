import { useState } from "react";
import { ListEntry } from "./ListEntry";
import { Switch } from "./Switch";

const SceneOrder = ["scene.uit", "scene.ochtend", "scene.middag", "scene.avond", "scene.nacht"];

type ScenesListProps = {
    onActivateScene: (scene: Switch) => void;
    scenes: Switch[];
};

export function ScenesList({ onActivateScene, scenes }: ScenesListProps) {
    const [loadingScene, setLoadingScene] = useState<Switch | undefined>();

    const sortedScenes = sortScenes(scenes);

    const onSceneClick = (scene: Switch) => {
        onActivateScene(scene);
        setLoadingScene(scene);
    };

    return (
        <ol>
            {sortedScenes.map((scene) => (
                <ListEntry
                    key={scene.id}
                    isActive={scene.state === "on"}
                    isLoading={scene.id === loadingScene?.id}
                    onClick={() => onSceneClick(scene)}
                >
                    {scene.friendlyName}
                </ListEntry>
            ))}
        </ol>
    );
}

const sortScenes: (scenes: Switch[]) => Switch[] = (scenes: Switch[]) => {
    const tempScenes = [...scenes];

    const result = [];

    for (const sceneInOrder of SceneOrder) {
        const sceneIndex = tempScenes.findIndex(sc => sc.id === sceneInOrder);

        result.push(tempScenes[sceneIndex]);

        tempScenes.splice(sceneIndex, 1);
    }

    return [...result, ...tempScenes];
}
