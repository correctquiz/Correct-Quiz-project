<script lang="ts">
    import type { QuizQuestion } from "../../model/quiz";
    import Button from "../Button.svelte";
    import SidebarItem from "./SidebarItem.svelte";

    export let questions: QuizQuestion[];
    export let selectedQuestion: QuizQuestion | null;
    export let invalidQuestionIds: Set<number>;

    function addNew() {
        const tempId = Date.now();

        const newQuestion: QuizQuestion = {
            id: tempId,
            name: "New Question",
            time: 60,
            imageUrl: "",
            choices: [
                { id: tempId + 1, name: "", correct: false },
                { id: tempId + 2, name: "", correct: false },
                { id: tempId + 3, name: "", correct: false },
                { id: tempId + 4, name: "", correct: false },
            ],
            index: questions.length,
            correctAnswerIndex: -1,
        };

        questions = [...questions, newQuestion];
    }
</script>

<div
    class="min-h-screen min-w-64 p-2 flex flex-col gap-2"
    style="background-color: #EEEEFD;"
>
    {#each questions as question}
        <SidebarItem
            {question}
            bind:selectedQuestion
            isInvalid={invalidQuestionIds.has(question.id)}
        />
    {/each}
    <Button on:click={addNew}>Add new</Button>
</div>
