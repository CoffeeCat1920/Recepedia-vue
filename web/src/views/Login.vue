<script>
  import { ref } from 'vue'
  import { useAuthStore } from '@/stores/authStore.ts'
  import { useRouter } from "vue-router";

  export default {

    setup() {
      
      const router = useRouter();

      const name = ref('');
      const password = ref('');

      const auth = useAuthStore();

      const login = async () => {
        await auth.login(name.value, password.value, router);
      }; 

      return { auth, name, password, login };
    }, 
    async mount() {
      // await this.auth.checkLogin(); 
    }
  }


</script>

<template>

  <div class="flex flex-col items-center text-center justify-center h-screen">

    <div class="card-container">

      <div class="card-content">

        <h1 class="title">Log In</h1>
        <h2 class="title-secondary">Welcome Back</h2>

        <form @submit.prevent="login" class="flex flex-col gap-4">

          <p class="flex flex-col">
            <label for="name" class="f-label">Name</label>
            <br>
            <input id="name" v-model="name" type="text" class="f-input" required /> 
          </p>

          <p class="flex flex-col">
            <label for="password" class="f-label">Password</label>
            <br>
            <input id="password" v-model="password" type="password" class="f-input" /> 
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
