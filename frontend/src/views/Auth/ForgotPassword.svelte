<script lang="ts">
    import { sendPasswordResetEmail } from "firebase/auth";
    import { auth } from "../../service/firebase";

    let email = "";
    let message = "";
    let isError = false;
    let isLoading = false;

    export let userType: "Host" | "Player";

    async function handleResetPassword() {
        if (!email) {
            isError = true;
            message = "Please enter your email address.";
            return;
        }

        isLoading = true;
        isError = false;
        message = "";

        try {
            await sendPasswordResetEmail(auth, email);
            message = "A password reset link has been sent to your email.";
        } catch (error: any) {
            isError = true;
            message = error.message;
        } finally {
            isLoading = false;
        }
    }
</script>

<div
    class="min-h-screen flex flex-col justify-center items-center bg-[#F2F2F2]"
>
    <div
        class="p-8 rounded-xl shadow-lg w-full max-w-md"
        style="background-color:#EEEEFD;"
    >
        <h2 class="text-3xl font-bold mb-6 text-center text-gray-800">
            Forgot Password
        </h2>

        <form on:submit|preventDefault={handleResetPassword}>
            <div class="mb-4">
                <label for="email" class="block text-gray-700 font-medium mb-2"
                    >Email</label
                >
                <input
                    type="email"
                    id="email"
                    bind:value={email}
                    class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    style="background-color: #F87923;"
                    placeholder="Enter your registered email"
                    disabled={isLoading}
                />
            </div>

            {#if message}
                <p
                    class:text-red-500={isError}
                    class:text-green-600={!isError}
                    class="text-center mb-4"
                >
                    {message}
                </p>
            {/if}

            <button
                type="submit"
                class="w-full text-white py-2 rounded-lg hover:bg-orange-600 transition-colors disabled:bg-gray-400 cursor-pointer"
                style="background-color: #464AA2;"
                disabled={isLoading}
            >
                {isLoading ? "Sending..." : "Send Reset Link"}
            </button>
        </form>

        <div class="text-center mt-4">
            {#if userType === "Player"}
                <a href="/#/player/login" class="text-blue-600 hover:underline"
                    >Back to Login</a
                >
            {:else}
                <a href="/#/host/login" class="text-blue-600 hover:underline"
                    >Back to Login</a
                >
            {/if}
        </div>
    </div>
</div>
