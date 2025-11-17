<script lang="ts">
    import { onMount } from 'svelte';
    import { querystring, push } from 'svelte-spa-router';
    import { apiService } from '../../service/api';

    let message = "กำลังยืนยันบัญชีของคุณ...";
    let isError = false;
    let isSuccess = false; 

    onMount(async () => {
        const params = new URLSearchParams($querystring);
        const token = params.get('token');

        if (!token) {
            message = "ไม่พบ Token ใน URL";
            isError = true;
            return;
        }

        try {
            await apiService.verifyEmailToken(token);
            message = "ยืนยันอีเมลสำเร็จ! บัญชีของคุณพร้อมใช้งานแล้ว";
            isSuccess = true; 
            

        } catch (err: any) {
            isError = true;
            message = err.response?.data?.error || "การยืนยันล้มเหลว";
        }
    });

    function goToLogin() {
        push('/player/login'); 
    }
</script>

<div class="flex flex-col justify-center items-center min-h-screen bg-[#FCFCFC] gap-6">
    <h1 class="text-3xl font-bold text-center" 
        class:text-red-500={isError} 
        class:text-green-600={isSuccess}
        class:text-gray-700={!isError && !isSuccess}
    >
        {message}
    </h1>

    {#if isSuccess}
        <div class="text-center">
            <p class="mb-4 text-gray-600">คุณสามารถปิดหน้านี้ แล้วกลับไปเข้าสู่ระบบที่หน้าเดิมได้เลย</p>
            <button 
                on:click={goToLogin}
                class="px-6 py-3 text-white bg-blue-600 rounded-lg shadow hover:bg-blue-700 transition-all font-semibold"
            >
                ไปหน้าเข้าสู่ระบบ
            </button>
        </div>
    {/if}
</div>