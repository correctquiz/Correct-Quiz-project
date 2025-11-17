<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { Quiz, QuizQuestion } from "../../model/quiz";
    import { apiService } from "../../service/api";
    import Button from "../../lib/Button.svelte";
    import EditSidebar from "../../lib/edit/EditSidebar.svelte";
    import EditQuestion from "../../lib/edit/EditQuestion.svelte";
    import { push } from "svelte-spa-router";
    import { onMount } from "svelte";

    export let params: Record<string, string>;

    let quiz: Quiz | null;
    let selectedQuestion: QuizQuestion | null = null;
    let imagePreviews = new Map<number, string>();
    let quizNameError: string | null = null;
    let questionErrors = new Map<number, string>();
    const dispatch = createEventDispatcher();

    onMount(async () => {
        if (params && params["quizId"]) {
            const quizId = params["quizId"];
            if (!isNaN(parseInt(quizId))) {
                const quizIdAsNumber = parseInt(quizId);
                quiz = await apiService.getQuizById(quizIdAsNumber);
            }
        }
    });

    $: {
        if (quiz) {
            const newErrors = new Map<number, string>();

            if (!quiz.name.trim()) {
                quizNameError = "กรุณาตั้งชื่อชุดคำถาม";
            } else {
                quizNameError = null;
            }

            for (const question of quiz.questions) {
                if (!question.name.trim()) {
                    newErrors.set(question.id, "กรุณาตั้งชื่อคำถาม");
                    continue;
                }
                if (question.time < 20) {
                    newErrors.set(
                        question.id,
                        "ต้องกำหนดเวลาอย่างน้อย 20 วินาที",
                    );
                    continue;
                }

                if (question.time > 60) {
                    newErrors.set(
                        question.id,
                        "ต้องกำหนดเวลาไม่มากกว่า 60 วินาที",
                    );
                    continue;
                }

                if (isNaN(question.time)) {
                    newErrors.set(question.id, "กรุณากำหนดเวลาเป็นตัวเลข");
                    continue;
                }

                const hasCorrectChoice = question.choices.some(
                    (c) => c.correct,
                );
                if (!hasCorrectChoice) {
                    newErrors.set(
                        question.id,
                        "ต้องมีคำตอบที่ถูกต้องอย่างน้อย 1 ข้อ",
                    );
                    continue;
                }
            }
            questionErrors = newErrors;
        }
    }
    $: isSaveDisabled = quizNameError !== null || questionErrors.size > 0;

    async function onQuestionDelete() {
        if (!quiz || !selectedQuestion) return;

        const confirmed = confirm(
            `คุณแน่ใจหรือไม่ว่าจะลบคำถาม "${selectedQuestion.name}"?`,
        );
        if (!confirmed) return;

        const success = await apiService.deleteQuestion(selectedQuestion.id);

        if (success) {
            quiz.questions = quiz.questions.filter(
                (q) => q.id !== selectedQuestion!.id,
            );
            selectedQuestion = null;
            alert("ลบคำถามสำเร็จ!");
        } else {
            alert("เกิดข้อผิดพลาดในการลบคำถาม");
        }
    }

    function goBack() {
        push("/host");
    }

    async function save() {
        if (!quiz || isSaveDisabled) {
            alert("ไม่สามารถบันทึกได้ กรุณาแก้ไขข้อผิดพลาดทั้งหมดก่อน");
            return;
        }

        let shouldRefetch = false;

        quiz.questions.forEach((question) => {
            if (imagePreviews.has(question.id)) {
                question.imageUrl = imagePreviews.get(question.id);
                shouldRefetch = true;
            }
            if (question.id > 1000000000) {
                shouldRefetch = true;
            }
        });

        const success = await apiService.saveQuiz(quiz.id, quiz);

        if (success) {
            alert("บันทึกสำเร็จ!");
            imagePreviews.clear();

            if (shouldRefetch) {
                quiz = await apiService.getQuizById(quiz.id);
            }
        } else {
            alert("เกิดข้อผิดพลาดในการบันทึก");
        }
    }
</script>

{#if quiz != null}
    <div
        class="w-full p-2 flex justify-between items-center"
        style="background-color:#EEEEFD"
    >
        <button on:click={goBack} class="border-none cursor-pointer">
            <img
                src="../image/back.png"
                alt="back"
                class="w-10 h-10 hover:opacity-80"
            />
        </button>
        <div class="flex gap-2 pr-2">
            <input
                type="text"
                class="rounded px-3 shadow-inner"
                style="background-color: #F87923"
                placeholder="Quiz name"
                bind:value={quiz.name}
            />
            <div></div>
            <Button on:click={save} disabled={isSaveDisabled}>Save</Button>
        </div>
    </div>
    <div class="flex">
        <EditSidebar
            bind:questions={quiz.questions}
            bind:selectedQuestion
            invalidQuestionIds={new Set(questionErrors.keys())}
        />
        {#if selectedQuestion != null}
            <EditQuestion
                on:delete={onQuestionDelete}
                on:change={() => (quiz = quiz)}
                bind:selectedQuestion
                {imagePreviews}
                errorForThisQuestion={questionErrors.get(selectedQuestion.id)}
            />
        {/if}
    </div>
{:else}
    quiz not found
{/if}
