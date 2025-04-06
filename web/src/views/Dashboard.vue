<script setup>
import DashboardRecipes from '@/components/DashboardRecipeBoard.vue';

import { ref, onMounted } from 'vue';

const user = ref("");
const uuid = ref("");
const loggedIn = ref(false);
const menuOpen = ref(false);

const LoginData = async () => {
  try {
    const response = await fetch("/api/data/login", {
      credentials: "include",
      method: "GET",
    });

    if (!response.ok) throw new Error("Unauthorized");

    const data = await response.json();

    user.value = data.name;
    uuid.value = data.uuid;
    console.log(data.uuid);
    loggedIn.value = data.loggedIn;

  } catch (error) {
    console.error("Error fetching login info:", error);
  }
}

const toggleMenu = () => {
  menuOpen.value = !menuOpen.value; // Toggle the menu open/close state
};

const deleteUser = async () => {
  try {
    const response = await fetch('/api/user/' + uuid.value, {
      credentials: "include",
      method: "DELETE",
    });

    if (response.ok) alert('Deleted user')
    else throw new Error("Can't delete user");

  } catch (error) {
    console.error(error);
  }
}

onMounted(LoginData);

</script>

<template>

  <div class="flex flex-col">

    <div class="relative">
      <!-- Hamburger menu icon -->
      <button @click="toggleMenu" class="absolute top-4 right-4 text-2xl">
        &#9776; <!-- Unicode character for hamburger icon -->
      </button>

      <!-- Dropdown menu for options -->
      <div v-if="menuOpen" class="cursor:pointer absolute top-12 right-4 bg-white border rounded-lg shadow-lg">
        <ul>
          <div class="card-container relative">
            <div class="card-content">
              <li class="cursor:pointer hover:font-bold" @click="deleteUser">
                Delete User
              </li>
            </div>
          </div>
        </ul>
      </div>
    </div>

    <div v-if="loggedIn">
      <h1 class="text-amber-950 font-monomakh text-6xl text-center pt-12 pb-5">{{ user }}</h1>
    </div>

    <DashboardRecipes />

  </div>

</template>
