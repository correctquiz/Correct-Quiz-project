<script lang="ts">
    import { apiService } from "../../service/api";
    import { auth } from "../../service/firebase";
    import { signInWithEmailAndPassword, signOut } from "firebase/auth";
    import { isLoggingIn, userStore } from "../../service/userStore";
    import { push } from "svelte-spa-router";

    export let userType: "Host" | "Player";

    let identifier = "";
    let password = "";
    let errorMessage = "";
    let isLoading = false;

    let emailForResend: string | null = null;

    async function handleLogin() {
        isLoggingIn.set(true);
        isLoading = true;
        errorMessage = "";
        emailForResend = null;

        try {

            if (!identifier || !password) {
                errorMessage = "กรุณาใส่ อีเมลหรือชื่อผู้ใช้และรหัสผ่าน";
                return; 
            }

            let emailToLogin = identifier;
            if (!identifier.includes("@")) {
                const foundEmail =
                    await apiService.getEmailForUsername(identifier);
                if (!foundEmail) {
                    errorMessage = "ไม่พบชื่อผู้ใช้นี้";
                    return; 
                }
                emailToLogin = foundEmail;
            } 

            emailForResend = emailToLogin;

            const userCredential = await signInWithEmailAndPassword(
                auth,
                emailToLogin,
                password,
            );

            const user = userCredential.user;

            if (!user.emailVerified) {
                errorMessage =
                    "กรุณายืนยันอีเมลของคุณก่อนเข้าสู่ระบบ ตรวจสอบกล่องจดหมายสำหรับลิงก์ยืนยัน";
                return; 
            } 

            emailForResend = null;

            const idTokenResult = await user.getIdTokenResult();
            const actualRole = idTokenResult.claims.role as "host" | "player";
            const intendedRole = userType.toLowerCase();
            if (actualRole === "player" && intendedRole === "host") {
                errorMessage =
                    "บัญชีผู้เล่น Player ไม่สามารถเข้าสู่ระบบ Host ได้";
                await signOut(auth);
                return;
            }
            const result = await apiService.login({
                idToken: idTokenResult.token,
            });
            if (result && result.token) {
                localStorage.setItem("jwt_token", result.token);
                userStore.set({
                    loggedIn: true,
                    user: user,
                    userType: actualRole === "host" ? "Host" : "Player",
                });
                if (actualRole === "host" && intendedRole === "host") {
                    push("/host");
                } else {
                    push("/");
                }
            } else {
                errorMessage = "การยืนยันตัวตนฝั่ง Backend ล้มเหลว";
            }
        } catch (error: any) {
            if (error.code === "auth/email-not-verified") {
                errorMessage =
                    "กรุณายืนยันอีเมลของคุณก่อนเข้าสู่ระบบ ตรวจสอบกล่องจดหมายสำหรับลิงก์ยืนยัน";
            } else {
                errorMessage =
                    "อีเมลหรือรหัสผ่านไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง";
            }
            console.error("Login error:", error);
        } finally {
            isLoading = false;
            isLoggingIn.set(false);
        }
    }

    async function resendVerificationEmail() {
        if (!emailForResend) {
            errorMessage = "เกิดข้อผิดพลาด: ไม่พบอีเมลสำหรับส่งซ้ำ";
            return;
        }
        isLoading = true;
        errorMessage = "";
        try {
            await apiService.resendVerificationEmail(emailForResend);
            alert(
                "ส่งอีเมลยืนยัน (5 นาที) อีกครั้งแล้ว กรุณาตรวจสอบกล่องจดหมาย (รวมถึงสแปม)",
            );
        } catch (error: any) {
            console.error("Error resending verification email:", error);
            errorMessage =
                error.response?.data?.error ||
                "เกิดข้อผิดพลาดในการส่งอีเมลยืนยันซ้ำ";
        } finally {
            isLoading = false;
        }
    }
</script>

<h2
    class="text-4xl font-bold mt-4 mb-8 text-center text-white"
    style="text-shadow:  -1px -1px 0 #000,  1px -1px 0 #000,-1px  1px 0 #000,1px  1px 0 #000;"
>
    Login ({userType})
</h2>

<div
    class="bg-white p-8 rounded-xl shadow-lg w-full max-w-md"
    style="background-color: #EEEEFD;"
>
    <form on:submit|preventDefault={handleLogin}>
        <div class="mb-4">
            <label
                for="confirm-password"
                class="block text-gray-700 font-medium mb-2"
                >Email or Username</label
            >
            <input
                type="text"
                bind:value={identifier}
                class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 shadow-inner"
                style="background-color:#F87923;"
                disabled={isLoading}
                placeholder="Enter your email or username"
            />
        </div>

        <div class="mb-10">
            <label
                for="confirm-password"
                class="block text-gray-700 font-medium mb-2">Password</label
            >
            <input
                type="password"
                bind:value={password}
                class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 shadow-inner"
                style="background-color:#F87923;"
                disabled={isLoading}
                placeholder="********"
            />
        </div>

        {#if errorMessage}
            <p class="text-red-500 text-center mb-2">{errorMessage}</p>
            {#if errorMessage.includes("ยืนยันอีเมล")}
                <button
                    type="button"
                    on:click={resendVerificationEmail}
                    class="w-full text-sm text-blue-600 hover:underline text-center mb-4 disabled:text-gray-400 disabled:no-underline"
                    disabled={!emailForResend || isLoading}
                >
                    (ส่งอีเมลยืนยันอีกครั้ง)
                </button>
            {/if}
        {/if}

        <button
            type="submit"
            class="w-full text-white py-2 rounded-lg transition-colors disabled:bg-gray-400 cursor-pointer"
            style="background-color:#464AA2;"
            disabled={isLoading}
        >
            {isLoading ? "Logging in..." : "Log In"}
        </button>
        <div class="text-center mt-4">
            {#if userType === "Player"}
                <a
                    href="/#/player/forgot-password"
                    class="text-sm text-blue-600 hover:underline"
                >
                    Forgot Password?
                </a>
            {:else if userType === "Host"}
                <a
                    href="/#/host/forgot-password"
                    class="text-sm text-blue-600 hover:underline"
                >
                    Forgot Password?
                </a>
            {/if}
        </div>
    </form>
</div>
