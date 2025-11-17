<script lang="ts">
    import {currentPlayer } from "../../service/player/player";
</script>

<div
    class="min-h-screen flex justify-center bg-[#FCFCFC]"
>
    <div class="px-4">
        <img
            src="../image/CorrectQuiz.png"
            alt="Quiz icon"
            class="w-52 h-53 mx-auto mt-18"
            style=" filter: drop-shadow(4px 4px 2px rgba(0,0,0,0.2));"
        />
        <p class="text-2xl text-center font-bold pb-4">Your Name is:</p>

        <h3 class="text-3xl text-white font-bold px-4 py-2 rounded-xl text-center bg-[#F87923]">
            {$currentPlayer.name}
        </h3>

        <p class="text-gray-800 pt-5">Please wait for the host to start the game....</p>
    </div>
</div>
