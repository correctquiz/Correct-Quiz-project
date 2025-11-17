<script lang="ts">
    export let color: string;
    export let size: "large" | "small" = "large";
    const largeClasses = "w-full md:w-1/2 h-60 text-white font-bold flex pl-4 items-center text-3xl";
    const smallClasses = "w-fit h-fit h-40 text-white font-bold flex pr-20 items-center text-xl";
</script>

<div class="{color} {size === 'large' ? largeClasses : smallClasses}">
    <slot />
</div>