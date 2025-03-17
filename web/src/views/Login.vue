<script setup>
  import { ref } from 'vue';

  const name = ref('');
  const password = ref('');

  const submitForm = async () => {
    const response = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: name.value, password: password.value }),
    })

    if (response.ok) {
      alert('Login successful!');
    } else {
      alert('Login failed. Please try again.');
    }
  }  
</script>

<template>

  <h1>Login</h1>
  <h2>Welcome back</h2>

  <form @submit.prevent="submitForm">

    <p>
      <label for="name">Name</label>
      <br>
      <input id="name" v-model="name" type="text" required /> 
    </p>

    <p>
      <label for="password">Password</label>
      <br>
      <input id="password" v-model="password" type="password" /> 
    </p>

    <p>
      <button type="submit">Submit</button>
    </p>
    
  </form>

</template>
