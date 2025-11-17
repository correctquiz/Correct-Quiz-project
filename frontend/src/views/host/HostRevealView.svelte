<script lang="ts">
	import Button from "../../lib/Button.svelte";
	import {
		players,
		answerCounts,
		currentQuestion,
		correctAnswerIndex,
		type HostGame,
	} from "../../service/host/host";
	import { COLORS } from "../../model/quiz";

	const choiceImages = [
		"../image/choices/circle.png",
		"../image/choices/triangle.png",
		"../image/choices/star.png",
		"../image/choices/square.png",
	];

	export let game: HostGame;

</script>

<div class="min-h-screen bg-[#FCFCFC]">

	<div class="pt-10"></div>
	<div
		class="text-4xl text-white h-25 border-b font-bold text-center"
		style="background-color: #F87923; display: flex; justify-content: center; align-items: center;"
	>
		{$currentQuestion.name}
	</div>
	<div class="flex justify-between items-center pt-10 px-5">
		<div
			class="text-3xl text-black h-15 p-2 font-bold text-center"
			style="border: 5px solid #464AA2; border-radius:10px; width: fit-content;  margin-right: 20px;"
		>
			<h2>ข้อ {$currentQuestion.index + 1}</h2>
		</div>
		<Button on:click={() => game.intermission()}>ต่อไป</Button>
	</div>

	<div class="max-w-4xl mx-auto">

		<div class="max-w-5xl mx-auto mt-8 flex flex-col gap-4">
        
        {#each $currentQuestion?.choices || [] as choice, i}
            {@const totalPlayers = $players.length}
			{@const currentAnswerCount = $answerCounts[i] || 0}
			{@const percentage = totalPlayers > 0 ? (($answerCounts[i] || 0) / totalPlayers) * 100 : 0}
            
            <div class="bg-white/60 p-4 rounded-lg shadow-lg flex items-center gap-4 transition-all duration-300">
                
                <div class="flex-shrink-0 w-1/10 font-bold text-2xl flex items-center gap" style="color: {COLORS[i]}">
                     <span class="text-5xl">{['●', '▲', '★', '■'][i]}</span>
                </div>

                <div class="flex-grow flex items-center gap-4">
                    <div class="w-full bg-gray-300 rounded-full h-10 border-2 border-gray-400">
						{#if percentage > 0}
                        <div 
                            class="h-full rounded-full text-white flex items-center justify-end text-xl font-bold pr-3 transition-all duration-500 {COLORS[i]}"
                            style="width: {percentage}%;"
                        >
                        </div>
						{/if}
                    </div>
                    <span class="font-bold text-3xl w-16 text-right">{currentAnswerCount}</span>
                </div>

                <div class="w-12 h-12 text-5xl flex items-center justify-center">
                    {#if $correctAnswerIndex.includes(i)}
                        <span><svg xmlns="http://www.w3.org/2000/svg" height="40px" viewBox="0 -960 960 960" width="40px" fill="#44BF0B"><path d="M400-291.91 228.15-464.14l68.08-67.71L400-428.09l264.14-263.76 67.71 67.71L400-291.91Z"/></svg></span>
                    {/if}
                </div>
            </div>
        {/each}
		</div>
	</div>
</div>
