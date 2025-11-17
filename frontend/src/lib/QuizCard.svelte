<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import Button from "./Button.svelte";
  import type { Quiz } from "../model/quiz";
  import { push } from "svelte-spa-router";

  const dispatch = createEventDispatcher();

  export let quiz: Quiz;

  export let isPlayable: boolean;

  function edit() {
    push(`/edit/${quiz.id}`);
  }
</script>

<div
  class="flex justify-between items-center border p-4 rounded-xl"
  style="background-color: #F87923;"
>
  <p>{quiz.name}</p>
  <div class="flex gap-2 items-center">
    <Button
      on:click={() => dispatch("host", quiz)}
      disabled={!isPlayable}
      
      title="{isPlayable
        ? 'เริ่มชุดคำถาม'
        : 'ชุดคำถามต้องมีอย่างน้อย 2 ข้อ'}">Host</Button
    >
    <Button on:click={edit}>Edit</Button>
    <button
      on:click={() => dispatch("delete", quiz)}
      class="text-white px-4 py-2 rounded shadow-md font-bold bg-[#FF0000] hover:bg-[#FF5050] cursor-pointer"
      >Delete</button
    >
  </div>
</div>
