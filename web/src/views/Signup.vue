<script setup>
  import { ref } from 'vue';

  const name = ref('');
  const password = ref('');

  const submitForm = async () => {
    const response = await fetch('/api/signup', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: name.value, password: password.value }),
    });

    if (response.ok) {
      alert('Signup successful!');
    } else {
      alert('Signup failed. Please try again.');
    }
  };
</script>

<template>

  <div>

    <div class="flex flex-col items-center text-center justify-center 
                  h-screen">

      <div class="card-container">

        <div class="card-content">

            <h1 class="title">Sign Up</h1>
            <h2 class="title-secondary">Welcome</h2>

            <form @submit.prevent="submitForm" class="flex flex-col gap-4">
              <p class="flex flex-col">
                <label for="name" class="f-label">Name</label>
                <input id="name" v-model="name" type="text" required class="f-input"/>
              </p>

              <p class="flex flex-col">
                <label for="password" class="f-label">Password</label>
                <input id="password" v-model="password" type="password" required class="f-input"/>
              </p>
              
              <div class="relative inline-block mt-5">
                <button type="submit" class="taped-button">Submit</button> 
                
                <div class="tape tape-top"></div>
                <div class="tape tape-bottom"></div>
              </div>
              
              
            </form>

          </div>

        </div>

      
    </div>

  </div>
</template>
