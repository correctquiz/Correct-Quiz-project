<script lang="ts">
    import { push } from "svelte-spa-router";
    import Leaderboard from "../../lib/Learderboard.svelte";
    import Button from "../../lib/Button.svelte";
    import {
        leaderboard,
        gameCode,
        game,
        isHostNavigating,
    } from "../../service/host/host";
    import { BASE_URL } from "../../service/api";
    import { onDestroy } from "svelte";

    let isNavigatingBack = false;

    function goBack() {
        isNavigatingBack = true;
        isHostNavigating.set(true);
        game.unhost({ broadcastEnd: false });
        push("/host");
    }

    function downloadResults() {
        if (!$gameCode) {
            alert("ไม่พบรหัสเกมสำหรับ Export");
            return;
        }
        const downloadUrl = `${BASE_URL}/api/games/${$gameCode}/export/csv`;

        window.open(downloadUrl, "_blank");
    }

    onDestroy(() => {
        if (!isNavigatingBack) {
            game.unhost();
        }
    });
</script>

<div class="relative min-h-screen w-full bg-[#FCFCFC]">
    <div>
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
        <div class="absolute top-8 right-8 z-10">
            <Button on:click={downloadResults}>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5 inline mr-2"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                >
                    <path
                        fill-rule="evenodd"
                        d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
                        clip-rule="evenodd"
                    />
                </svg>
                Export ผลลัพธ์ (CSV)
            </Button>
        </div>
    </div>
    <div class="flex flex-col items-center min-h-screen w-full">
        <div class="mt-32">
            <h2
                class="text-center text-white text-5xl font-bold"
                style="text-shadow:  -1px -1px 0 #000,  1px -1px 0 #000,-1px  1px 0 #000,1px  1px 0 #000;"
            >
                Game ended!
            </h2>
            <div class="flex flex-wrap gap-2 mt-10">
                <Leaderboard finish={true} leaderboard={$leaderboard} />
            </div>
        </div>
    </div>
</div>
