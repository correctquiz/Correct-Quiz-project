<script lang="ts">
    import LoginForm from "../../lib/auth/LoginForm.svelte";
</script>

<div
    class="min-h-screen flex flex-col justify-start items-center pt-10 bg-[#FCFCFC]"
>
    <LoginForm userType="Player" />

    <p class="text-sm text-gray-700 mt-10">
        New player?
        <a
            href="/#/player/signup"
            class="font-bold text-blue-600 hover:underline">Create an account</a
        >
    </p>
</div>
