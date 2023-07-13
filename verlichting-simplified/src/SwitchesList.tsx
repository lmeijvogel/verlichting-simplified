import { ListEntry } from "./ListEntry";
import { Switch } from "./Switch";

type Props = {
    switches: Switch[];
    onToggleSwitch: (theSwitch: Switch, newState: boolean) => void;
};

export function SwitchesList(props: Props) {
    return <ol>
        {props.switches.map(theSwitch => {
            return <ListEntry key={theSwitch.id} isActive={theSwitch.state === "on"} isLoading={false} onClick={() => props.onToggleSwitch(theSwitch, theSwitch.state === "off")}>{theSwitch.friendlyName}</ListEntry>
        })}
    </ol>;
}
