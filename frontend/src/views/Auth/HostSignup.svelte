<script lang="ts">
    import SignupForm from "../../lib/auth/SignupForm.svelte";
</script>

<div
    class="min-h-screen flex flex-col justify-center items-center bg-[#FCFCFC]"
>
    <SignupForm userType="Host" />
    <p class="text-sm text-gray-700 mt-10">
        Already have an account?
        <a href="/#/host/login" class="font-bold text-blue-600 hover:underline"
            >login</a
        >
    </p>
</div>
