<script lang="ts">
    import { push } from "svelte-spa-router";
    import {
        currentPlayer,
        maxStreak,
        points,
        rank,
        resetPlayerStores,
    } from "../../service/player/player";

    function goBack() {
        push("/");
        resetPlayerStores();
    }
</script>

<div class="min-h-screen flex justify-center bg-[#FCFCFC]">
    <button
        on:click={goBack}
        class="absolute top-8 left-8 border-none cursor-pointer z-10"
    >
        <img
            src="../image/back.png"
            alt="back"
            class="w-10 h-10 hover:opacity-80"
        />
    </button>
    <div>
        <img
            src="../image/CorrectQuiz.png"
            alt="Quiz icon"
            class="w-52 h-53 mx-auto mt-18"
            style=" filter: drop-shadow(4px 4px 2px rgba(0,0,0,0.2));"
        />
         <h1 class="text-3xl text-black font-bold text-center">
            {$currentPlayer.name}
        </h1>
        <p class="text-2xl text-center pb-4">Your learderboad is:</p>

        <h3 class="text-center text-5xl font-bold">
            NO. {$rank > 0 ? $rank : "-"}
        </h3>

        <p class="text-center mt-3 text-xl">&</p>

        
        <div class="flex justify-center mt-4">
            <p class="text-2xl text-gray-800 px-4 py-2">
                Your score is:
                <span
                    class="text-2xl text-white font-bold rounded-xl text-center px-4 py-1 bg-[#F87923]"
                    >{$points}</span
                >
            </p>
            <div
            class="flex items-center justify-center gap-2 {$maxStreak > 0
                ? 'text-orange-500'
                : 'text-gray-800'}"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-8 w-8"
                viewBox="0 -960 960 960"
                fill="currentColor"
            >
                <path
                    d="M160-400q0-105 50-187t110-138q60-56 110-85.5l50-29.5v132q0 37 25 58.5t56 21.5q17 0 32.5-7t28.5-23l18-22q72 42 116 116.5T800-400q0 88-43 160.5T644-125q17-24 26.5-52.5T680-238q0-40-15-75.5T622-377L480-516 339-377q-29 29-44 64t-15 75q0 32 9.5 60.5T316-125q-70-42-113-114.5T160-400Zm320-4 85 83q17 17 26 38t9 45q0 49-35 83.5T480-120q-50 0-85-34.5T360-238q0-23 9-44.5t26-38.5l85-83Z"
                />
            </svg>
            <span class="text-2xl font-bold text-[#79ACD9]">{$maxStreak}</span>
        </div>
        </div>
    </div>
</div>
