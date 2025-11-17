<script lang="ts">
    import type { QuizQuestion } from "../../model/quiz";

    export let question: QuizQuestion;
    export let selectedQuestion: QuizQuestion | null;

    export let isInvalid: boolean = false;
    
    $: isSelected = selectedQuestion?.id === question.id;

    function onClick() {
        selectedQuestion = question;
    }
</script>

<button
    class="bg-white rounded-md p-2 border hover:bg-gray-50"
    on:click={onClick}
    class:border-red-500={isInvalid}
>
    {question.name}
</button>
