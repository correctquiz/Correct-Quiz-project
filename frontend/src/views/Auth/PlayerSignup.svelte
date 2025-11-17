<script lang="ts">
    import SignupForm from "../../lib/auth/SignupForm.svelte";
</script>

<div
    class="min-h-screen flex flex-col justify-start items-center pt-6 bg-[#FCFCFC]"
>
    <SignupForm userType="Player" />
    <p class="text-sm text-gray-700 pt-4">
        Already have an account?
        <a
            href="/#/player/login"
            class="font-bold text-blue-600 hover:underline">Login</a
        >
    </p>
</div>
