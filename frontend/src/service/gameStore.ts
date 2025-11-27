import { get, Writable, writable } from "svelte/store";
import { PlayerGame } from "./player/player";

export const playerGameStore: Writable<PlayerGame | null> = writable(null);

export function initializePlayerGame(
    navigateFunction: (path: string) => void,
) {
    if (get(playerGameStore) === null) {
        playerGameStore.set(new PlayerGame(navigateFunction));
    }
}
