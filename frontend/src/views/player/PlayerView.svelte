<script lang="ts">
    import {GameState} from "../../service/net";
    import PlayerPlayView from "./PlayerPlayView.svelte";
    import PlayerLobbyView from "./PlayerLobbyView.svelte";
    import PlayerEndView from "./PlayerEndView.svelte";
    import { onDestroy } from 'svelte';
    import { state as gameStateStore} from '../../service/player/player';
    import {playerGameStore } from '../../service/gameStore';
    import { get } from "svelte/store";

    let views: Record<GameState,any>={
        [GameState.Lobby]: PlayerLobbyView,
        [GameState.Play]: PlayerPlayView,
        [GameState.Reveal]: PlayerLobbyView,
        [GameState.End]: PlayerEndView,
        [GameState.GameEndedState]: undefined,
        [GameState.Intermission]: undefined
    }

    onDestroy(() => {
        const game = get(playerGameStore);
        if (game) {
            game.signalPlayerLeaving();
        }
    });

</script>


<svelte:component this={views[$gameStateStore]} />