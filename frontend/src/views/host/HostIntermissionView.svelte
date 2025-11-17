<script lang="ts">
    import { leaderboard,type HostGame } from "../../service/host/host";
    export let game: HostGame;
</script>

<div class="relative min-h-screen w-full bg-[#FCFCFC]">
    <button
        on:click={() => game.nextQuestion()}
        class="absolute top-5 right-5 z-10 px-6 py-2 bg-[#FF6A00] text-white font-bold text-xl rounded-lg hover:bg-[#E05C00] mr-4 mt-30 cursor-pointer"
    >
        Next
    </button>

    <div class="absolute top-0 left-1/2 -translate-x-1/2 mt-30">
        <div
            class="bg-white flex items-center rounded-xl p-4 flex-col gap-4 w-96 border-3 border-[#FF6A00]"
        >
            <h2 class="text-center text-4xl font-bold" style="color:#FF6A00">
                Leaderboard
            </h2>
            {#each $leaderboard as entry}
                <div
                    class="flex w-full max-w-xs items-center justify-between rounded-full p-1.5"
                    style="background-color:#464AA2"
                >
                    <div class="flex items-center ml-4">
                        <p class="text-lg font-bold text-white">
                            {entry.name}
                        </p>
                    </div>
                    <div
                        class="flex items-center justify-center rounded-xl bg-white px-3 py-1 mr-1.5"
                    >
                        <p class="text-xl font-bold text-black">
                            {entry.points}
                        </p>
                    </div>
                </div>
            {/each}
        </div>
    </div>
</div>
