<script lang="ts">
    import Button from "../../lib/Button.svelte";
    import { userStore, logout } from "../../service/userStore";
    import { playerGameStore } from "../../service/gameStore";
    import { push, querystring } from "svelte-spa-router";
    import { BASE_URL } from "../../service/api";
    import { getHeaders } from "../../service/api";

    let code: string = "";
    let name: string = "";

    let joinError: string | null = null;
    let isLoading = false;
    let codeProcessed = false;

    const hasspace = /\s/;

    const forbiddenWords = [
        "admin",
        "god",
        "fuck",
        "shit",
        "dick",
        "hitler",
        "pussy",
        "penis",
        "ass",
        "ควย",
        "หี",
        "เย็ด",
        "สัส",
        "แม่ง",
        "เหี้ย",
    ];

    async function handlePlayerLogout() {
        await logout();
        push("/");
    }

    async function join() {
        joinError = null;

        if (!name || !code) {
            joinError = "กรุณากรอกชื่อและรหัส PIN";
            return;
        }

        if (/[^0-9]/.test(code)) {
            joinError = "Game pin ต้องเป็นตัวเลขเท่านั้น";
            return;
        }

        if (code.length < 6 || code.length > 6) {
            joinError = "Game pin ต้องมี 6 หลัก";
            return;
        }

        if (hasspace.test(name)) {
            joinError = "ชื่อต้องไม่มีเว้นวรรค";
            return;
        }

        if (forbiddenWords.some((word) => name.toLowerCase().includes(word))) {
            joinError = "ชื่อไม่เหมาะสมกรุณากรอกชื่อใหม่";
            return;
        }

        if (/^[0-9]/.test(name)) {
            joinError = "ชื่อห้ามเริ่มด้วยตัวเลข";
            return;
        }

        isLoading = true;

        try {
            let token = localStorage.getItem("jwt_token");

            if (!token) {
                const guestResponse = await fetch(
                    `${BASE_URL}/api/auth/guest-login`,
                    {
                        method: "POST",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify({ name: name }),
                    },
                );

                if (!guestResponse.ok)
                    throw new Error("Failed to create guest user");

                const guestData = await guestResponse.json();
                token = guestData.token;

                localStorage.setItem("jwt_token", token);
            }

            const pinCheckResponse = await fetch(`${BASE_URL}/api/game/check`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ code: code }),
            });
            if (!pinCheckResponse.ok) {
                throw new Error("Invalid game PIN");
            } else if (!pinCheckResponse.ok) {
                throw new Error("Invalid game PIN");
            }
            const nameCheckResponse = await fetch(
                `${BASE_URL}/api/game/check-name`,
                {
                    method: "POST",
                    headers: getHeaders(),
                    body: JSON.stringify({ code: code, name: name }),
                },
            );
            if (nameCheckResponse.status === 409) {
                throw new Error("Username already taken");
            } else if (!nameCheckResponse.ok) {
                throw new Error("Error checking name");
            }

            $playerGameStore.join(code, name, token);
        } catch (error: any) {
            if (error.message === "Invalid game PIN") {
                joinError = "รหัส PIN ของเกมไม่ถูกต้อง กรุณาลองใหม่";
            } else if (error.message === "Username already taken") {
                joinError = "ชื่อนี้ใช้ไปแล้ว";
            } else {
                joinError = "เกิดข้อผิดพลาดที่ไม่ทราบสาเหตุ";
            }
        } finally {
            isLoading = false;
        }
    }

    $: {
        if ($querystring && !codeProcessed) {
            const params = new URLSearchParams($querystring);
            const codeFromUrl = params.get("code");
            if (codeFromUrl) {
                code = codeFromUrl;
                codeProcessed = true;
            }
        }
    }
</script>

<div class="min-h-screen flex justify-center bg-[#FCFCFC]">
    <div class="items-center mt-23">
        <img
            src="../image/CorrectQuiz.png"
            alt="Quiz icon"
            class="w-52 h-53 mx-auto"
            style=" filter: drop-shadow(4px 4px 2px rgba(0,0,0,0.2));"
        />
        <div class="flex flex-col gap-4 item-center">
            <input
                bind:value={code}
                type="text"
                placeholder="Game code"
                class="p-2 font-bold rounded-lg w-65 max-w-xs bg-[#F87923]"
                style="box-shadow: inset 0 2px 4px 0 rgba(0, 0, 0, 0.2);"
            />
            <input
                bind:value={name}
                type="text"
                placeholder="Name"
                class="p-2 font-bold rounded-lg w-65 max-w-xs bg-[#F87923]"
                style="box-shadow: inset 0 2px 4px 0 rgba(0, 0, 0, 0.2);"
            />
            {#if joinError}
                <p class="text-red-500 font-bold text-center">{joinError}</p>
            {/if}
            <div class="flex justify-center mt-2 mb-18">
                <Button on:click={join}>Join game</Button>
            </div>
            {#if !$userStore.loggedIn}
                <div class="text-center border-t mt-6 pt-4">
                    <p class="text-sm text-gray-700">
                        Already have an account?
                        <a
                            href="/#/player/login"
                            class="font-bold text-blue-600 hover:underline"
                            >Login</a
                        >
                    </p>
                    <p class="text-sm text-gray-700">
                        New player?
                        <a
                            href="/#/player/signup"
                            class="font-bold text-blue-600 hover:underline"
                            >Create an account</a
                        >
                    </p>
                </div>
            {/if}
            {#if $userStore.loggedIn}
                <button
                    on:click={handlePlayerLogout}
                    class="text-red-600 font-bold text-sm pt-15 cursor-pointer underline"
                >
                    Logout
                </button>
            {/if}
        </div>
    </div>
</div>
