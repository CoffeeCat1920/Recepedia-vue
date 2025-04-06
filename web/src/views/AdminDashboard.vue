<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter, RouterLink } from "vue-router";

const loggedIn = ref(false);
const numberOfRecipes = ref(false);
const numberOfUsers = ref(false);

const router = useRouter();

const LoginData = async () => {
  try {
    const response = await fetch("/api/data/admin/login", {
      credentials: "include",
      method: "GET",
    });

    if (!response.ok) throw new Error("Unauthorized");

    const data = await response.json();

    loggedIn.value = data.loggedIn;

  } catch (error) {
    console.error("Error fetching login info:", error);
  }
}

const DashboardData = async () => {

  try {
    const response = await fetch("/api/data/admin/dashboard", {
      credentials: "include",
      method: "GET",
    })

    if (!response.ok) throw new Error("Can't get data")

    const data = await response.json();

    numberOfRecipes.value = data.numberOfRecipes;
    numberOfUsers.value = data.numberOfUsers;

  } catch (error) {
    console.error(error)
  }
}

const editRecipes = async () => {
  router.push(`/admin/editRecipes`);
}

const editUsers = async () => {
  router.push(`/admin/editUsers`);
}

onMounted(() => {
  LoginData();
  DashboardData();
});

</script>

<template>

  <div class="flex flex-col items-center text-center">
    <h1 class="text-6xl font-bold text-amber-900 font-[Monomakh] pt-[100px]">Admin Pannel</h1>
  </div>

  <div class="grid md:grid-cols-1 place-items-center gap-y-12">

    <!-- Cards -->
    <div class="flex">
      <div class="card-container relative">
        <div class="card-content h-64 w-60">
          <h3 class="recipe-name font-monomakh text-center text-2xl text-amber-950">Recipes</h3>
          <h2 class="recipe-name font-monomakh text-center text-1xl text-amber-900">Total: {{ numberOfRecipes }}</h2>
          <div class="absolute bottom-2 right-2 flex flex-col space-y-1">
            <button
              class="text-sm bg-green-400 shadow-md shadow-amber-900 px-2 hover:scale-110 hover:rotate-3 transition-transform"
              @click="editRecipes">
              Edit Recipes
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Cards -->
    <div class="flex">
      <div class="card-container relative">
        <div class="card-content h-64 w-60">
          <h3 class="recipe-name font-monomakh text-center text-2xl text-amber-950">Users</h3>
          <h2 class="recipe-name font-monomakh text-center text-1xl text-amber-900">Total: {{ numberOfUsers }}</h2>
          <h2></h2>
          <div class="absolute bottom-2 right-2 flex flex-col space-y-1">
            <button class="text-sm bg-green-400 shadow-md shadow-amber-900 px-2 hover:scale-110 hover:rotate-3
              transition-transform" @click="editUsers">
              Edit Users
            </button>
          </div>
        </div>
      </div>
    </div>

  </div>

</template>
