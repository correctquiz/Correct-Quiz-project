<script lang="ts">
    import QuizCard from "../../lib/QuizCard.svelte";
    import type { Quiz } from "../../model/quiz";
    import { apiService } from "../../service/api";
    import { location, push } from "svelte-spa-router";
    import Button from "../../lib/Button.svelte";
    import { logout, userStore, isLoggingOut } from "../../service/userStore";
    import { game } from "../../service/host/host";

    let quizzes: Quiz[] = [];
    let isLoading = true;
    game.navigate = push;

    async function handleHostLogout() {
        isLoggingOut.set(true);
        push("/host/login");
        await logout();
    }

    function onHost(event: { detail: Quiz }) {
        const quizIdString = String(event.detail.id);
        game.hostQuiz(quizIdString);
    }

    async function loadQuizzes() {
        if (!$userStore.loggedIn) return;
        isLoading = true;
        try {
            quizzes = await apiService.getQuizzes();
        } catch (error) {
            alert("Could not load your quizzes.");
        } finally {
            isLoading = false;
        }
    }

    async function onQuizHost(event: { detail: Quiz }) {
        const quizToHost = event.detail;

        if (quizToHost.questions.length < 1) {
            return;
        }

        try {
            await game.hostQuiz(String(quizToHost.id));

            push("/host/game");
        } catch (error) {
            alert("Could not start hosting the game.");
        }
    }

    async function createNewQuiz() {
        const newQuizName = prompt("กรุณาตั้งชื่อ ชุดคำถาม:", "New Quiz");
        if (newQuizName) {
            const newQuiz = await apiService.createQuiz(newQuizName);
            if (newQuiz) {
                push(`/edit/${newQuiz.id}`);
            }
        }
    }
    async function onQuizDelete(event: { detail: Quiz }) {
        const quizToDelete = event.detail;
        const confirmed = confirm(
            `คุณแน่ใจว่าจะลบชุดคำถาม"${quizToDelete.name}" ไหม ?`,
        );

        if (confirmed) {
            const success = await apiService.deleteQuiz(quizToDelete.id);
            if (success) {
                quizzes = quizzes.filter((q) => q.id !== quizToDelete.id);
                alert("ลบชุดคำถามสำเร็จ!");
            } else {
                alert("เกิดข้อผิดพลาดในการลบคำถาม");
            }
        }
    }

    $: if ($location === "/host" && $userStore.loggedIn) {
        loadQuizzes();
    }
</script>

<div class="p-8 min-h-screen bg-[#FCFCFC]">
    <div class="flex justify-between items-center">
        <h3 class="text-4xl text-[#F87923] font-bold">Your Quizzes</h3>
        <div class="flex items-center gap-4">
            <img
                src="../image/CorrectQuiz.png"
                alt="Quiz icon"
                class="w-26 h-28"
                style="filter: drop-shadow(4px 4px 2px rgba(0,0,0,0.2));"
            />
            <button
                on:click={handleHostLogout}
                class="text-red-600 font-bold text-sm cursor-pointer underline"
                >Logout</button
            >
        </div>
    </div>
    <div class="flex flex-col gap-2 mt-4">
        {#each quizzes as quiz (quiz.id)}
            <QuizCard
                on:host={onQuizHost}
                on:delete={onQuizDelete}
                {quiz}
                isPlayable={quiz.questions?.length >= 1}
            />
        {/each}
    </div>
    <div class="flex justify-center mt-4">
        <Button on:click={createNewQuiz}>Create New</Button>
    </div>
</div>
