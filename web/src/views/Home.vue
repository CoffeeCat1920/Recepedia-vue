<!-- /src/views/Home.vue -->
<script setup>
import Brand from '@/components/Brand.vue';
import Recommendation from '@/components/RecipeRecommendation.vue';
import Searchbar from '@/components/Searchbar.vue';
</script>


<template>
  <div class="flex flex-col items-center text-center pt-10 px-4 sm:px-6 lg:px-8">
    <div class="card-container w-full max-w-6xl">
      <div class="card-content pt-10">
        <Brand />
        <Searchbar />
      </div>
    </div>
  </div>

  <Recommendation />
</template>
