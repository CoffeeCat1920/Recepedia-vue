<script setup>
import { useAuthStore } from "@/stores/authStore";
import { useRouter } from "vue-router";
import { onMounted, ref } from "vue";

const router = useRouter();
const auth = useAuthStore();
// Mobile menu state
const mobileMenuOpen = ref(false);

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value;
};

const logout = async () => {
  await auth.logout(router);
  mobileMenuOpen.value = false; // Close menu after logout
};

onMounted(async () => {
  await auth.checkLogin();
  await auth.checkAdminLogin();
});
</script>

<template>
  <nav class="bg-amber-100 shadow-md w-full py-4 px-6">
    <div class="flex justify-between items-center">
      <!-- Brand Logo -->
      <router-link to="/" class="text-2xl font-bold text-gray-900 hover:text-gray-700">
        •ᴗ• Recipidia
      </router-link>

      <!-- Desktop Navigation Links -->
      <div class="hidden md:flex space-x-6 items-center">
        <router-link to="/" class="text-gray-700 hover:text-black font-medium">Home</router-link>
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

      <!-- Mobile Menu Button -->
      <button class="md:hidden text-gray-700 focus:outline-none" @click="toggleMobileMenu">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Mobile Menu -->
    <div 
      v-show="mobileMenuOpen" 
      class="md:hidden mt-4 transition-all duration-300 ease-in-out"
    >
      <div class="flex flex-col space-y-3 py-2">
        <router-link @click="mobileMenuOpen = false" to="/" class="text-gray-700 hover:text-black font-medium px-2 py-1 rounded hover:bg-amber-200">
          Home
        </router-link>
        <template v-if="auth.loggedIn">
          <router-link @click="mobileMenuOpen = false" to="/uploadrecipe" class="text-gray-700 hover:text-black font-medium px-2 py-1 rounded hover:bg-amber-200">
            Upload Recipe
          </router-link>
          <router-link @click="mobileMenuOpen = false" to="/dashboard" class="text-gray-700 hover:text-black font-medium px-2 py-1 rounded hover:bg-amber-200">
            Dashboard
          </router-link>
          <a href="#" @click.prevent="logout" class="text-red-600 hover:text-red-800 font-medium px-2 py-1 rounded hover:bg-amber-200">
            Logout
          </a>
        </template>
        <template v-else>
          <router-link @click="mobileMenuOpen = false" to="/signup" class="text-blue-600 hover:text-blue-800 font-medium px-2 py-1 rounded hover:bg-amber-200">
            Sign Up
          </router-link>
          <router-link @click="mobileMenuOpen = false" to="/login" class="text-blue-600 hover:text-blue-800 font-medium px-2 py-1 rounded hover:bg-amber-200">
            Log In
          </router-link>
        </template>
        <template v-if="auth.adminLoggedIn">
          <div class="border-t border-gray-300 my-2"></div>
          <router-link @click="mobileMenuOpen = false" to="/admindashboard" class="text-blue-600 hover:text-blue-800 font-medium px-2 py-1 rounded hover:bg-amber-200">
            Admin Menu
          </router-link>
        </template>
      </div>
    </div>
  </nav>
</template>
