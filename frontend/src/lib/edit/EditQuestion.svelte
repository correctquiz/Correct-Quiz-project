<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import { COLORS, type QuizQuestion } from "../../model/quiz";
    import Button from "../Button.svelte";
    import QuizChoiceCard from "../play/QuizChoiceCard.svelte";
    import Clock from "../Clock.svelte";

    const dispatch = createEventDispatcher();

    export let selectedQuestion: QuizQuestion;
    export let imagePreviews: Map<number, string>;
    export let errorForThisQuestion: string | undefined = undefined;

    let fileInput: HTMLInputElement;
    let timeError: string | null = null;

    function notifyParentOfChange() {
        dispatch("change");
    }

    function triggerFileInput() {
        fileInput.click();
    }

    async function uploadImageForChoice(event: Event, choiceIndex: number) {
        const input = event.target as HTMLInputElement;
        if (!input.files || !input.files[0] || !selectedQuestion) {
            return;
        }

        const file = input.files[0];
        const reader = new FileReader();

        reader.onload = (e) => {
            const resultUrl = e.target?.result as string;

            selectedQuestion.choices[choiceIndex].imageUrl = resultUrl;

            selectedQuestion = selectedQuestion;

            notifyParentOfChange();
        };

        reader.readAsDataURL(file);
    }

    function onFileSelected(event: Event) {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files[0]) {
            const file = input.files[0];
            const reader = new FileReader();

            reader.onload = (e) => {
                const resultUrl = e.target?.result as string;
                imagePreviews.set(selectedQuestion.id, resultUrl);
                imagePreviews = imagePreviews;
                notifyParentOfChange();
            };

            reader.readAsDataURL(file);
        }
    }

    function onDelete() {
        dispatch("delete");
    }

    function onTimeChange(e: Event) {
        let target = e.target as HTMLInputElement;
        const newTime = parseInt(target.value);
        selectedQuestion.time = newTime;
        notifyParentOfChange();
    }
</script>

<input
    type="file"
    class="hidden"
    accept="image/*"
    bind:this={fileInput}
    on:change={onFileSelected}
/>

<div class="p-14 w-full flex-1 bg-[#FCFCFC]">
    <div class="flex-1 flex flex-col min-h-[95%] border">
        <div
            class="flex font-bold text-center items-center border-b p-2 justify-between"
            style="background-color:#F87923;"
        >
            <div class="w-32"></div>
            <input
                class="p-4 text-3xl text-center"
                style="background-color:#F87923;"
                on:change
                bind:value={selectedQuestion.name}
            />
            <div class="w-32">
                <Button on:click={onDelete}>Delete</Button>
            </div>
        </div>

        <div class="flex justify-between items-center w-full p-4 gap-4">
            <div class="w-1/4 flex justify-start">
                <Clock>
                    <input
                        value={selectedQuestion.time}
                        on:change={onTimeChange}
                        type="text"
                        class="w-[70%] text-3xl p-2 text-center text-white"
                        style="background-color: #464AA2;"
                    />
                </Clock>
            </div>
            <div class="flex-1 flex justify-center">
                {#if imagePreviews.get(selectedQuestion.id) || selectedQuestion.imageUrl}
                    <button
                        type="button"
                        on:click={triggerFileInput}
                        class="w-[500px] h-[250px] rounded-lg overflow-hidden group flex justify-center items-center"
                    >
                        <img
                            src={imagePreviews.get(selectedQuestion.id) ||
                                selectedQuestion.imageUrl}
                            alt="Quiz question visual"
                            class="w-full h-full object-contain transition-transform duration-300 group-hover:scale-105"
                        />
                    </button>
                {:else}
                    <button
                        type="button"
                        on:click={triggerFileInput}
                        class="w-[500px] h-[250px] bg-white border-2 border-dashed border-gray-300 rounded-lg
                           flex flex-col justify-center items-center text-gray-500
                           hover:border-blue-500 hover:text-blue-600 transition-colors"
                    >
                        <svg class="w-12 h-12 mb-2"></svg>
                        <span class="font-bold">Upload Image</span>
                    </button>
                {/if}

                <div class="w-1/4"></div>

                <div class="w-1/4 flex items-center justify-end">
                    {#if errorForThisQuestion}
                        <div
                            class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-md"
                            role="alert"
                        >
                            <p class="font-bold">ข้อผิดพลาด:</p>
                            <p>{errorForThisQuestion}</p>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
        <div class="flex flex-wrap w-full">
            {#each selectedQuestion.choices as choice, i}
                <QuizChoiceCard color={COLORS[i]}>
                    <div
                        class="px-14 w-full h-full flex flex-col items-center justify-center gap-2"
                    >
                        <div class="w-full flex-grow min-h-0">
                            {#if choice.imageUrl}
                                <label
                                    class="w-full h-full cursor-pointer group flex items-center justify-center"
                                >
                                    <img
                                        src={choice.imageUrl}
                                        alt="Choice {i + 1}"
                                        class="w-full h-full object-contain pt-4 rounded-md transition-opacity group-hover:opacity-75"
                                    />
                                    <input
                                        type="file"
                                        class="hidden"
                                        accept="image/*"
                                        on:change={(e) =>
                                            uploadImageForChoice(e, i)}
                                    />
                                </label>
                            {:else}
                                <label
                                    class="w-full h-full flex items-center justify-center pt-4 bg-gray-200 text-gray-500 rounded-md cursor-pointer hover:bg-gray-300 hover:text-gray-600"
                                >
                                    Upload Image<input
                                        type="file"
                                        class="hidden"
                                        accept="image/*"
                                        on:change={(e) =>
                                            uploadImageForChoice(e, i)}
                                    />
                                </label>
                            {/if}
                        </div>
                        <div class="w-full flex-shrink-0 flex gap-2">
                            <input
                                class="rounded px-2 py-1 w-full text-black"
                                placeholder="Choice Text"
                                bind:value={selectedQuestion.choices[i].name}
                                on:change={notifyParentOfChange}
                            />
                            <input
                                type="checkbox"
                                class="w-16 h-18"
                                bind:checked={
                                    selectedQuestion.choices[i].correct
                                }
                                on:change={notifyParentOfChange}
                            />
                        </div>
                    </div>
                </QuizChoiceCard>
            {/each}
        </div>
    </div>
</div>
