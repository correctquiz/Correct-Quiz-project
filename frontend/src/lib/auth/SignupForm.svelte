<script lang="ts">
    import { FirebaseError } from "firebase/app";
    import { apiService } from "../../service/api";
    import { createUserWithEmailAndPassword,deleteUser, } from "firebase/auth";
    import { auth } from "../../service/firebase";
    import { isSigningUp } from "../../service/userStore";

    let email = "";
    let username = "";
    let password = "";
    let confirmPassword = "";
    let errorMessage = "";
    let isLoading = false;
    export let userType: "Host" | "Player";

    const specialChars = /[!@#$%^&*()_+=\[\]{};':"\\|,.<>\/?~`]/;

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

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

    async function handleSignup() {
        errorMessage = "";
        isSigningUp.set(true);
        isLoading = true;
        let userCredential = null;

        try {
            if (!email || !username || !password) {
                errorMessage = "กรุณากรอกอีเมล ชื่อผู้ใช้ และรหัสผ่าน";
                throw new Error(errorMessage);
            }
            if (!emailRegex.test(email)) {
                errorMessage =
                    "โปรดระบุรูปแบบอีเมลที่ถูกต้อง (เช่น Correct@quiz.com)";
                throw new Error(errorMessage);
            }

            if (password !== confirmPassword) {
                errorMessage = "รหัสไม่ตรงกัน";
                throw new Error(errorMessage);
            }

            let firebaseUid = "";
            try {
                userCredential = await createUserWithEmailAndPassword(
                    auth,
                    email,
                    password,
                );
                firebaseUid = userCredential.user.uid;
            } catch (error) {
                if (
                    error instanceof FirebaseError &&
                    error.code === "auth/email-already-in-use"
                ) {
                    errorMessage = "อีเมลนี้ถูกใช้งานแล้ว กรุณาใช้อีเมลอื่น";
                } else {
                    errorMessage = "เกิดข้อผิดพลาดในการสร้างบัญชี Firebase";
                }
                throw error;
            }

            try {
                await apiService.register({
                    firebaseUid: firebaseUid,
                    username: username,
                    email: email,
                    role: userType.toLowerCase(),
                });
            } catch (backendError: any) {
                if (userCredential) {
                    await deleteUser(userCredential.user); 
                }
                if (backendError.response?.status === 409) {
                    errorMessage =
                        backendError.response.data.error ||
                        "ชื่อผู้ใช้หรืออีเมลนี้ถูกใช้งานแล้ว";
                } else {
                    errorMessage =
                        "เกิดข้อผิดพลาดในการตรวจสอบข้อมูลกับเซิร์ฟเวอร์";
                }
                throw backendError;
            }

            alert(
                "การสมัครสำเร็จ! กรุณาตรวจสอบอีเมลและคลิกลิงก์ยืนยันก่อนเข้าสู่ระบบ",
            );

        } catch (error: any) {
            if (errorMessage === "") {
                errorMessage = "เกิดข้อผิดพลาดไม่ทราบสาเหตุ กรุณาลองใหม่";
            }
            console.error("Signup failed:", error);
        } finally {
            isLoading = false;
            isSigningUp.set(false);
        }
    }
</script>

<h2
    class="text-3xl font-bold text-center text-white mb-4 mt-10"
    style="text-shadow:  -1px -1px 0 #000,  1px -1px 0 #000,-1px  1px 0 #000,1px  1px 0 #000;"
>
    Create Account ({userType})
</h2>

<div
    class="px-10 py-3 rounded-xl shadow-lg w-full max-w-md"
    style="background-color: #EEEEFD;"
>
    <form on:submit|preventDefault={handleSignup}>
        <div>
            <label for="email" class="block text-gray-700 font-medium"
                >Email</label
            >
            <input
                type="email"
                id="email"
                bind:value={email}
                class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 shadow-inner"
                style="background-color:#F87923;"
                placeholder="Enter your email address"
                disabled={isLoading}
            />
        </div>

        <div>
            <label for="username" class="block text-gray-700 font-medium mb-2"
                >Username</label
            >
            <input
                type="text"
                id="username"
                bind:value={username}
                class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 shadow-inner"
                style="background-color:#F87923;"
                placeholder="Username"
                disabled={isLoading}
            />
        </div>

        <div>
            <label for="password" class="block text-gray-700 font-medium mb-2"
                >Password</label
            >
            <input
                type="password"
                id="password"
                bind:value={password}
                class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 shadow-inner"
                style="background-color:#F87923;"
                placeholder="********"
                disabled={isLoading}
            />
            <ul class="text-[10px] text-gray-500 list-disc list-inside mt-1">
                <li>อย่างน้อย 6 ตัวอักษร</li>
                <li>มีตัวพิมพ์ใหญ่อย่างน้อย 1 ตัว (A-Z)</li>
                <li>รหัสผ่านต้องมีอักขระพิเศษอย่างน้อย 1 ตัว</li>
                <li>มีตัวเลขอย่างน้อย 1 ตัว (0-9)</li>
            </ul>
        </div>

        <div class="mb-5">
            <label
                for="confirm-password"
                class="block text-gray-700 font-medium mb-2"
                >Confirm Password</label
            >
            <input
                type="password"
                id="confirm-password"
                bind:value={confirmPassword}
                class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 shadow-inner"
                style="background-color:#F87923;"
                placeholder="********"
                disabled={isLoading}
            />
        </div>

        {#if errorMessage}
            <p class="text-red-500 text-center mb-4">{errorMessage}</p>
        {/if}

        <button
            type="submit"
            class="w-full text-white py-2 rounded-lg transition-colors disabled:bg-gray-400 cursor-pointer"
            style="background-color: #464AA2;"
            disabled={isLoading}
        >
            {isLoading ? "Creating Account..." : "Sign Up"}
        </button>
    </form>
</div>
