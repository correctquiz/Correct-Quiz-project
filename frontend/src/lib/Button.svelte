<button
    on:click
    class="text-white px-4 py-2 rounded shadow-md font-bold bg-[#464AA2] hover:bg-[#73759A] cursor-pointer disabled:bg-gray-400 disabled:cursor-not-allowed"
    {...$$props}
>
    <div>
        <slot />
    </div>
</button>
