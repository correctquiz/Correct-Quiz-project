<script lang="ts">
    import Button from "../../lib/Button.svelte";
    import PlayerNameCard from "../../lib/lobby/PlayerNameCard.svelte";
    import { gameCode, game, players } from "../../service/host/host";
    import QRCode from "svelte-qrcode";
    import { push } from "svelte-spa-router";
    import { hostGameStore } from "../../service/gameStore";
    import { get } from "svelte/store";

    game.navigate = push;

    $: joinUrl = `${window.location.origin}/#/?code=${$gameCode}`;

    function goBack() {
        game.unhost();
        push("/host");
    }

    function start() {
        const game = get(hostGameStore);
        if (game) {
            console.log("🚀 Host clicking START..."); 
            game.start();
        } else {
            console.error("❌ Game instance not found in store!");
        }
    }

    function handleKick(event: CustomEvent<string>) {
        const playerIdToKick = event.detail;
        game.kickPlayer(playerIdToKick);
    }
</script>

<div class="p-8 min-h-screen w-full bg-[#FCFCFC]">
    <button on:click={goBack} class="border-none cursor-pointer">
        <img
            src="../image/back.png"
            alt="back"
            class="w-10 h-10 hover:opacity-80"
        />
    </button>
    <div class="text-center">
        <h2 class="text-4xl font-bold" style="color: #464AA2">
            Join with game code
        </h2>
        <h2 class="text-6xl font-bold mt-4 text-black-100;">{$gameCode}</h2>
    </div>

    <div class="flex justify-end">
        <Button on:click={start} disabled={$players.length === 0}>Start Game</Button>
    </div>

    <div class="w-full flex justify-center pt-2">
        <div
            class="bg-white p-4 rounded-lg shadow-lg flex justify-center items-center"
        >
            <QRCode value={joinUrl} size={180} />
        </div>
    </div>

    <h2 class="mt-10 text-4xl font-bold" style="color: #464AA2;">
        players ({$players.length})
    </h2>

    <div class="flex flex-wrap gap-5 mt-10 ml-3">
        {#each $players as player (player.id)}
            <PlayerNameCard {player} on:kick={handleKick} />
        {:else}
            <p class="text-black text-2xl">No players have joined yet</p>
        {/each}
    </div>
</div>
