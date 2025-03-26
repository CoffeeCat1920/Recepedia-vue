<script setup>
import { defineProps } from "vue";
import { RouterLink } from "vue-router";

const props = defineProps({
  name: String,
  uuid: String,
  views: Number,
});
</script>

<template>
  <div class="flex">
    <RouterLink :to="`/recipe/${uuid}`" class="recipe-card">
      <div class="card-container relative hover:scale-110 transition-transform cursor-pointer">
        <div class="card-content h-52 w-60">
          <h3 class="recipe-name font-monomakh text-center text-2xl text-amber-950">{{ name }}</h3>
          <p class="recipe-views absolute bottom-2 left-2 text-sm text-amber-900">Views: {{ views }}</p>
        </div>
      </div>
    </RouterLink>
  </div>
</template>
