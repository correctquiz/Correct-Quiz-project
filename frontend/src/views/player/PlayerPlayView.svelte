<script lang="ts">
    import { COLORS } from "../../model/quiz";
    import {
        currentPlayer,
        currentQuestion,
        points,
        maxStreak,
        streakBonus,
    } from "../../service/player/player";
    import {
        PacketTypes,
        type PlayerAnswerFeedbackPacket,
    } from "../../service/net";
    import { onMount } from "svelte";
    import { playerGameStore } from "../../service/gameStore";

    let selectedAnswerIndex: number | null = null;
    let isAnswerCorrect: boolean | null = null;
    let showResult: boolean = false;
    let actualCorrectIndex: number[] = [];
    let awardedStreakBonus: number = 0;

    function onClick(i: number) {
        selectedAnswerIndex = i;
        const questionIndex = $currentQuestion?.index;
        if (questionIndex !== undefined) {
            $playerGameStore.answer(questionIndex, i);
        } else {
            console.error("Cannot send answer: question index is undefined.");
        }
    }

    onMount(() => {
        $playerGameStore.onMessage((packet) => {
            if (packet.id === PacketTypes.PlayerAnswerFeedback) {
                const feedback = packet as PlayerAnswerFeedbackPacket;
                isAnswerCorrect = feedback.isCorrect;
                actualCorrectIndex = feedback.correctAnswerIndex;
                awardedStreakBonus = feedback.streakBonus;
                showResult = true;
            }
        });
    });

    function resetForNewQuestion() {
        selectedAnswerIndex = null;
        isAnswerCorrect = null;
        showResult = false;
        actualCorrectIndex = [];
        awardedStreakBonus = 0;
    }
    $: if ($currentQuestion) {
        resetForNewQuestion();
    }
</script>

<div class="min-h-screen h-screen flex flex-col bg-gray-200">
    <div class="relative p-4 shadow-xl flex items-center">
        <div class="flex-shrink-0 mr-auto">
            <h3 class="text-xl font-bold text-gray-800">คะแนน</h3>
            <p class="text-3xl font-bold text-gray-800">{$points}</p>
        </div>

        <div class="absolute inset-0 flex flex-col justify-center items-center">
            <p class="text-lg font-semibold text-gray-600">Name :</p>
            <p class="text-4xl font-bold text-black">
                {$currentPlayer.name}
            </p>
        </div>

        <div class="flex-1 text-right">
            <h3 class="text-xl font-bold text-gray-800">Streak สูงสุด</h3>

            <div
                class="text-3xl font-bold flex items-center justify-end {$maxStreak >
                0
                    ? 'text-orange-500'
                    : 'text-gray-800'}"
            >
                <span>{$maxStreak}</span>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="ml-1"
                    height="24px"
                    viewBox="0 -960 960 960"
                    width="24px"
                    fill="#FFA32F"
                    ><path
                        d="M160-400q0-105 50-187t110-138q60-56 110-85.5l50-29.5v132q0 37 25 58.5t56 21.5q17 0 32.5-7t28.5-23l18-22q72 42 116 116.5T800-400q0 88-43 160.5T644-125q17-24 26.5-52.5T680-238q0-40-15-75.5T622-377L480-516 339-377q-29 29-44 64t-15 75q0 32 9.5 60.5T316-125q-70-42-113-114.5T160-400Zm320-4 85 83q17 17 26 38t9 45q0 49-35 83.5T480-120q-50 0-85-34.5T360-238q0-23 9-44.5t26-38.5l85-83Z"
                    /></svg
                >

                <span
                    class="text-xl text-[#79ACD9] transition-opacity ml-2 {$streakBonus >
                        0 && showResult
                        ? 'opacity-100'
                        : 'opacity-0'}"
                >
                    (+{$streakBonus})
                </span>
            </div>
        </div>
    </div>

    <div class="flex-grow grid grid-cols-2">
        {#each $currentQuestion?.choices || [] as choice, i}
            <button
                class="relative w-full h-full p-4 text-left transition-opacity duration-300 {showResult &&
                actualCorrectIndex.includes(i)
                    ? 'bg-green-500'
                    : showResult &&
                        i === selectedAnswerIndex &&
                        !isAnswerCorrect
                      ? 'bg-red-500'
                      : COLORS[i]}"
                class:opacity-40={(selectedAnswerIndex !== null &&
                    !showResult) ||
                    (showResult &&
                        !actualCorrectIndex.includes(i) &&
                        i !== selectedAnswerIndex)}
                on:click={() => onClick(i)}
                disabled={selectedAnswerIndex !== null || showResult}
            >
                {#if choice.imageUrl}
                    <div
                        class="absolute top-0 left-0 right-0 h-2/3 flex justify-center items-start mt-2"
                    >
                        <img
                            src={choice.imageUrl}
                            alt={`Choice ${i + 1} Image`}
                            class="w-full h-full object-contain"
                        />
                    </div>
                {/if}

                <div
                    class="absolute bottom-2 left-4 flex items-center space-x-4"
                >
                    {#if i === 0}
                        <img
                            src="../image/choices/circle.png"
                            alt="Circle"
                            class="w-24 h-22"
                        />
                    {:else if i === 1}
                        <img
                            src="../image/choices/triangle.png"
                            alt="Triangle"
                            class="w-29 h-27"
                        />
                    {:else if i === 2}
                        <img
                            src="../image/choices/star.png"
                            alt="Star"
                            class="w-29 h-27"
                        />
                    {:else if i === 3}
                        <img
                            src="../image/choices/square.png"
                            alt="Square"
                            class="w-24 h-22"
                        />
                    {/if}

                    <p
                        class="text-5xl font-bold text-white"
                        style="text-shadow: 2px 2px 4px rgba(0,0,0,0.4);"
                    >
                        {choice.name}
                    </p>
                </div>
            </button>
        {/each}
    </div>
</div>
