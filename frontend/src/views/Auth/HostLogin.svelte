<script lang="ts">
    import LoginForm from "../../lib/auth/LoginForm.svelte";
</script>

<div
    class="min-h-screen flex flex-col justify-center items-center bg-[#FCFCFC]"
>
    <LoginForm userType="Host" />
    <p class="text-sm text-gray-700 mt-10">
        New Host?
        <a href="/#/host/signup" class="font-bold text-blue-600 hover:underline"
            >Create an account</a
        >
    </p>
</div>
