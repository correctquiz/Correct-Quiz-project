import { get, Writable, writable } from "svelte/store";
import { PlayerGame } from "./player/player";
import { HostGame } from "./host/host";

export const playerGameStore: Writable<PlayerGame | null> = writable(null);

export function initializePlayerGame(
    navigateFunction: (path: string) => void,
) {
    if (get(playerGameStore) === null) {
        playerGameStore.set(new PlayerGame(navigateFunction));
    }
}
export const hostGameStore: Writable<HostGame | null> = writable(null);

export function initializeHostGame(navigateFunction: (path: string) => void) {
    if (get(hostGameStore) === null) {
        const hostGame = new HostGame(navigateFunction);
        hostGameStore.set(hostGame);
        const token = localStorage.getItem("jwt_token");
        if (token) {
            hostGame.connect(token);
        } else {
            console.error("No token found for Host!");
        }
    }
}