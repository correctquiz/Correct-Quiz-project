<script lang="ts">
    import Clock from "../../lib/Clock.svelte";
    import { currentQuestion, tick, } from "../../service/host/host"
</script>

{#if $currentQuestion != null}
    <div
        class="min-h-screen h-screen flex flex-col bg-[#FCFCFC]"
    >
        <div class="h-6"></div>
        <div class="flex justify-between pt-20 pb-12.5">
            <div class="flex justify-between items-center">
                <Clock>
                    <span class="text-4xl">{$tick}</span>
                </Clock>
            </div>
            {#if $currentQuestion.imageUrl}
                <img
                    src={$currentQuestion.imageUrl}
                    alt="Question visual"
                    class="max-h-[250px] max-w-[500px] object-contain rounded-lg shadow-md"
                />
            {/if}
            <div
                class="text-3xl text-black h-15 p-2 font-bold text-center"
                style="border: 5px solid #464AA2; border-radius:10px; width: fit-content; margin-right: 24px;"
            >
                <h2>ข้อ {$currentQuestion.index + 1}</h2>
            </div>
        </div>

        <div
            class="text-4xl text-white h-25 border-b p-4 font-bold text-center"
            style="background-color: #F87923; display: flex; justify-content: center; align-items: center;"
        >
            {$currentQuestion.name}
        </div>
    </div>
{/if}
