<script lang="ts">
    import type { LeaderboardEntry } from "../service/net";

    export let leaderboard: LeaderboardEntry[];
    export let finish = false;
    const finishColors = ["bg-yellow-400", "bg-gray-300", "bg-yellow-600"];
</script>

<div
    class="bg-white flex items-center rounded-xl p-4 flex-col gap-4 w-96 border-4 border-orange-500"
>
    <h2 class="text-center text-4xl font-bold" style="color:#FF6A00">
        Leaderboard
    </h2>
    <div class="w-full flex flex-col gap-3">
        {#each leaderboard as entry, i}
            <div class="flex items-center w-full">
                {#if finish && i < 3}
                    <div
                        class="{finishColors[
                            i
                        ]} flex-shrink-0 w-12 h-12 rounded-full text-white text-2xl font-bold flex items-center justify-center mr-4"
                    >
                        {i + 1}
                    </div>
                {:else}
                    <div
                        class="flex-shrink-0 w-12 text-center text-2xl font-bold text-gray-500 mr-4"
                    >
                        {i + 1}.
                    </div>
                {/if}
                <div
                    class="flex w-full items-center justify-between rounded-full p-1.5"
                    style="background-color:#464AA2"
                >
                    <div class="flex items-center">
                        <p class="text-lg font-bold text-white ml-4">
                            {entry.name}
                        </p>
                    </div>
                    <div
                        class="flex items-center justify-center rounded-xl bg-white px-4 py-1 mr-1.5"
                    >
                        <p
                            class="text-xl font-bold text-black pr-3 border-r mr-3"
                        >
                            {entry.points}
                        </p>
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-8 w-8 text-orange-500"
                            viewBox="0 -960 960 960"
                            fill="currentColor"
                        >
                            <path
                                d="M160-400q0-105 50-187t110-138q60-56 110-85.5l50-29.5v132q0 37 25 58.5t56 21.5q17 0 32.5-7t28.5-23l18-22q72 42 116 116.5T800-400q0 88-43 160.5T644-125q17-24 26.5-52.5T680-238q0-40-15-75.5T622-377L480-516 339-377q-29 29-44 64t-15 75q0 32 9.5 60.5T316-125q-70-42-113-114.5T160-400Zm320-4 85 83q17 17 26 38t9 45q0 49-35 83.5T480-120q-50 0-85-34.5T360-238q0-23 9-44.5t26-38.5l85-83Z"
                            />
                        </svg>
                        <p
                            class="text-lg font-bold text-[#79ACD9] border-gray-300"
                        >
                            {entry.correctCount}
                        </p>
                    </div>
                </div>
            </div>
        {/each}
    </div>
</div>
