<script setup>
  import { ref } from 'vue'; 

  const name = ref('');
  const content = ref('');

  const submitForm = async () => {
    const response = await fetch('/api/uploadrecipe', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: name.value, content: content.value}),
    });   

    if (response.ok) {
      alert('Uploaded Recipe');
    } else {
      alert('Can\'t Upload Recipe');
    } 
  };  
</script>

<template>

  <div class="flex flex-col text-center items-start justify-center pl-12 h-screen">

    <div class="card-container">

      <div class="card-content">

        <h1 class="title">Upload a Recipe</h1>
        <h2 class="title-secondary">What looks yummy?</h2>

        <form @submit.prevent="submitForm" class="flex flex-col gap-4">
          <p class="flex flex-col">
            <label for="name" class="f-label">Name</label>
            <br>
            <input v-model="name" type="text" name="name" class="f-input" required />
          </p>

          <p class="flex flex-col">
            <label for="content" class="f-label">Write here</label>
            <br>
            <textarea v-model="content" name="content" rows="10" class="border-2 border-dashed border-amber-700 font-semibold  text-amber-950" required></textarea>
          </p>

          <div class="relative inline-block mt-10">
            <button type="submit" class="taped-button">Submit</button>
            <div class="tape tape-top"></div>
            <div class="tape tape-bottom"></div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
