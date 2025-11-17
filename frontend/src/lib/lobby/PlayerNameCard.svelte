<script lang="ts">
    import type { Player } from "../../model/quiz";
    import { createEventDispatcher } from 'svelte';

    export let player: Player;

     const dispatch = createEventDispatcher();

    function handleKick() {
        dispatch("kick", player.id); 
    }
</script>

<button on:click={handleKick}>
    <h3
        class="text-3xl text-white font-bold px-4 py-2 rounded-xl hover:line-through cursor-pointer bg-[#F87923]"
    >
        {player.name}
    </h3>
</button>