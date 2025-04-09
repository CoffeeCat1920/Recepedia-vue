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
  await auth.checkAdminLogin();
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

      <template v-if="auth.adminLoggedIn">
        <p class="text-gray-700">|</p>

        <router-link to="/admindashboard" class="text-blue-600 hover:text-blue-800 font-medium">Admin Menu</router-link>
      </template>

    </div>

  </nav>

</template>
