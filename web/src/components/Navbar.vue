<!-- src/components/Navbar.vue -->

<script setup>
import { useAuthStore } from "@/stores/authStore";
import { useRouter } from "vue-router";
import { onMounted } from "vue";

const router = useRouter();
const auth = useAuthStore();

const logout = async () => {
  await auth.logout(router);
};

onMounted(async () => {
  await auth.checkLogin();
});
</script>

<template>
  <nav class="bg-amber-100 shadow-md w-full py-4 px-6 flex justify-between items-center">
    <!-- Brand Logo -->
    <router-link to="/" class="text-2xl font-bold text-gray-900 hover:text-gray-700">
      •ᴗ• Recipidia
    </router-link>

    <!-- Navigation Links -->
    <div class="hidden md:flex space-x-6">
      <router-link to="/" class="text-gray-700 hover:text-black font-medium">Home</router-link>
      <router-link to="/about" class="text-gray-700 hover:text-black font-medium">About</router-link>

      <template v-if="auth.loggedIn">
        <router-link to="/uploadrecipe" class="text-gray-700 hover:text-black font-medium">Upload Recipe</router-link>
        <router-link to="/dashboard" class="text-gray-700 hover:text-black font-medium">Dashboard</router-link>
        <a href="#" @click.prevent="logout" class="text-red-600 hover:text-red-800 font-medium">Logout</a>
      </template>

      <template v-else>
        <router-link to="/signup" class="text-blue-600 hover:text-blue-800 font-medium">Sign Up</router-link>
        <router-link to="/login" class="text-blue-600 hover:text-blue-800 font-medium">Log In</router-link>
      </template>
    </div>

    <!-- Mobile Menu Button -->
    <button class="md:hidden text-gray-700 focus:outline-none">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16m-7 6h7" />
      </svg>
    </button>

  </nav>
</template>
