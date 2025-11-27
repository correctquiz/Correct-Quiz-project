import { Writable, writable } from "svelte/store";
import { PlayerGame } from "./player/player";
import { push } from "svelte-spa-router";

export const playerGameStore: Writable<PlayerGame | null> = writable(null);