<script>
import { defineComponent, ref, onMounted } from 'vue';
import { useAuthStore } from '@/stores/authStore.ts';
import { useRouter } from 'vue-router';

export default defineComponent({
  setup() {
    const router = useRouter();
    
    // Declare refs for name and password
    const password = ref('');
    
    // Access the auth store
    const auth = useAuthStore();

    // Login function
    const login = async () => {
      await auth.adminLogin(password.value, router);
    };

    // Optional: Check login status when the component is mounted
    onMounted(async () => {
      // You can call your checkLogin method here if needed
      // await auth.checkLogin();
    });

    return {
      password,
      login,
    };
  }
});
</script>

<template>

  <div class="flex flex-col items-center text-center justify-center h-screen">

    <div class="card-container">

      <div class="card-content">

        <h1 class="title">Admin Log In</h1>

        <form @submit.prevent="login" class="flex flex-col gap-4">

          <p class="flex flex-col">
            <label for="password" class="f-label">Admin Password</label>
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
