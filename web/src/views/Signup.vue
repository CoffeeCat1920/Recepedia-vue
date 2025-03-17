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
    <h1>Sign Up</h1>
    <h2>Welcome</h2>

    <form @submit.prevent="submitForm">
      <p>
        <label for="name">Name</label>
        <br>
        <input id="name" v-model="name" type="text" required />
      </p>

      <p>
        <label for="password">Password</label>
        <br>
        <input id="password" v-model="password" type="password" required />
      </p>

      <p>
        <button type="submit">Submit</button>
      </p>
    </form>
  </div>
</template>
