<script setup>
import { defineProps, watch } from "vue";
import { useRouter, RouterLink } from "vue-router";

const props = defineProps({
  name: String,
  uuid: String,
  views: Number,
});

const router = useRouter();

const deleteRecipes = async (id) => {
  try {
    const response = await fetch(`/api/admin/recipe/delete/${id} `, {
      method: 'PATCH',
      credentials: "include",
    });

    if (!response.ok) {
      throw new Error('Failed to delete recipe');
    }

    console.log(`Recipe with ID ${id} deleted successfully.`);

    location.reload();
  } catch (error) {
    console.error('Error:', error);
  }
};

</script>

<template>
  <div class="flex">
    <div class="card-container relative">
      <div class="card-content h-64 w-60">
        <h3 class="recipe-name font-monomakh text-center text-2xl text-amber-950">{{ name }}</h3>

        <div class="absolute bottom-2 left-2 text-sm text-amber-900">
          <p class="recipe-views flex">
            <span class="material-symbols-outlined"> visibility </span>
            {{ views }}
          </p>
        </div>

        <div class="absolute bottom-2 right-2 flex flex-col space-y-1">
          <button
            class="text-sm bg-red-400 shadow-md shadow-amber-900 px-2 hover:scale-110 hover:rotate-3 transition-transform"
            @click="deleteRecipes(uuid)">
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
