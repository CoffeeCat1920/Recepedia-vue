<script setup>
  import DashboardRecipes from '@/components/DashboardRecipeBoard.vue'; 

  import { ref, onMounted } from 'vue';
  
  const user =  ref("");
  const loggedIn = ref(false);

  const LoginData = async () => {
    try {
      const response = await fetch("/api/data/login", {
        credentials: "include",
        method: "GET",
      });

      if (!response.ok) throw new Error("Unauthorized");

      const data = await response.json();

      user.value = data.name;  
      loggedIn.value = data.loggedIn;

    } catch (error) {
      console.error("Error fetching login info:", error);
    }
  }

  onMounted(LoginData);

</script>

<template>

  <div class="flex flex-col">

    <div v-if="loggedIn">
      <h1 class="text-amber-950 font-monomakh text-6xl text-center pt-12 pb-5">{{ user }}</h1>
    </div>

    <DashboardRecipes />

  </div>

</template>
