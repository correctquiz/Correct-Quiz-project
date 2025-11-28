<script lang="ts">
    import { gameCode, isHostNavigating, state } from "../../service/host/host";
    import { GameState } from "../../service/net";
    import HostEndView from "./HostEndView.svelte";
    import HostLobbyView from "./HostLobbyView.svelte";
    import HostPlayView from "./HostPlayView.svelte";
    import HostRevealView from "./HostRevealView.svelte";
    import HostIntermissionView from "./HostIntermissionView.svelte";
    import { onDestroy } from "svelte";
    import { get } from "svelte/store";
    import { hostGameStore } from "../../service/gameStore";

    onDestroy(() => {
        if (!get(isHostNavigating)) {
            const currentGame = get(hostGameStore);
            if (currentGame) {
                currentGame.unhost();
            }
        }
        isHostNavigating.set(false);
    });

    let views: Record<GameState, any> = {
        [GameState.Lobby]: HostLobbyView,
        [GameState.Play]: HostPlayView,
        [GameState.Reveal]: HostRevealView,
        [GameState.Intermission]: HostIntermissionView,
        [GameState.End]: HostEndView,
        [GameState.GameEndedState]: HostEndView,
    };
</script>

{#if $gameCode != null && $hostGameStore != null}
    <svelte:component this={views[$state]} game={$hostGameStore} />
{:else}
    <div class="min-h-screen w-full flex justify-center items-center">
        <h2 class="text-4xl font-bold" style="color: #464AA2">
            Loading...
        </h2>
    </div>
{/if}
