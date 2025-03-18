<script>
  import { ref } from 'vue'
  import { useAuthStore } from '@/stores/authStore.ts'

  export default {

    setup() {

      const name = ref('');
      const password = ref('');

      const auth = useAuthStore();

      const login = async () => {
        await auth.login(name.value, password.value);
      }; 

      return { auth, name, password, login };
    }, 
    async mount() {
      // await this.auth.checkLogin(); 
    }
  }


</script>

<template>

  <h1>Login</h1>
  <h2>Welcome back</h2>

  <form @submit.prevent="login">

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
