import { writable } from "svelte/store";
import { PlayerGame } from "./player/player";
import { push } from "svelte-spa-router";

export const playerGameStore = writable(new PlayerGame(push));