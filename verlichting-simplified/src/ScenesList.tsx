import { useState } from "react";
import { Scene } from "./Scene";
import { SceneEntry } from "./SceneEntry";

const SceneOrder = ["scene.uit", "scene.ochtend", "scene.middag", "scene.avond", "scene.nacht"];

type ScenesListProps = {
    onActivateScene: (scene: Scene) => void;
    scenes: Scene[];
};

export function ScenesList({ onActivateScene, scenes }: ScenesListProps) {
    const [loadingScene, setLoadingScene] = useState<Scene | undefined>();

    const sortedScenes = sortScenes(scenes);

    const currentScene = maxBy(scenes, "lastActivated");

    const onSceneClick = (scene: Scene) => {
        onActivateScene(scene);
        setLoadingScene(scene);
    };

    return (
        <ol>
            {sortedScenes.map((scene) => (
                <SceneEntry
                    key={scene.id}
                    isCurrent={scene.id === currentScene.id}
                    isLoading={scene.id === loadingScene?.id}
                    onClick={() => onSceneClick(scene)}
                >
                    {scene.friendlyName}
                </SceneEntry>
            ))}
        </ol>
    );
}

const sortScenes: (scenes: Scene[]) => Scene[] = (scenes: Scene[]) => {
    const tempScenes = [...scenes];

    const result = [];

    for (const sceneInOrder of SceneOrder) {
        const sceneIndex = tempScenes.findIndex(sc => sc.id === sceneInOrder);

        result.push(tempScenes[sceneIndex]);

        tempScenes.splice(sceneIndex, 1);
    }

    return [...result, ...tempScenes];
}

function maxBy<T>(input: T[], field: keyof T): T {
    let result: T = input[0];

    for (const element of input) {
        if (element[field] > result[field]) {
            result = element;
        }
    }

    return result;
}
